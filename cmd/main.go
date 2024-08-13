package main

import (
	"context"

	"github.com/outofforest/build"

	"github.com/sei-protocol/build/pkg/tools"
	_ "github.com/sei-protocol/build/pkg/tools/golang"
)

func main() {
	build.Main("sei", map[string]build.Command{
		"build/me": {
			Description: "Rebuilds the builder",
			Fn: func(ctx context.Context, deps build.DepsFunc) error {
				// TODO (wojciech): Implement
				return nil
			},
		},
		"setup": {
			Description: "Installs all the tools for the host operating system",
			Fn:          tools.EnsureAll,
		},
	})
}
