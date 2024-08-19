package localnet

import (
	"context"

	"github.com/sei-protocol/build/pkg/localnet/infra"
)

// Start starts environment.
func Start(ctx context.Context, apps ...*infra.App) (infra.AppSet, error) {
	appSet := make(infra.AppSet, 0, len(apps))
	for _, app := range apps {
		appSet = append(appSet, app)
	}

	if err := infra.NewDocker().Deploy(ctx, appSet); err != nil {
		return nil, err
	}
	return appSet, nil
}

// Stop stops environment.
func Stop(ctx context.Context, appSet infra.AppSet) error {
	return infra.NewDocker().Stop(ctx, appSet)
}

// Remove removes environment.
func Remove(ctx context.Context) error {
	return infra.NewDocker().Remove(ctx)
}
