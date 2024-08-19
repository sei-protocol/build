package main

import (
	"context"
	"path/filepath"

	"github.com/outofforest/build"

	"github.com/sei-protocol/build/pkg/tools"
	"github.com/sei-protocol/build/pkg/tools/golang"
	"github.com/sei-protocol/build/pkg/tools/protobuf"
	"github.com/sei-protocol/build/pkg/tools/rust"
)

func buildGo(ctx context.Context, deps build.DepsFunc) error {
	deps(generateGo)

	if err := golang.Build(ctx, deps, golang.BuildConfig{
		Platform:      tools.PlatformLocal,
		PackagePath:   "testapps/golang",
		BinOutputPath: "bin/test-go",
	}); err != nil {
		return err
	}

	return golang.Build(ctx, deps, golang.BuildConfig{
		Platform:      tools.PlatformDocker,
		PackagePath:   "testapps/golang",
		BinOutputPath: filepath.Join("bin", tools.PlatformDocker.String(), "test-go"),
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

func generateGo(ctx context.Context, deps build.DepsFunc) error {
	const (
		protoDir = "testapps/golang/proto"
		outDir   = "testapps/golang/generated"
	)
	if err := protobuf.GenerateGo(ctx, deps, protoDir, outDir); err != nil {
		return err
	}

	return protobuf.GenerateGoGRPC(ctx, deps, protoDir, outDir)
}
