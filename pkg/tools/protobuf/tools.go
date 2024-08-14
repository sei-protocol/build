package protobuf

import (
	"context"

	"github.com/outofforest/build"

	"github.com/sei-protocol/build/pkg/tools"
	"github.com/sei-protocol/build/pkg/tools/golang"
)

// Tool names.
const (
	Protoc          tools.Name = "protoc"
	ProtocGenGo     tools.Name = "protoc-gen-go"
	ProtocGenGoGRPC tools.Name = "protoc-gen-go-grpc"
)

var t = []tools.Tool{
	// https://github.com/protocolbuffers/protobuf/releases
	tools.BinaryTool{
		Name:    Protoc,
		Version: "v25.0",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://github.com/protocolbuffers/protobuf/releases/download/v25.0/protoc-25.0-linux-x86_64.zip",
				Hash: "sha256:d26c4efe0eae3066bb560625b33b8fc427f55bd35b16f246b7932dc851554e67",
				Links: map[string]string{
					"bin/protoc": "bin/protoc",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://github.com/protocolbuffers/protobuf/releases/download/v25.0/protoc-25.0-osx-x86_64.zip",
				Hash: "sha256:15eefb30ba913e8dc4dd21d2ccb34ce04a2b33124f7d9460e5fd815a5d6459e3",
				Links: map[string]string{
					"bin/protoc": "bin/protoc",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://github.com/protocolbuffers/protobuf/releases/download/v25.0/protoc-25.0-osx-aarch_64.zip",
				Hash: "sha256:76a997df5dacc0608e880a8e9069acaec961828a47bde16c06116ed2e570588b",
				Links: map[string]string{
					"bin/protoc": "bin/protoc",
				},
			},
		},
	},

	// https://github.com/protocolbuffers/protobuf-go
	golang.GoPackageTool{
		Name:    ProtocGenGo,
		Version: "v1.34.2",
		Package: "google.golang.org/protobuf/cmd/protoc-gen-go",
	},

	// https://github.com/grpc/grpc-go/releases
	golang.GoPackageTool{
		Name:    ProtocGenGoGRPC,
		Version: "v1.5.1",
		Package: "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
	},
}

// EnsureProtoc ensures that protoc is available.
func EnsureProtoc(ctx context.Context, _ build.DepsFunc) error {
	return tools.Ensure(ctx, Protoc, tools.PlatformLocal)
}

// EnsureProtocGenGo ensures that protoc-gen-go is available.
func EnsureProtocGenGo(ctx context.Context, _ build.DepsFunc) error {
	return tools.Ensure(ctx, ProtocGenGo, tools.PlatformLocal)
}

// EnsureProtocGenGoGRPC ensures that protoc-gen-go is available.
func EnsureProtocGenGoGRPC(ctx context.Context, _ build.DepsFunc) error {
	return tools.Ensure(ctx, ProtocGenGoGRPC, tools.PlatformLocal)
}

func init() {
	tools.Add(t...)
}
