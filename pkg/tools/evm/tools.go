package evm

import (
	"context"

	"github.com/outofforest/build"

	"github.com/sei-protocol/build/pkg/tools"
)

// Tool names.
const (
	Foundry tools.Name = "foundry"
)

var t = []tools.Tool{
	// https://github.com/foundry-rs/foundry/releases/tag/nightly-2b1f8d6dd90f9790faf0528e05e60e573a7569ce
	tools.BinaryTool{
		Name:    Foundry,
		Version: "nightly-2b1f8d6dd90f9790faf0528e05e60e573a7569ce",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://github.com/foundry-rs/foundry/releases/download/nightly-2b1f8d6dd90f9790faf0528e05e60e573a7569ce/foundry_nightly_linux_amd64.tar.gz", //nolint:lll
				Hash: "sha256:2c75c62fab2a521938fd2a5eec6e97f9718eb0b6802852f54f1d886100fe8eb0",                                                                     //nolint:lll
				Links: map[string]string{
					"bin/anvil": "anvil",
					"bin/forge": "forge",
					"bin/cast":  "cast",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://github.com/foundry-rs/foundry/releases/download/nightly-2b1f8d6dd90f9790faf0528e05e60e573a7569ce/foundry_nightly_darwin_amd64.tar.gz", //nolint:lll
				Hash: "sha256:cf853e416cf9358174bf4fcf603b5c263aed456842b9c78661c4d77654133b7a",                                                                      //nolint:lll
				Links: map[string]string{
					"bin/anvil": "anvil",
					"bin/forge": "forge",
					"bin/cast":  "cast",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://github.com/foundry-rs/foundry/releases/download/nightly-2b1f8d6dd90f9790faf0528e05e60e573a7569ce/foundry_nightly_darwin_arm64.tar.gz", //nolint:lll
				Hash: "sha256:af157f6daac33bb4b955875e777c52d7d022e6471ed2bf1cddba9869ed5707f0",                                                                      //nolint:lll
				Links: map[string]string{
					"bin/anvil": "anvil",
					"bin/forge": "forge",
					"bin/cast":  "cast",
				},
			},
		},
	},
}

// EnsureFoundry ensures that foundry is available.
func EnsureFoundry(ctx context.Context, _ build.DepsFunc) error {
	return tools.Ensure(ctx, Foundry, tools.PlatformLocal)
}

func init() {
	tools.Add(t...)
}
