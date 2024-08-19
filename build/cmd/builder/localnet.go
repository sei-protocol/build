package main

import (
	"context"

	"github.com/outofforest/build"

	"github.com/sei-protocol/build/localnet/testgo"
	"github.com/sei-protocol/build/localnet/testrust"
	"github.com/sei-protocol/build/pkg/localnet"
)

func localnetStart(ctx context.Context, deps build.DepsFunc) error {
	deps(buildGo, buildRust)

	_, err := localnet.Start(ctx,
		testgo.New(testgo.Config{Name: "testgo"}),
		testrust.New(testrust.Config{Name: "testrust"}),
	)
	return err
}

func localnetRemove(ctx context.Context, _ build.DepsFunc) error {
	return localnet.Remove(ctx)
}
