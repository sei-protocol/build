package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	osexec "os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/outofforest/libexec"
	"github.com/outofforest/logger"
	"github.com/outofforest/parallel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/sei-protocol/build/pkg/tools/docker"
)

const (
	labelEnv    = "io.seinetwork.localnet.env"
	envName     = "localnet"
	networkName = "localnet"
)

// NewDocker creates new docker target.
func NewDocker() *Docker {
	return &Docker{}
}

// Docker is the target deploying apps to docker.
type Docker struct {
	mu            sync.Mutex
	networkExists bool
}

// Stop stops running applications.
func (d *Docker) Stop(ctx context.Context, appSet AppSet) error {
	dependencies := map[string][]chan struct{}{}
	readyChs := map[string]chan struct{}{}
	for _, app := range appSet {
		readyCh := make(chan struct{})
		readyChs[app.Name] = readyCh

		cID, err := containerExists(ctx, app.Name)
		if err != nil {
			return err
		}
		if cID == "" {
			close(readyCh)
		}

		for _, dep := range app.Requires.Dependencies {
			dependencies[dep.Name] = append(dependencies[dep.Name], readyCh)
		}
	}

	return forContainer(ctx, func(ctx context.Context, info container) error {
		log := logger.Get(ctx).With(zap.String("id", info.ID), zap.String("name", info.Name))

		if appSet.App(info.Name) == nil {
			log.Info("Unexpected container found, deleting it")

			if err := removeContainer(ctx, info); err != nil {
				return err
			}

			log.Info("Container deleted")
			return nil
		}

		if deps := dependencies[info.Name]; len(deps) > 0 {
			log.Info("Waiting for dependencies to be stopped")
			for _, depCh := range deps {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-depCh:
				}
			}
		}

		log.Info("Stopping container")

		if err := libexec.Exec(ctx, noStdout(docker.Cmd("stop", "--time", "60", info.ID))); err != nil {
			return errors.Wrapf(err, "stopping container `%s` failed", info.Name)
		}

		log.Info("Container stopped")
		close(readyChs[info.Name])
		return nil
	})
}

// Remove removes running applications.
func (d *Docker) Remove(ctx context.Context) error {
	err := forContainer(ctx, func(ctx context.Context, info container) error {
		log := logger.Get(ctx).With(zap.String("id", info.ID), zap.String("name", info.Name))
		log.Info("Deleting container")

		if err := removeContainer(ctx, info); err != nil {
			return err
		}

		log.Info("Container deleted")
		return nil
	})
	if err != nil {
		return err
	}
	return d.deleteNetwork(ctx)
}

// Deploy deploys environment to docker target.
func (d *Docker) Deploy(ctx context.Context, appSet AppSet) error {
	if err := appSet.Deploy(ctx, d); err != nil {
		return err
	}

	log := logger.Get(ctx)
	log.Info("Waiting until all applications start.")
	waitCtx, waitCancel := context.WithTimeout(ctx, 5*time.Minute)
	defer waitCancel()
	defer log.Info("All applications are healthy.")

	return WaitUntilHealthy(waitCtx, AppSetToHealthChecks(appSet))
}

// DeployContainer starts container in docker.
func (d *Docker) DeployContainer(ctx context.Context, app *App) (DeploymentInfo, error) {
	if err := d.ensureNetwork(ctx); err != nil {
		return DeploymentInfo{}, err
	}

	log := logger.Get(ctx).With(zap.String("appName", app.Name))
	log.Info("Starting container")

	id, err := containerExists(ctx, app.Name)
	if err != nil {
		return DeploymentInfo{}, err
	}

	var startCmd *osexec.Cmd
	if id != "" {
		startCmd = docker.Cmd("start", id)
	} else {
		runArgs := d.prepareRunArgs(app)
		startCmd = docker.Cmd(runArgs...)
	}
	idBuf := &bytes.Buffer{}
	startCmd.Stdout = idBuf

	if err := libexec.Exec(ctx, startCmd); err != nil {
		return DeploymentInfo{}, err
	}

	log.Info("Container started", zap.String("id", strings.TrimSuffix(idBuf.String(), "\n")))

	// FromHostIP = ipLocalhost here means that application is available on host's localhost, not container's localhost
	return DeploymentInfo{
		Status:            AppStatusRunning,
		HostFromHost:      "localhost",
		HostFromContainer: app.Name,
	}, nil
}

