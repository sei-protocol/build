package rust

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/outofforest/build"
	"github.com/outofforest/libexec"
	"github.com/outofforest/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/sei-protocol/build/pkg/helpers"
	"github.com/sei-protocol/build/pkg/tools"
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

// Build builds go binary.
func Build(ctx context.Context, deps build.DepsFunc, config BuildConfig) error {
	if config.Platform.OS == tools.OSDocker {
		return errors.New("building in docker hasn't been implemented yet")
		// return buildInDocker(ctx, config)
	}
	return buildLocally(ctx, deps, config)
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
		"PATH=" + os.Getenv("PATH"),
		"CARGO_HOME=" + filepath.Join(tools.DevDir(ctx), "rust", "cargo"),
	}
}

func targetDir(ctx context.Context) string {
	return filepath.Join(tools.DevDir(ctx), "rust", "target")
}
