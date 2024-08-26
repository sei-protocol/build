package main

import (
	"github.com/outofforest/build"

	me "github.com/sei-protocol/build"
	_ "github.com/sei-protocol/build/pkg/tools/evm"
	"github.com/sei-protocol/build/pkg/tools/git"
	"github.com/sei-protocol/build/pkg/tools/golang"
	_ "github.com/sei-protocol/build/pkg/tools/protobuf"
	"github.com/sei-protocol/build/pkg/tools/rust"
)

func main() {
	build.RegisterCommands(
		me.Commands,
		git.Commands,
		golang.Commands,
		rust.Commands,
		commands,
	)
	build.Main("sei")
}