func (d *Docker) prepareRunArgs(app *App) []string {
	runArgs := []string{
		"run", "--name", app.Name, "-d", "--label", labelEnv + "=" + envName, "--network", networkName,
	}
	if app.RunAsUser {
		runArgs = append(runArgs, "--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()))
	}
	for _, port := range app.Ports {
		portStr := strconv.Itoa(port)
		runArgs = append(runArgs, "-p", "127.0.0.1:"+portStr+":"+portStr+"/tcp")
	}
	for _, v := range app.Volumes {
		runArgs = append(runArgs, "-v", lo.Must(filepath.EvalSymlinks(lo.Must(filepath.Abs(v.Source))))+":"+
			v.Destination)
	}
	if app.EnvVarsFunc != nil {
		for _, env := range app.EnvVarsFunc() {
			runArgs = append(runArgs, "-e", env.Name+"="+env.Value)
		}
	}

	if app.Entrypoint != "" {
		runArgs = append(runArgs, "--entrypoint", app.Entrypoint)
	}

	runArgs = append(runArgs, app.Image)
	if app.ArgsFunc != nil {
		runArgs = append(runArgs, app.ArgsFunc()...)
	}

	return runArgs
}

func (d *Docker) ensureNetwork(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.networkExists {
		return nil
	}

	log := logger.Get(ctx).With(zap.String("network", networkName))

	var err error
	d.networkExists, err = networkExists(ctx)
	if err != nil {
		return err
	}
	if d.networkExists {
		log.Info("Docker network exists")
		return nil
	}

	log.Info("Creating docker network")

	if err := libexec.Exec(ctx, noStdout(docker.Cmd("network", "create", networkName))); err != nil {
		return errors.Wrapf(err, "creating network '%s' failed", networkName)
	}

	d.networkExists = true
	log.Info("Docker network created")
	return nil
}

func (d *Docker) deleteNetwork(ctx context.Context) error {
	exists, err := networkExists(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	log := logger.Get(ctx).With(zap.String("network", networkName))
	log.Info("Deleting docker network")

	if err := libexec.Exec(ctx, noStdout(docker.Cmd("network", "rm", networkName))); err != nil {
		return errors.Wrapf(err, "deleting network '%s' failed", networkName)
	}

	log.Info("Docker network deleted")
	return nil
}

func containerExists(ctx context.Context, name string) (string, error) {
	idBuf := &bytes.Buffer{}
	existsCmd := docker.Cmd("ps", "-aq", "--no-trunc", "--filter", "name="+name)
	existsCmd.Stdout = idBuf
	if err := libexec.Exec(ctx, existsCmd); err != nil {
		return "", err
	}
	return strings.TrimSuffix(idBuf.String(), "\n"), nil
}

type container struct {
	ID      string
	Name    string
	Running bool
}

func forContainer(ctx context.Context, fn func(ctx context.Context, info container) error) error {
	listBuf := &bytes.Buffer{}
	listCmd := docker.Cmd("ps", "-aq", "--no-trunc", "--filter", "label="+labelEnv+"="+envName)
	listCmd.Stdout = listBuf
	if err := libexec.Exec(ctx, listCmd); err != nil {
		return err
	}

	listStr := strings.TrimSuffix(listBuf.String(), "\n")
	if listStr == "" {
		return nil
	}

	inspectBuf := &bytes.Buffer{}
	inspectCmd := docker.Cmd(append([]string{"inspect"}, strings.Split(listStr, "\n")...)...)
	inspectCmd.Stdout = inspectBuf

	if err := libexec.Exec(ctx, inspectCmd); err != nil {
		return err
	}

	var info []struct {
		ID    string `json:"Id"` //nolint:tagliatelle // `Id` is defined by docker
		Name  string
		State struct {
			Running bool
		}
		Config struct {
			Labels map[string]string
		}
	}

	if err := json.Unmarshal(inspectBuf.Bytes(), &info); err != nil {
		return errors.Wrap(err, "unmarshalling container properties failed")
	}

	return parallel.Run(ctx, func(ctx context.Context, spawn parallel.SpawnFn) error {
		for _, cInfo := range info {
			spawn("container."+cInfo.ID, parallel.Continue, func(ctx context.Context) error {
				return fn(ctx, container{
					ID:      cInfo.ID,
					Name:    strings.TrimPrefix(cInfo.Name, "/"),
					Running: cInfo.State.Running,
				})
			})
		}
		return nil
	})
}

func noStdout(cmd *osexec.Cmd) *osexec.Cmd {
	cmd.Stdout = io.Discard
	return cmd
}

func removeContainer(ctx context.Context, info container) error {
	cmds := []*osexec.Cmd{}
	if info.Running {
		// Everything will be removed, so we don't care about graceful shutdown
		cmds = append(cmds, noStdout(docker.Cmd("kill", info.ID)))
	}
	if err := libexec.Exec(ctx, append(cmds, noStdout(docker.Cmd("rm", info.ID)))...); err != nil {
		return errors.Wrapf(err, "deleting container `%s` failed", info.Name)
	}
	return nil
}

func networkExists(ctx context.Context) (bool, error) {
	buf := &bytes.Buffer{}
	cmd := docker.Cmd("network", "ls", "-q", "--no-trunc", "--filter", "name="+networkName)
	cmd.Stdout = buf
	if err := libexec.Exec(ctx, cmd); err != nil {
		return false, err
	}
	return strings.TrimSuffix(buf.String(), "\n") != "", nil
}
