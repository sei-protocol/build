package golang

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/outofforest/build"
	"github.com/outofforest/libexec"
	"github.com/outofforest/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/sei-protocol/build/pkg/tools"
)

// Tool names.
const (
	Go        tools.Name = "go"
	GolangCI  tools.Name = "golangci"
	LibEVMOne tools.Name = "libevmone"
)

var t = []tools.Tool{
	// https://go.dev/dl/
	tools.BinaryTool{
		Name:    Go,
		Version: "1.23.0",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://go.dev/dl/go1.23.0.linux-amd64.tar.gz",
				Hash: "sha256:905a297f19ead44780548933e0ff1a1b86e8327bb459e92f9c0012569f76f5e3",
				Links: map[string]string{
					"bin/go":    "go/bin/go",
					"bin/gofmt": "go/bin/gofmt",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://go.dev/dl/go1.23.0.darwin-amd64.tar.gz",
				Hash: "sha256:ffd070acf59f054e8691b838f274d540572db0bd09654af851e4e76ab88403dc",
				Links: map[string]string{
					"bin/go":    "go/bin/go",
					"bin/gofmt": "go/bin/gofmt",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://go.dev/dl/go1.23.0.darwin-arm64.tar.gz",
				Hash: "sha256:b770812aef17d7b2ea406588e2b97689e9557aac7e646fe76218b216e2c51406",
				Links: map[string]string{
					"bin/go":    "go/bin/go",
					"bin/gofmt": "go/bin/gofmt",
				},
			},
		},
	},

	// https://github.com/golangci/golangci-lint/releases/
	tools.BinaryTool{
		Name:    GolangCI,
		Version: "1.60.2",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://github.com/golangci/golangci-lint/releases/download/v1.60.2/golangci-lint-1.60.2-linux-amd64.tar.gz", //nolint:lll
				Hash: "sha256:607be92de8519aa88de0688e62b02ef87899386e9dfed320a04422bbe352d124",
				Links: map[string]string{
					"bin/golangci-lint": "golangci-lint-1.60.2-linux-amd64/golangci-lint",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://github.com/golangci/golangci-lint/releases/download/v1.60.2/golangci-lint-1.60.2-darwin-amd64.tar.gz", //nolint:lll // breaking down urls is not beneficial
				Hash: "sha256:875c91b7e3f00d48142920a15ba336b98a65bce65f93dd8d2f037aed81fed953",
				Links: map[string]string{
					"bin/golangci-lint": "golangci-lint-1.60.2-darwin-amd64/golangci-lint",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://github.com/golangci/golangci-lint/releases/download/v1.60.2/golangci-lint-1.60.2-darwin-arm64.tar.gz", //nolint:lll // breaking down urls is not beneficial
				Hash: "sha256:f29c40e162c704b6ca51f328f1c378f4740e50d9a27b17a975873df52cbceb72",
				Links: map[string]string{
					"bin/golangci-lint": "golangci-lint-1.60.2-darwin-arm64/golangci-lint",
				},
			},
		},
	},

	// https://github.com/ethereum/evmone/releases
	tools.BinaryTool{
		Name:    LibEVMOne,
		Version: "0.12.0",
		Sources: tools.Sources{
			tools.PlatformDockerAMD64: {
				URL:  "https://github.com/ethereum/evmone/releases/download/v0.12.0/evmone-0.12.0-linux-x86_64.tar.gz",
				Hash: "sha256:1c7b5eba0c8c3b3b2a7a05101e2d01a13a2f84b323989a29be66285dba4136ce",
				Links: map[string]string{
					"lib/libevmone.so": "lib/libevmone.so",
				},
			},
			tools.PlatformLinuxAMD64: {
				URL:  "https://github.com/ethereum/evmone/releases/download/v0.12.0/evmone-0.12.0-linux-x86_64.tar.gz",
				Hash: "sha256:1c7b5eba0c8c3b3b2a7a05101e2d01a13a2f84b323989a29be66285dba4136ce",
				Links: map[string]string{
					"lib/libevmone.so": "lib/libevmone.so",
				},
			},
			tools.PlatformDockerARM64: {
				URL:  "https://github.com/ethereum/evmone/releases/download/v0.12.0/evmone-0.12.0-linux-x86_64.tar.gz",
				Hash: "sha256:1c7b5eba0c8c3b3b2a7a05101e2d01a13a2f84b323989a29be66285dba4136ce",
				Links: map[string]string{
					"lib/libevmone.so": "lib/libevmone.so",
				},
			},
		},
	},
}

