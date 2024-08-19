package localnet

import (
	"context"
	"path/filepath"

	"github.com/sei-protocol/build/pkg/localnet/infra"
	"github.com/sei-protocol/build/pkg/tools"
)

// DockerAppDir is in-docker directory for application-specific data.
const DockerAppDir = "/app"

// AppDir returns on-host directory for application-specific data.
func AppDir(ctx context.Context, appName string) string {
	return filepath.Join(rootDir(ctx), appName)
}

// AppDirVolume returns the volume to be mounted inside docker for application-specific data.
func AppDirVolume(ctx context.Context, appName string) infra.Volume {
	return infra.Volume{
		Source:      AppDir(ctx, appName),
		Destination: DockerAppDir,
	}
}

func rootDir(ctx context.Context) string {
	return filepath.Join(tools.EnvDir(ctx), "localnet")
}
