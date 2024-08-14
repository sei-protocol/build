package main

import (
	"context"

	"github.com/outofforest/build"

	"github.com/sei-protocol/build/pkg/tools"
	"github.com/sei-protocol/build/pkg/tools/golang"
	"github.com/sei-protocol/build/pkg/tools/rust"
)

func buildGo(ctx context.Context, deps build.DepsFunc) error {
	return golang.Build(ctx, deps, golang.BuildConfig{
		Platform:      tools.PlatformLocal,
		PackagePath:   "testapps/golang",
		BinOutputPath: "bin/test-go",
	})
}

func buildRust(ctx context.Context, deps build.DepsFunc) error {
	return rust.Build(ctx, deps, rust.BuildConfig{
		Platform:      tools.PlatformLocal,
		PackagePath:   "testapps/rust",
		Binary:        "rust",
		BinOutputPath: "bin/test-rust",
	})
}
