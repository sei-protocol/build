package docker

import (
	"context"
	"os/exec"

	"github.com/outofforest/build"
	"github.com/pkg/errors"

	"github.com/sei-protocol/build/pkg/helpers"
)

// AlpineVersion is the version of the alpine docker image.
const AlpineVersion = "3.20"

// Label used to tag docker resources created by localnet.
const (
	LabelKey   = "io.seinetwork.build"
	LabelValue = "build"
)

// EnsureDocker verifies that docker is installed.
func EnsureDocker(_ context.Context, _ build.DepsFunc) error {
	if _, err := exec.LookPath("docker"); err != nil {
		return errors.Wrap(err, "docker command is not available in PATH")
	}
	return nil
}

// Cmd returns docker command.
func Cmd(args ...string) *exec.Cmd {
	return helpers.ToolCmd("docker", args)
}
