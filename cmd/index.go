package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/outofforest/build"
	"github.com/outofforest/libexec"
	"github.com/samber/lo"

	"github.com/sei-protocol/build/pkg/tools"
	"github.com/sei-protocol/build/pkg/tools/golang"
)

var commands = map[string]build.Command{
	"enter": {
		Description: "Enters the environment",
		Fn:          enter,
	},
	"build/me": {
		Description: "Rebuilds the builder",
		Fn: func(ctx context.Context, deps build.DepsFunc) error {
			return golang.Build(ctx, deps, golang.BuildConfig{
				Platform:      tools.PlatformLocal,
				PackagePath:   "cmd",
				BinOutputPath: filepath.Join("bin", ".cache", filepath.Base(lo.Must(os.Executable()))),
			})
		},
	},
	"setup": {
		Description: "Installs all the tools for the host operating system",
		Fn:          tools.EnsureAll,
	},
	"lint": {
		Description: "Lints code",
		Fn:          golang.Lint,
	},
	"test": {
		Description: "Runs unit tests",
		Fn:          golang.UnitTests,
	},
	"tidy": {
		Description: "Tidies up the code",
		Fn:          golang.Tidy,
	},
}

func enter(ctx context.Context, deps build.DepsFunc) error {
	bash := exec.Command("bash")
	bash.Env = append(os.Environ(),
		fmt.Sprintf("PS1=%s", "("+build.GetName(ctx)+`) [\u@\h \W]\$ `),
		fmt.Sprintf("PATH=%s:%s", filepath.Join(tools.VersionDir(ctx, tools.PlatformLocal), "bin"), os.Getenv("PATH")),
	)
	bash.Stdin = os.Stdin
	bash.Stdout = os.Stdout
	bash.Stderr = os.Stderr
	err := libexec.Exec(ctx, bash)
	if bash.ProcessState != nil && bash.ProcessState.ExitCode() != 0 {
		return nil
	}
	return err
}
