package testrust

import (
	"path/filepath"

	"github.com/sei-protocol/build/pkg/localnet"
	"github.com/sei-protocol/build/pkg/localnet/infra"
	"github.com/sei-protocol/build/pkg/tools"
)

// Config stores rust app config.
type Config struct {
	Name string
}

// New creates new rust app.
func New(config Config) *infra.App {
	binVolume, _, dockerBinDir := localnet.BinDirMount(tools.PlatformDocker)

	return &infra.App{
		RunAsUser: true,
		Name:      config.Name,
		Volumes: []infra.Volume{
			binVolume,
		},
		Entrypoint: filepath.Join(dockerBinDir, "test-rust"),
	}
}
