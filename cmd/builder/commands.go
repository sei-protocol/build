package main

import (
	"context"

	"github.com/outofforest/build"

	"github.com/sei-protocol/build/pkg/tools/golang"
	"github.com/sei-protocol/build/pkg/tools/rust"
)

var commands = map[string]build.Command{
	"build": {
		Description: "Builds all binaries",
		Fn: func(_ context.Context, deps build.DepsFunc) error {
			deps(buildGo, buildRust)
			return nil
		},
	},
	"build/go": {
		Description: "Builds test go app",
		Fn:          buildGo,
	},
	"build/rust": {
		Description: "Builds test rust app",
		Fn:          buildRust,
	},
	"lint": {
		Description: "Lints code",
		Fn: func(_ context.Context, deps build.DepsFunc) error {
			deps(golang.Lint, rust.Lint)
			return nil
		},
	},
	"test": {
		Description: "Runs unit tests",
		Fn: func(_ context.Context, deps build.DepsFunc) error {
			deps(golang.UnitTests, rust.UnitTests)
			return nil
		},
	},
	"tidy": {
		Description: "Tidies up the code",
		Fn: func(_ context.Context, deps build.DepsFunc) error {
			deps(golang.Tidy)
			return nil
		},
	},
}
