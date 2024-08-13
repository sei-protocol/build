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

	"github.com/sei-protocol/build/pkg/tools"
	"github.com/sei-protocol/build/pkg/tools/git"
)

const coverageReportDir = "coverage"

// Lint lints.
func Lint(ctx context.Context, deps build.DepsFunc) error {
	deps(Tidy, LintCode, git.StatusClean)
	log := logger.Get(ctx)
	config := lintConfigPath(ctx)

	return onModule(func(path string) error {
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
		cmd.Dir = path
		if err := libexec.Exec(ctx, cmd); err != nil {
			return errors.Wrapf(err, "linter errors found in module '%s'", path)
		}
		return nil
	})
}

// LintCode lints the code.
func LintCode(ctx context.Context, deps build.DepsFunc) error {
	deps(EnsureGo, EnsureGolangCI, storeLintConfig)
	log := logger.Get(ctx)
	config := lintConfigPath(ctx)

	fmt.Println(config)

	return onModule(func(path string) error {
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
	return onModule(func(path string) error {
		log.Info("Running go mod tidy", zap.String("path", path))

		cmd := exec.Command(tools.Bin(ctx, "bin/go", tools.PlatformLocal), "mod", "tidy")
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
	return onModule(func(path string) error {
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

		if filepath.Base(path) == "integration-tests" {
			log.Info("Skipping integration-tests", zap.String("path", path))
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
		cmd.Dir = path
		if err := libexec.Exec(ctx, cmd); err != nil {
			return errors.Wrapf(err, "unit tests failed in module '%s'", path)
		}
		return nil
	})
}

func onModule(fn func(path string) error) error {
	return filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "go.mod" {
			return nil
		}
		return fn(filepath.Dir(path))
	})
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

func init() {
	tools.Add(t...)
}
