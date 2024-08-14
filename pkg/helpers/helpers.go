package helpers

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// CopyFile copies the file.
func CopyFile(dst, src string, perm os.FileMode) error {
	fr, err := os.Open(src)
	if err != nil {
		return errors.WithStack(err)
	}
	defer fr.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return errors.WithStack(err)
	}

	fw, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
	if err != nil {
		return errors.WithStack(err)
	}
	defer fw.Close()

	if _, err = io.Copy(fw, fr); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// OnModule iterates over all the modules in the source code.
func OnModule(fileName string, fn func(path string) error) error {
	return filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != fileName {
			return nil
		}
		return fn(filepath.Dir(path))
	})
}
