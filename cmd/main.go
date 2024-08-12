package main

import (
	"context"
	"fmt"

	"github.com/outofforest/build"
)

func main() {
	build.Main("sei", map[string]build.Command{
		"hello": {
			Description: "Hello world",
			Fn: func(ctx context.Context, deps build.DepsFunc) error {
				fmt.Println("Hello world")
				return nil
			},
		},
	})
}
