package rust

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/outofforest/build"
	"github.com/outofforest/libexec"
	"github.com/outofforest/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/sei-protocol/build/pkg/tools"
)

// Tool names.
const (
	RustUpInit tools.Name = "rustup-init"
	Rust       tools.Name = "rust"
)

var t = []tools.Tool{
	// https://rust-lang.github.io/rustup/installation/other.html
	tools.BinaryTool{
		Name: RustUpInit,
		// update GCP bin source when update the version
		Version: "1.27.1",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://static.rust-lang.org/rustup/dist/x86_64-unknown-linux-gnu/rustup-init", //nolint:lll // breaking down urls is not beneficial
				Hash: "sha256:6aeece6993e902708983b209d04c0d1dbb14ebb405ddb87def578d41f920f56d",
				Links: map[string]string{
					"bin/rustup-init": "rustup-init",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://static.rust-lang.org/rustup/dist/x86_64-apple-darwin/rustup-init", //nolint:lll // breaking down urls is not beneficial
				Hash: "sha256:f547d77c32d50d82b8228899b936bf2b3c72ce0a70fb3b364e7fba8891eba781",
				Links: map[string]string{
					"bin/rustup-init": "rustup-init",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://static.rust-lang.org/rustup/dist/aarch64-apple-darwin/rustup-init", //nolint:lll // breaking down urls is not beneficial
				Hash: "sha256:760b18611021deee1a859c345d17200e0087d47f68dfe58278c57abe3a0d3dd0",
				Links: map[string]string{
					"bin/rustup-init": "rustup-init",
				},
			},
		},
	},

	// https://releases.rs
	RustInstaller{
		Version: "1.80.1",
	},
}

// RustInstaller installs rust.
//
//nolint:revive
type RustInstaller struct {
	Version string
}

// GetName returns the name of the tool.
func (ri RustInstaller) GetName() tools.Name {
	return Rust
}

// GetVersion returns the version of the tool.
func (ri RustInstaller) GetVersion() string {
	return ri.Version
}

// IsCompatible tells if tool is defined for the platform.
func (ri RustInstaller) IsCompatible(platform tools.Platform) (bool, error) {
	rustUpInit, err := tools.Get(RustUpInit)
	if err != nil {
		return false, err
	}
	return rustUpInit.IsCompatible(platform)
}

// Verify verifies the cheksums.
func (ri RustInstaller) Verify(ctx context.Context) ([]error, error) {
	return nil, nil
}

// Ensure ensures that tool is installed.
func (ri RustInstaller) Ensure(ctx context.Context, platform tools.Platform) error {
	binaries := ri.binaries()

	toolchain, err := ri.toolchain(ctx, platform)
	if err != nil {
		return err
	}

	install := toolchain == ""
	if !install {
		toolchainDir := filepath.Join(
			"rustup",
			"toolchains",
			toolchain,
		)

		for _, binary := range binaries {
			if tools.ShouldReinstall(ctx, platform, ri, binary, filepath.Join(toolchainDir, binary)) {
				install = true
				break
			}
		}
	}

	if install {
		if err := ri.install(ctx, platform); err != nil {
			return err
		}
	}

	return tools.LinkFiles(ctx, platform, ri, binaries)
}

func (ri RustInstaller) binaries() []string {
	return []string{
		"bin/rustc",
		"bin/cargo",
		"bin/cargo-clippy",
	}
}

func (ri RustInstaller) install(ctx context.Context, platform tools.Platform) (retErr error) {
	if err := tools.Ensure(ctx, RustUpInit, platform); err != nil {
		return errors.Wrapf(err, "ensuring rustup-installer failed")
	}

	log := logger.Get(ctx)
	log.Info("Installing binaries")

	downloadDir := tools.ToolDownloadDir(ctx, platform, ri)
	rustupHome := filepath.Join(downloadDir, "rustup")
	toolchainsDir := filepath.Join(rustupHome, "toolchains")
	cargoHome := filepath.Join(downloadDir, "cargo")
	rustupInstaller := tools.Bin(ctx, "bin/rustup-init", platform)
	rustup := filepath.Join(cargoHome, "bin", "rustup")
	env := append(
		os.Environ(),
		"RUSTUP_HOME="+rustupHome,
		"CARGO_HOME="+cargoHome,
	)

	cmdRustupInstaller := exec.Command(rustupInstaller,
		"-y",
		"--no-update-default-toolchain",
		"--no-modify-path",
	)
	cmdRustupInstaller.Env = env

	cmdRustDefault := exec.Command(rustup, "default", ri.Version)
	cmdRustDefault.Env = env

	if err := libexec.Exec(ctx, cmdRustupInstaller, cmdRustDefault); err != nil {
		return err
	}

	toolchain, err := ri.toolchain(ctx, platform)
	if err != nil {
		return err
	}

	toolchainDir := filepath.Join(toolchainsDir, toolchain)
	linksDir := tools.ToolLinksDir(ctx, platform, ri)
	for _, binary := range ri.binaries() {
		binChecksum, err := tools.Checksum(filepath.Join(toolchainDir, binary))
		if err != nil {
			return err
		}

		dstPath := filepath.Join(linksDir, binary)
		dstPathChecksum := dstPath + ":" + binChecksum
		if err := os.Remove(dstPath); err != nil && !os.IsNotExist(err) {
			return errors.WithStack(err)
		}
		if err := os.Remove(dstPathChecksum); err != nil && !os.IsNotExist(err) {
			return errors.WithStack(err)
		}

		if err := os.MkdirAll(filepath.Dir(dstPath), 0o700); err != nil {
			return errors.WithStack(err)
		}

		srcLinkPath, err := filepath.Rel(filepath.Dir(dstPathChecksum), filepath.Join(toolchainDir, binary))
		if err != nil {
			return errors.WithStack(err)
		}
		if err := os.Symlink(srcLinkPath, dstPathChecksum); err != nil {
			return errors.WithStack(err)
		}
		if err := os.Symlink(filepath.Base(dstPathChecksum), dstPath); err != nil {
			return errors.WithStack(err)
		}

		log.Info("Binary installed to path", zap.String("path", dstPath))
	}

	log.Info("Binaries installed")

	return nil
}

func (ri RustInstaller) toolchain(ctx context.Context, platform tools.Platform) (string, error) {
	downloadDir := tools.ToolDownloadDir(ctx, platform, ri)
	toolchainsDir := filepath.Join(downloadDir, "rustup", "toolchains")

	toolchains, err := os.ReadDir(toolchainsDir)
	switch {
	case err == nil:
		for _, dir := range toolchains {
			if dir.IsDir() && strings.HasPrefix(dir.Name(), ri.Version) {
				return dir.Name(), nil
			}
		}

		return "", nil
	case os.IsNotExist(err):
		return "", nil
	default:
		return "", errors.WithStack(err)
	}
}

// EnsureRust ensures that rust is available.
func EnsureRust(ctx context.Context, _ build.DepsFunc) error {
	return tools.Ensure(ctx, Rust, tools.PlatformLocal)
}

func init() {
	tools.Add(t...)
}
