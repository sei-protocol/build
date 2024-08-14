package main

import (
	"context"
	"time"

	"github.com/outofforest/logger"
	"github.com/outofforest/run"
	"github.com/pkg/errors"
)

func main() {
	run.New().Run("test", func(ctx context.Context) error {
		log := logger.Get(ctx)
		for {
			select {
			case <-ctx.Done():
				return errors.WithStack(ctx.Err())
			case <-time.After(time.Second):
				log.Info("I'm running!")
			}
		}
	})
}
