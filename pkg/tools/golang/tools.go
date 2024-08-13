package golang

import (
	"context"

	"github.com/outofforest/build"

	"github.com/sei-protocol/build/pkg/tools"
)

// Tool names.
const (
	Go tools.Name = "go"
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
}

// EnsureGo ensures that go is available.
func EnsureGo(ctx context.Context, deps build.DepsFunc) error {
	return tools.Ensure(ctx, Go, tools.PlatformLocal)
}

func init() {
	tools.Add(t...)
}
