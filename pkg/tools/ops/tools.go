package ops

import (
	"context"
	"os/exec"

	"github.com/outofforest/build"
	"github.com/outofforest/libexec"

	"github.com/sei-protocol/build/pkg/tools"
)

// Tool names.
const (
	Terraform tools.Name = "terraform"
)

var t = []tools.Tool{
	// https://developer.hashicorp.com/terraform/install
	tools.BinaryTool{
		Name:    Terraform,
		Version: "1.9.5",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://releases.hashicorp.com/terraform/1.9.5/terraform_1.9.5_linux_amd64.zip",
				Hash: "sha256:9cf727b4d6bd2d4d2908f08bd282f9e4809d6c3071c3b8ebe53558bee6dc913b",
				Links: map[string]string{
					"bin/terraform": "terraform",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://releases.hashicorp.com/terraform/1.9.5/terraform_1.9.5_darwin_amd64.zip",
				Hash: "sha256:c28945c377d04b1d237f704729258234c471c8c4f617a1303042862f708ebbc6",
				Links: map[string]string{
					"bin/terraform": "terraform",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://releases.hashicorp.com/terraform/1.9.5/terraform_1.9.5_darwin_arm64.zip",
				Hash: "sha256:b7eca5cd6f0f6644d45d8708c1b864e64a9e26c355d2c9b585faa049f640fe71",
				Links: map[string]string{
					"bin/terraform": "terraform",
				},
			},
		},
	},
}

// EnsureTerraform ensures that terraform is available.
func EnsureTerraform(ctx context.Context, _ build.DepsFunc) error {
	return tools.Ensure(ctx, Terraform, tools.PlatformLocal)
}

// TerraformApply applies changes to the deployment.
func TerraformApply(ctx context.Context, deps build.DepsFunc, path string) error {
	return runTerraform(ctx, deps, path, "apply")
}

// TerraformDestroy destroys deployment.
func TerraformDestroy(ctx context.Context, deps build.DepsFunc, path string) error {
	return runTerraform(ctx, deps, path, "destroy")
}

func runTerraform(ctx context.Context, deps build.DepsFunc, path, action string) error {
	deps(EnsureTerraform)

	cmd := exec.Command(tools.Bin(ctx, "bin/terraform", tools.PlatformLocal), action,
		"-parallelism", "100",
		"-auto-approve",
	)
	cmd.Dir = path

	return libexec.Exec(ctx, cmd)
}

func init() {
	tools.Add(t...)
}
