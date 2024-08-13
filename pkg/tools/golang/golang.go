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
	"go.uber.org/zap"

	"github.com/sei-protocol/build/pkg/tools"
)

// Lint lints the code.
func Lint(ctx context.Context, deps build.DepsFunc) error {
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
