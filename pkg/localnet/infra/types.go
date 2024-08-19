package infra

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/outofforest/libexec"
	"github.com/outofforest/logger"
	"github.com/outofforest/parallel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/sei-protocol/build/pkg/tools/docker"
)

// AppSet is the list of applications to deploy.
type AppSet []*App

// Deploy deploys app in environment to the target.
func (as AppSet) Deploy(ctx context.Context, t *Docker) error {
	log := logger.Get(ctx)
	log.Info(fmt.Sprintf("Staring AppSet deployment, apps: %s", strings.Join(lo.Map(as, func(app *App, _ int) string {
		return app.Name
	}), ",")))

	return parallel.Run(ctx, func(ctx context.Context, spawn parallel.SpawnFn) error {
		deploymentSlots := make(chan struct{}, runtime.NumCPU())
		for i := 0; i < cap(deploymentSlots); i++ {
			deploymentSlots <- struct{}{}
		}
		imagePullSlots := make(chan struct{}, 3)
		for i := 0; i < cap(imagePullSlots); i++ {
			imagePullSlots <- struct{}{}
		}

		deployments := map[string]struct {
			App          *App
			ImageReadyCh chan struct{}
			ReadyCh      chan struct{}
		}{}
		images := map[string]chan struct{}{}
		for _, app := range as {
			if _, exists := images[app.Image]; !exists {
				ch := make(chan struct{}, 1)
				ch <- struct{}{}
				images[app.Image] = ch
			}
			deployments[app.Name] = struct {
				App          *App
				ImageReadyCh chan struct{}
				ReadyCh      chan struct{}
			}{
				App:          app,
				ImageReadyCh: images[app.Image],
				ReadyCh:      make(chan struct{}),
			}
		}
		for name, toDeploy := range deployments {
			if toDeploy.App.info.Status == AppStatusRunning {
				close(toDeploy.ReadyCh)
				continue
			}

			toDeploy := toDeploy
			spawn("deploy."+name, parallel.Continue, func(ctx context.Context) error {
				log.Info("Deployment initialized")

				if err := ensureDockerImage(ctx, toDeploy.App.Image, imagePullSlots, toDeploy.ImageReadyCh); err != nil {
					return err
				}

				if dependencies := toDeploy.App.Requires.Dependencies; len(dependencies) > 0 {
					depNames := make([]string, 0, len(dependencies))
					for _, d := range dependencies {
						depNames = append(depNames, d.Name)
					}
					log.Info("Waiting for dependencies", zap.Strings("dependencies", depNames))
					for _, name := range depNames {
						select {
						case <-ctx.Done():
							return errors.WithStack(ctx.Err())
						case <-deployments[name].ReadyCh:
						}
					}
					log.Info("Dependencies are running now")
				}

				log.Info("Waiting for free slot for deploying the application")
				select {
				case <-ctx.Done():
					return errors.WithStack(ctx.Err())
				case <-deploymentSlots:
				}

				log.Info("Deployment started")

				if err := toDeploy.App.Deploy(ctx, t); err != nil {
					return err
				}

				log.Info("Deployment succeeded")

				close(toDeploy.ReadyCh)
				deploymentSlots <- struct{}{}
				return nil
			})
		}
		return nil
	})
}

// App returns an app by name.
func (as AppSet) App(name string) *App {
	for _, app := range as {
		if app.Name == name {
			return app
		}
	}
	return nil
}

func ensureDockerImage(ctx context.Context, image string, slots, readyCh chan struct{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case _, ok := <-readyCh:
		// If <-imageReadyCh blocks it means another goroutine is pulling the image
		if !ok {
			// Channel is closed, image has been already pulled, we are ready to go
			return nil
		}
	}

	log := logger.Get(ctx).With(zap.String("image", image))

	imageBuf := &bytes.Buffer{}
	imageCmd := docker.Cmd("images", "-q", image)
	imageCmd.Stdout = imageBuf
	if err := libexec.Exec(ctx, imageCmd); err != nil {
		return errors.Wrapf(err, "failed to list image '%s'", image)
	}
	if imageBuf.Len() > 0 {
		log.Info("Docker image exists")
		close(readyCh)
		return nil
	}

	log.Info("Waiting for free slot for pulling the docker image")

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-slots:
	}

	log.Info("Pulling docker image")

	if err := libexec.Exec(ctx, docker.Cmd("pull", image)); err != nil {
		return errors.Wrapf(err, "failed to pull docker image '%s'", image)
	}

	log.Info("Image pulled")
	close(readyCh)
	slots <- struct{}{}
	return nil
}

