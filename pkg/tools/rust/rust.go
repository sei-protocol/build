package rust

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/outofforest/build"
	"github.com/outofforest/libexec"
	"github.com/outofforest/logger"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/sei-protocol/build/pkg/helpers"
	"github.com/sei-protocol/build/pkg/tools"
	"github.com/sei-protocol/build/pkg/tools/docker"
)

// BuildConfig is the configuration for building binaries.
type BuildConfig struct {
	// Platform is the platform to build the binary for.
	Platform tools.Platform

	// PackagePath is the path to package to build relative to the ModulePath.
	PackagePath string

	// Binary is the name of the binary to build as specified in Cargo.toml
	Binary string

	// BinOutputPath is the path for compiled binary file.
	BinOutputPath string
}

// Build builds rust binary.
func Build(ctx context.Context, deps build.DepsFunc, config BuildConfig) error {
	if config.Platform.OS == tools.OSDocker {
		return buildInDocker(ctx, deps, config)
	}
	return buildLocally(ctx, deps, config)
}

// Lint lints the rust code.
func Lint(ctx context.Context, deps build.DepsFunc) error {
	deps(EnsureRust)

	log := logger.Get(ctx)

	return helpers.OnModule("Cargo.toml", func(path string) error {
		log.Info("Running linter", zap.String("path", path))
		cmd := exec.Command(tools.Bin(ctx, "bin/cargo", tools.PlatformLocal), "clippy",
			"--target-dir", targetDir(ctx))
		cmd.Env = env(ctx)
		cmd.Dir = path
		if err := libexec.Exec(ctx, cmd); err != nil {
			return errors.Wrapf(err, "linter errors found in module '%s'", path)
		}
		return nil
	})
}

// UnitTests runs rust unit tests in repository.
func UnitTests(ctx context.Context, deps build.DepsFunc) error {
	deps(EnsureRust)

	log := logger.Get(ctx)

	return helpers.OnModule("Cargo.toml", func(path string) error {
		path = lo.Must(filepath.Abs(lo.Must(filepath.EvalSymlinks(path))))

		log.Info("Running rust tests", zap.String("path", path))
		cmd := exec.Command(tools.Bin(ctx, "bin/cargo", tools.PlatformLocal), "test",
			"--target-dir", targetDir(ctx))
		cmd.Env = env(ctx)
		cmd.Dir = path
		if err := libexec.Exec(ctx, cmd); err != nil {
			return errors.Wrapf(err, "unit tests failed in module '%s'", path)
		}
		return nil
	})
}

func buildLocally(ctx context.Context, deps build.DepsFunc, config BuildConfig) error {
	deps(EnsureRust)

	if config.Platform != tools.PlatformLocal {
		return errors.Errorf("building requested for platform %s while only %s is supported",
			config.Platform, tools.PlatformLocal)
	}

	args, envs := buildArgsAndEnvs(ctx)
	args = append(args, "--bin", config.Binary)

	cmd := exec.Command(tools.Bin(ctx, "bin/cargo", config.Platform), args...)
	cmd.Dir = config.PackagePath
	cmd.Env = append(os.Environ(), envs...)

	logger.Get(ctx).Info(
		"Building rust package locally",
		zap.String("package", config.PackagePath),
		zap.String("binary", config.Binary),
		zap.String("command", cmd.String()),
	)
	if err := libexec.Exec(ctx, cmd); err != nil {
		return errors.Wrapf(err, "building rust package '%s' failed", config.PackagePath)
	}

	return helpers.CopyFile(config.BinOutputPath, filepath.Join(targetDir(ctx), "release", config.Binary), 0o755)
}

func buildInDocker(ctx context.Context, deps build.DepsFunc, config BuildConfig) error {
	deps(docker.EnsureDocker)

	rustTool, err := tools.Get(Rust)
	if err != nil {
		return err
	}

	image := fmt.Sprintf("rust:%s-alpine%s", rustTool.GetVersion(), docker.AlpineVersion)

	srcDir := lo.Must(filepath.EvalSymlinks(lo.Must(filepath.Abs("."))))
	envDir := tools.EnvDir(ctx)

	if err := os.MkdirAll(envDir, 0o755); err != nil {
		return errors.WithStack(err)
	}

	args, envs := buildArgsAndEnvs(ctx)
	if err != nil {
		return err
	}
	runArgs := []string{
		"run", "--rm",
		"--label", docker.LabelKey + "=" + docker.LabelValue,
		"-v", srcDir + ":" + srcDir,
		"-v", envDir + ":" + envDir,
		"--workdir", filepath.Join(srcDir, config.PackagePath),
		"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		"--name", "sei-build-rust",
	}

	for _, env := range envs {
		if strings.HasPrefix(env, "PATH=") {
			env = strings.Replace(env, "=", "=/usr/local/cargo/bin:", 1)
		}
		runArgs = append(runArgs, "--env", env)
	}

	runArgs = append(runArgs, image, "cargo")
	runArgs = append(runArgs, args...)

	cmd := exec.Command("docker", runArgs...)
	logger.Get(ctx).Info(
		"Building rust package in docker",
		zap.String("package", config.PackagePath),
		zap.String("command", cmd.String()),
	)
	if err := libexec.Exec(ctx, cmd); err != nil {
		return errors.Wrapf(err, "building package '%s' failed", config.PackagePath)
	}

	return helpers.CopyFile(config.BinOutputPath, filepath.Join(targetDir(ctx), "release", config.Binary), 0o755)
}

func buildArgsAndEnvs(ctx context.Context) (args, envs []string) {
	args = []string{
		"build",
		"--release",
		"--target-dir", targetDir(ctx),
	}

	return args, env(ctx)
}

func env(ctx context.Context) []string {
	return []string{
		"PATH=" + filepath.Join(tools.VersionDir(ctx, tools.PlatformLocal), "bin") + ":" + os.Getenv("PATH"),
		"CARGO_HOME=" + filepath.Join(tools.DevDir(ctx), "rust", "cargo"),
	}
}

func targetDir(ctx context.Context) string {
	return filepath.Join(tools.DevDir(ctx), "rust", "target")
}