// GoPackageTool is the tool installed using go install command.
type GoPackageTool struct {
	Name    tools.Name
	Version string
	Package string
}

// GetName returns the name of the tool.
func (gpt GoPackageTool) GetName() tools.Name {
	return gpt.Name
}

// GetVersion returns the version of the tool.
func (gpt GoPackageTool) GetVersion() string {
	return gpt.Version
}

// IsCompatible tells if tool is defined for the platform.
func (gpt GoPackageTool) IsCompatible(platform tools.Platform) (bool, error) {
	golang, err := tools.Get(Go)
	if err != nil {
		return false, err
	}
	return golang.IsCompatible(platform)
}

// Verify verifies the cheksums.
func (gpt GoPackageTool) Verify(ctx context.Context) ([]error, error) {
	return nil, nil
}

// Ensure ensures that tool is installed.
func (gpt GoPackageTool) Ensure(ctx context.Context, platform tools.Platform) error {
	binName := filepath.Base(gpt.Package)
	downloadDir := tools.ToolDownloadDir(ctx, platform, gpt)
	dst := filepath.Join("bin", binName)

	//nolint:nestif // complexity comes from trivial error-handling ifs.
	if tools.ShouldReinstall(ctx, platform, gpt, dst, binName) {
		if err := tools.Ensure(ctx, Go, platform); err != nil {
			return errors.Wrapf(err, "ensuring go failed")
		}

		cmd := exec.Command(tools.Bin(ctx, "bin/go", platform), "install", gpt.Package+"@"+gpt.Version)
		cmd.Env = append(env(ctx), "GOBIN="+downloadDir)

		if err := libexec.Exec(ctx, cmd); err != nil {
			return err
		}

		srcPath := filepath.Join(downloadDir, binName)

		binChecksum, err := tools.Checksum(srcPath)
		if err != nil {
			return err
		}

		linksDir := tools.ToolLinksDir(ctx, platform, gpt)
		dstPath := filepath.Join(linksDir, dst)
		dstPathChecksum := dstPath + ":" + binChecksum

		if err := os.Remove(dstPath); err != nil && !os.IsNotExist(err) {
			panic(err)
		}
		if err := os.Remove(dstPathChecksum); err != nil && !os.IsNotExist(err) {
			return errors.WithStack(err)
		}

		if err := os.MkdirAll(filepath.Dir(dstPath), 0o700); err != nil {
			return errors.WithStack(err)
		}
		if err := os.Chmod(srcPath, 0o700); err != nil {
			return errors.WithStack(err)
		}
		srcLinkPath, err := filepath.Rel(filepath.Dir(dstPathChecksum), filepath.Join(downloadDir, binName))
		if err != nil {
			return errors.WithStack(err)
		}
		if err := os.Symlink(srcLinkPath, dstPathChecksum); err != nil {
			return errors.WithStack(err)
		}
		if err := os.Symlink(filepath.Base(dstPathChecksum), dstPath); err != nil {
			return errors.WithStack(err)
		}
		if _, err := filepath.EvalSymlinks(dstPath); err != nil {
			return errors.WithStack(err)
		}

		logger.Get(ctx).Info("Binary installed to path", zap.String("path", dstPath))
	}

	return tools.LinkFiles(ctx, platform, gpt, []string{dst})
}

// EnsureGo ensures that go is available.
func EnsureGo(ctx context.Context, _ build.DepsFunc) error {
	return tools.Ensure(ctx, Go, tools.PlatformLocal)
}

// EnsureGolangCI ensures that go linter is available.
func EnsureGolangCI(ctx context.Context, _ build.DepsFunc) error {
	return tools.Ensure(ctx, GolangCI, tools.PlatformLocal)
}

// EnsureLibEVMOne ensures that libevmone is available.
func EnsureLibEVMOne(ctx context.Context, _ build.DepsFunc) error {
	return tools.Ensure(ctx, LibEVMOne, tools.PlatformDockerARM64)
}

func init() {
	tools.Add(t...)
}