// Prerequisites specifies list of other apps which have to be healthy before app may be started.
type Prerequisites struct {
	// Timeout tells how long we should wait for prerequisite to become healthy
	Timeout time.Duration

	// Dependencies specifies a list of other apps this app depends on
	Dependencies AppSet
}

// EnvVar is used to define environment variable for docker container.
type EnvVar struct {
	Name  string
	Value string
}

// Volume defines volume to be mounted inside container.
type Volume struct {
	Source      string
	Destination string
}

// App represents application to be deployed.
type App struct {
	// Name of the application
	Name string

	// ArgsFunc is the function returning args passed to binary
	ArgsFunc func() []string

	// Ports are the network ports exposed by the application
	Ports map[string]int

	// Requires is the list of health checks to be required before app can be deployed
	Requires Prerequisites

	// PrepareFunc is the function called before application is deployed for the first time.
	// It is a good place to prepare configuration files and other things which must or might
	// be done before application runs.
	PrepareFunc func(ctx context.Context) error

	// ConfigureFunc is the function called after application is deployed for the first time.
	// It is a good place to connect to the application to configure it because at this stage
	// the app's IP address is known.
	ConfigureFunc func(ctx context.Context, deployment DeploymentInfo) error

	// HealthCheckFunc is the function called to check if the application is ready.
	HealthCheckFunc func(ctx context.Context) error

	// Image is the url of the container image
	Image string

	// EnvVarsFunc is a function defining environment variables for docker container
	EnvVarsFunc func() []EnvVar

	// Volumes defines volumes to be mounted inside the container
	Volumes []Volume

	// RunAsUser set to true causes the container to be run using uid and gid of current user.
	// It is required if container creates files inside mounted directory which is a part of app's home.
	// Otherwise, `localnet` won't be able to delete them.
	RunAsUser bool

	// Entrypoint is the custom entrypoint for the container.
	Entrypoint string

	info DeploymentInfo
}

// Deploy deploys container to the target.
func (app *App) Deploy(ctx context.Context, target *Docker) error {
	if err := app.preprocess(ctx); err != nil {
		return err
	}

	var err error
	app.info, err = target.DeployContainer(ctx, app)
	if err != nil {
		return err
	}

	return app.postprocess(ctx)
}

// Info returns deployment info.
func (app *App) Info() DeploymentInfo {
	return app.info
}

func (app *App) preprocess(ctx context.Context) error {
	if len(app.Requires.Dependencies) > 0 {
		waitCtx, waitCancel := context.WithTimeout(ctx, app.Requires.Timeout)
		defer waitCancel()
		if err := WaitUntilHealthy(waitCtx, AppSetToHealthChecks(app.Requires.Dependencies)); err != nil {
			return err
		}
	}

	if app.info.Status == AppStatusStopped {
		return nil
	}

	if app.PrepareFunc != nil {
		return app.PrepareFunc(ctx)
	}
	return nil
}

func (app *App) postprocess(ctx context.Context) error {
	if app.info.Status == AppStatusStopped {
		return nil
	}
	if app.ConfigureFunc != nil {
		return app.ConfigureFunc(ctx, app.info)
	}
	return nil
}

// AppStatus describes current status of an application.
type AppStatus string

const (
	// AppStatusNotDeployed ,eans that app has been never deployed.
	AppStatusNotDeployed AppStatus = ""

	// AppStatusRunning means that app is running.
	AppStatusRunning AppStatus = "running"

	// AppStatusStopped means app was running but now is stopped.
	AppStatusStopped AppStatus = "stopped"
)

// DeploymentInfo contains info about deployed application.
type DeploymentInfo struct {
	// HostFromHost is the host's hostname application binds to
	HostFromHost string

	// HostFromContainer is the container's hostname application is listening on
	HostFromContainer string

	// Status indicates the status of the application
	Status AppStatus
}
