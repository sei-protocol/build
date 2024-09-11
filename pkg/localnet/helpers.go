package localnet

import (
	"context"
	"os"
	"path/filepath"

	"github.com/outofforest/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"

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
	lo.Must0(os.MkdirAll(hostDir, 0o700))

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
	log := logger.Get(ctx)
	log.Info("PlatformDirMount", zap.String("hostPlatformDir", hostPlatformDir), zap.String("dockerPlatformDir", dockerPlatformDir), zap.String("envVersion", envVersion))
	return infra.Volume{
		Source:      hostPlatformDir,
		Destination: dockerPlatformDir,
	}, filepath.Join(hostPlatformDir, envVersion), filepath.Join(dockerPlatformDir, envVersion)
}

// FromHostAddress returns address of the app service seen from the host.
func FromHostAddress(proto string, app *infra.App, port infra.PortName) string {
	return infra.JoinNetAddr(proto, app.Info().HostFromHost, app.Ports[port])
}

// FromContainerAddress returns address of the app service seen from the container.
func FromContainerAddress(proto string, app *infra.App, port infra.PortName) string {
	return infra.JoinNetAddr(proto, app.Info().HostFromContainer, app.Ports[port])
}

func appDir(ctx context.Context, appName string) string {
	return filepath.Join(rootDir(ctx), appName)
}

func rootDir(ctx context.Context) string {
	return filepath.Join(tools.EnvDir(ctx), "localnet")
}
