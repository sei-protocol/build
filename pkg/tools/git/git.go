package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/outofforest/build"
	"github.com/outofforest/libexec"
	"github.com/pkg/errors"
)

// IsStatusClean checks that there are no uncommitted files in the repo.
func IsStatusClean(ctx context.Context, _ build.DepsFunc) error {
	buf := &bytes.Buffer{}
	cmd := exec.Command("git", "status", "-s")
	cmd.Stdout = buf
	if err := libexec.Exec(ctx, cmd); err != nil {
		return errors.Wrap(err, "git command failed")
	}
	if buf.Len() > 0 {
		fmt.Println(buf)
		return errors.New("repository contains uncommitted changes")
	}
	return nil
}
