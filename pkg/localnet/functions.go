package localnet

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"

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
	if err := infra.NewDocker().Remove(ctx); err != nil {
		return err
	}

	// It may happen that some files are flushed to disk even after processes are terminated
	// so let's try to delete dir a few times
	var err error
	for range 3 {
		if err = os.RemoveAll(rootDir(ctx)); err == nil || errors.Is(err, os.ErrNotExist) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
	return errors.WithStack(err)
}
