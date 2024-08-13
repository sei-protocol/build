package golang

import (
	"context"

	"github.com/outofforest/build"

	"github.com/sei-protocol/build/pkg/tools"
)

// Tool names.
const (
	Go       tools.Name = "go"
	GolangCI tools.Name = "golangci"
)

var t = []tools.Tool{
	tools.BinaryTool{
		Name:    Go,
		Version: "1.22.4",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://go.dev/dl/go1.22.4.linux-amd64.tar.gz",
				Hash: "sha256:ba79d4526102575196273416239cca418a651e049c2b099f3159db85e7bade7d",
				Links: map[string]string{
					"bin/go":    "go/bin/go",
					"bin/gofmt": "go/bin/gofmt",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://go.dev/dl/go1.22.4.darwin-amd64.tar.gz",
				Hash: "sha256:7788f40f3a46f201df1dc46ca640403eb535d5513fc33449164a90dbd229b761",
				Links: map[string]string{
					"bin/go":    "go/bin/go",
					"bin/gofmt": "go/bin/gofmt",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://go.dev/dl/go1.22.4.darwin-arm64.tar.gz",
				Hash: "sha256:4036c88faf57a6b096916f1827edcdbf5290a47cc5f59956e88cdd9b1b71088c",
				Links: map[string]string{
					"bin/go":    "go/bin/go",
					"bin/gofmt": "go/bin/gofmt",
				},
			},
		},
	},
	// https://github.com/golangci/golangci-lint/releases/
	tools.BinaryTool{
		Name:    GolangCI,
		Version: "1.55.2",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://github.com/golangci/golangci-lint/releases/download/v1.55.2/golangci-lint-1.55.2-linux-amd64.tar.gz",
				Hash: "sha256:ca21c961a33be3bc15e4292dc40c98c8dcc5463a7b6768a3afc123761630c09c",
				Links: map[string]string{
					"bin/golangci-lint": "golangci-lint-1.55.2-linux-amd64/golangci-lint",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://github.com/golangci/golangci-lint/releases/download/v1.55.2/golangci-lint-1.55.2-darwin-amd64.tar.gz", //nolint:lll // breaking down urls is not beneficial
				Hash: "sha256:632e96e6d5294fbbe7b2c410a49c8fa01c60712a0af85a567de85bcc1623ea21",
				Links: map[string]string{
					"bin/golangci-lint": "golangci-lint-1.55.2-darwin-amd64/golangci-lint",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://github.com/golangci/golangci-lint/releases/download/v1.55.2/golangci-lint-1.55.2-darwin-arm64.tar.gz", //nolint:lll // breaking down urls is not beneficial
				Hash: "sha256:234463f059249f82045824afdcdd5db5682d0593052f58f6a3039a0a1c3899f6",
				Links: map[string]string{
					"bin/golangci-lint": "golangci-lint-1.55.2-darwin-arm64/golangci-lint",
				},
			},
		},
	},
}

// EnsureGo ensures that go is available.
func EnsureGo(ctx context.Context, deps build.DepsFunc) error {
	return tools.Ensure(ctx, Go, tools.PlatformLocal)
}

// EnsureGolangCI ensures that go linter is available.
func EnsureGolangCI(ctx context.Context, deps build.DepsFunc) error {
	return tools.Ensure(ctx, GolangCI, tools.PlatformLocal)
}
