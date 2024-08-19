package golang

import (
	"context"
	_ "embed"
	"fmt"
	"io/fs"
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

const coverageReportDir = "coverage"

// BuildConfig is the configuration for building binaries.
type BuildConfig struct {
	// Platform is the platform to build the binary for.
	Platform tools.Platform

	// PackagePath is the path to package to build relative to the ModulePath.
	PackagePath string

	// BinOutputPath is the path for compiled binary file.
	BinOutputPath string

	// CGOEnabled builds cgo binary.
	CGOEnabled bool

	// Tags is go build tags.
	Tags []string
}

// Build builds go binary.
func Build(ctx context.Context, deps build.DepsFunc, config BuildConfig) error {
	if config.Platform.OS == tools.OSDocker {
		return buildInDocker(ctx, deps, config)
	}
	return buildLocally(ctx, deps, config)
}

// Lint lints the go code.
func Lint(ctx context.Context, deps build.DepsFunc) error {
	deps(EnsureGo, EnsureGolangCI, storeLintConfig)

	log := logger.Get(ctx)
	config := lintConfigPath(ctx)

	return helpers.OnModule("go.mod", func(path string) error {
		goCodePresent, err := containsGoCode(path)
		if err != nil {
			return err
		}
		if !goCodePresent {
			log.Info("No code to lint", zap.String("path", path))
			return nil
		}

		log.Info("Running linter", zap.String("path", path))
		cmd := exec.Command(tools.Bin(ctx, "bin/golangci-lint", tools.PlatformLocal), "run", "--config", config)
		cmd.Env = env(ctx)
		cmd.Dir = path
		if err := libexec.Exec(ctx, cmd); err != nil {
			return errors.Wrapf(err, "linter errors found in module '%s'", path)
		}
		return nil
	})
}

// Tidy runs go mod tidy in repository.
func Tidy(ctx context.Context, deps build.DepsFunc) error {
	deps(EnsureGo)

	log := logger.Get(ctx)
	return helpers.OnModule("go.mod", func(path string) error {
		log.Info("Running go mod tidy", zap.String("path", path))

		cmd := exec.Command(tools.Bin(ctx, "bin/go", tools.PlatformLocal), "mod", "tidy")
		cmd.Env = env(ctx)
		cmd.Dir = path
		if err := libexec.Exec(ctx, cmd); err != nil {
			return errors.Wrapf(err, "'go mod tidy' failed in module '%s'", path)
		}
		return nil
	})
}

// UnitTests runs go unit tests in repository.
func UnitTests(ctx context.Context, deps build.DepsFunc) error {
	deps(EnsureGo)

	log := logger.Get(ctx)

	if err := os.MkdirAll(coverageReportDir, 0o700); err != nil {
		return errors.WithStack(err)
	}

	rootDir := filepath.Dir(lo.Must(filepath.Abs(lo.Must(filepath.EvalSymlinks(lo.Must(os.Getwd()))))))
	return helpers.OnModule("go.mod", func(path string) error {
		path = lo.Must(filepath.Abs(lo.Must(filepath.EvalSymlinks(path))))
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return errors.WithStack(err)
		}

		goCodePresent, err := containsGoCode(path)
		if err != nil {
			return err
		}
		if !goCodePresent {
			log.Info("No code to test", zap.String("path", path))
			return nil
		}

		coverageName := strings.ReplaceAll(relPath, "/", "-")
		coverageProfile := filepath.Join(coverageReportDir, coverageName)

		log.Info("Running go tests", zap.String("path", path))
		cmd := exec.Command(
			tools.Bin(ctx, "bin/go", tools.PlatformLocal),
			"test",
			"-count=1",
			"-shuffle=on",
			"-race",
			"-cover", "./...",
			"-coverpkg", "./...",
			"-coverprofile", coverageProfile,
			"./...",
		)
		cmd.Env = env(ctx)
		cmd.Dir = path
		if err := libexec.Exec(ctx, cmd); err != nil {
			return errors.Wrapf(err, "unit tests failed in module '%s'", path)
		}
		return nil
	})
}

