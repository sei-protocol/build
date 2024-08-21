package localnet

import (
	"context"
	"path/filepath"

	"github.com/samber/lo"

	"github.com/sei-protocol/build/pkg/localnet/infra"
	"github.com/sei-protocol/build/pkg/tools"
)

const (
	dockerBinDir      = "/usr/local/localnet/bin"
	dockerAppDir      = "/app"
	dockerLocalnetDir = "/usr/local/localnet/"
)

// BinDirMount returns mount info for repo bin directory.
func BinDirMount(platform tools.Platform) (volume infra.Volume, hostDir, dockerDir string) {
	hostDir = lo.Must(filepath.EvalSymlinks(lo.Must(filepath.Abs(filepath.Join("bin", platform.String())))))
	return infra.Volume{
		Source:      hostDir,
		Destination: dockerBinDir,
	}, hostDir, dockerBinDir
}

// AppDirMount returns mount info for application data directory.
func AppDirMount(ctx context.Context, appName string) (volume infra.Volume, hostDir, dockerDir string) {
	hostDir = appDir(ctx, appName)
	return infra.Volume{
		Source:      hostDir,
		Destination: dockerAppDir,
	}, hostDir, dockerAppDir
}

// PlatformDirMount returns mount info for application data directory.
func PlatformDirMount(ctx context.Context, platform tools.Platform) (volume infra.Volume, hostDir, dockerDir string) {
	hostPlatformDir := tools.PlatformDir(ctx, platform)
	dockerPlatformDir := filepath.Join(dockerLocalnetDir, platform.String())
	envVersion := tools.EnvVersion()
	return infra.Volume{
		Source:      hostPlatformDir,
		Destination: dockerPlatformDir,
	}, filepath.Join(hostPlatformDir, envVersion), filepath.Join(dockerPlatformDir, envVersion)
}

func appDir(ctx context.Context, appName string) string {
	return filepath.Join(rootDir(ctx), appName)
}

func rootDir(ctx context.Context) string {
	return filepath.Join(tools.EnvDir(ctx), "localnet")
}
