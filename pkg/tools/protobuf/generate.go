package protobuf

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/outofforest/build"
	"github.com/outofforest/libexec"
	"github.com/pkg/errors"

	"github.com/sei-protocol/build/pkg/tools"
)

// GenerateGo generates go code from protobufs.
func GenerateGo(ctx context.Context, deps build.DepsFunc, protoDir, outDir string) error {
	deps(EnsureProtoc, EnsureProtocGenGo)

	protoFiles, err := findProtoFiles(protoDir)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, 0o700); err != nil {
		return errors.WithStack(err)
	}

	cmd := exec.Command(tools.Bin(ctx, "bin/protoc", tools.PlatformLocal),
		append([]string{
			"--proto_path", protoDir,
			"--plugin", tools.Bin(ctx, "bin/protoc-gen-go", tools.PlatformLocal),
			"--go_out", outDir,
		}, protoFiles...)...)

	return libexec.Exec(ctx, cmd)
}

// GenerateGoGRPC generates go GRPC service from protobufs.
func GenerateGoGRPC(ctx context.Context, deps build.DepsFunc, protoDir, outDir string) error {
	deps(EnsureProtoc, EnsureProtocGenGoGRPC)

	protoFiles, err := findProtoFiles(protoDir)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, 0o700); err != nil {
		return errors.WithStack(err)
	}

	cmd := exec.Command(tools.Bin(ctx, "bin/protoc", tools.PlatformLocal),
		append([]string{
			"--proto_path", protoDir,
			"--plugin", tools.Bin(ctx, "bin/protoc-gen-go-grpc", tools.PlatformLocal),
			"--go-grpc_out", outDir,
		}, protoFiles...)...)

	return libexec.Exec(ctx, cmd)
}

func findProtoFiles(dir string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".proto") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return files, nil
}
