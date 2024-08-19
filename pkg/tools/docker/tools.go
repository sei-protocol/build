package docker

import (
	"context"
	"os/exec"

	"github.com/outofforest/build"
	"github.com/pkg/errors"
)

// AlpineVersion is the version of the alpine docker image.
const AlpineVersion = "3.20"

// Label used to tag docker resources created by crust.
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
