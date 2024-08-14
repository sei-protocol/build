package main

import (
	"context"

	"github.com/outofforest/build"

	"github.com/sei-protocol/build/pkg/tools/golang"
)

var commands = map[string]build.Command{
	"lint": {
		Description: "Lints code",
		Fn: func(_ context.Context, deps build.DepsFunc) error {
			deps(golang.Lint)
			return nil
		},
	},
	"test": {
		Description: "Runs unit tests",
		Fn: func(_ context.Context, deps build.DepsFunc) error {
			deps(golang.UnitTests)
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