func buildLocally(ctx context.Context, deps build.DepsFunc, config BuildConfig) error {
	deps(EnsureGo)

	if config.Platform != tools.PlatformLocal {
		return errors.Errorf("building requested for platform %s while only %s is supported",
			config.Platform, tools.PlatformLocal)
	}

	args, envs := buildArgsAndEnvs(ctx, config)

	cmd := exec.Command(tools.Bin(ctx, "bin/go", config.Platform), args...)
	cmd.Dir = config.PackagePath
	cmd.Env = append(os.Environ(), envs...)

	logger.Get(ctx).Info(
		"Building go package locally",
		zap.String("package", config.PackagePath),
		zap.String("output", config.BinOutputPath),
		zap.String("command", cmd.String()),
	)
	if err := libexec.Exec(ctx, cmd); err != nil {
		return errors.Wrapf(err, "building go package '%s' failed", config.PackagePath)
	}
	return nil
}

func buildInDocker(ctx context.Context, deps build.DepsFunc, config BuildConfig) error {
	deps(docker.EnsureDocker)

	goTool, err := tools.Get(Go)
	if err != nil {
		return err
	}

	image := fmt.Sprintf("golang:%s-alpine%s", goTool.GetVersion(), docker.AlpineVersion)

	srcDir := lo.Must(filepath.EvalSymlinks(lo.Must(filepath.Abs("."))))
	envDir := tools.EnvDir(ctx)

	if err := os.MkdirAll(envDir, 0o755); err != nil {
		return errors.WithStack(err)
	}

	args, envs := buildArgsAndEnvs(ctx, config)
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
		"--name", "sei-build-golang",
	}

	for _, env := range envs {
		runArgs = append(runArgs, "--env", env)
	}

	runArgs = append(runArgs, image, "/usr/local/go/bin/go")
	runArgs = append(runArgs, args...)

	cmd := exec.Command("docker", runArgs...)
	logger.Get(ctx).Info(
		"Building go package in docker",
		zap.String("package", config.PackagePath),
		zap.String("command", cmd.String()),
	)
	if err := libexec.Exec(ctx, cmd); err != nil {
		return errors.Wrapf(err, "building package '%s' failed", config.PackagePath)
	}
	return nil
}

func buildArgsAndEnvs(ctx context.Context, config BuildConfig) (args, envs []string) {
	ldFlags := []string{"-w", "-s"}

	args = []string{
		"build",
		"-trimpath",
		"-buildvcs=false",
		"-ldflags=" + strings.Join(ldFlags, " "),
		"-o", lo.Must(filepath.Abs(config.BinOutputPath)),
		".",
	}
	if len(config.Tags) != 0 {
		args = append(args, "-tags="+strings.Join(config.Tags, ","))
	}

	goOS := config.Platform.OS
	if goOS == tools.OSDocker {
		goOS = tools.OSLinux
	}

	cgoEnabled := "0"
	if config.CGOEnabled {
		cgoEnabled = "1"
	}
	envs = append(env(ctx),
		"CGO_ENABLED="+cgoEnabled,
		"GOOS="+goOS,
		"GOARCH="+config.Platform.Arch,
	)

	return args, envs
}

func containsGoCode(path string) (bool, error) {
	errFound := errors.New("found")
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}
		return errFound
	})
	if errors.Is(err, errFound) {
		return true, nil
	}
	return false, errors.WithStack(err)
}

//go:embed "golangci.yaml"
var lintConfig []byte

func storeLintConfig(ctx context.Context, _ build.DepsFunc) error {
	return errors.WithStack(os.WriteFile(lintConfigPath(ctx), lintConfig, 0o600))
}

func lintConfigPath(ctx context.Context) string {
	return filepath.Join(tools.VersionDir(ctx, tools.PlatformLocal), "golangci.yaml")
}

func env(ctx context.Context) []string {
	return []string{
		"PATH=" + filepath.Join(tools.VersionDir(ctx, tools.PlatformLocal), "bin") + ":" + os.Getenv("PATH"),
		"GOPATH=" + filepath.Join(tools.DevDir(ctx), "go"),
		"GOCACHE=" + filepath.Join(tools.DevDir(ctx), "go", "cache", "gobuild"),
		"GOLANGCI_LINT_CACHE=" + filepath.Join(tools.DevDir(ctx), "go", "cache", "golangci"),
	}
}
