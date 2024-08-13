package tools

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/outofforest/build"
	"github.com/outofforest/logger"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// Name is the type used for defining tool names.
type Name string

// Platform defines platform to install tool on.
type Platform struct {
	OS   string
	Arch string
}

func (p Platform) String() string {
	return p.OS + "." + p.Arch
}

const (
	osLinux  = "linux"
	osDarwin = "darwin"
	osDocker = "docker"

	archAMD64 = "amd64"
	archARM64 = "arm64"
)

// Platform definitions.
var (
	PlatformLocal       = Platform{OS: runtime.GOOS, Arch: runtime.GOARCH}
	PlatformLinuxAMD64  = Platform{OS: osLinux, Arch: archAMD64}
	PlatformDarwinAMD64 = Platform{OS: osDarwin, Arch: archAMD64}
	PlatformDarwinARM64 = Platform{OS: osDarwin, Arch: archARM64}
	PlatformDockerAMD64 = Platform{OS: osDocker, Arch: archAMD64}
	PlatformDockerARM64 = Platform{OS: osDocker, Arch: archARM64}
)

// Tool represents a tool to be installed.
type Tool interface {
	GetName() Name
	GetVersion() string
	IsCompatible(platform Platform) bool
	Ensure(ctx context.Context, platform Platform) error
}

var toolsMap = map[Name]Tool{}

// Add adds tools to the toolset.
func Add(tools ...Tool) {
	for _, tool := range tools {
		toolsMap[tool.GetName()] = tool
	}
}

// Source represents source where tool is fetched from.
type Source struct {
	URL   string
	Hash  string
	Links map[string]string
}

// Sources is the map of sources.
type Sources map[Platform]Source

// BinaryTool is the tool having compiled binaries available on the internet.
type BinaryTool struct {
	Name    Name
	Version string
	Sources Sources
}

// GetName returns the anme of the tool.
func (bt BinaryTool) GetName() Name {
	return bt.Name
}

// GetVersion returns the version of the tool.
func (bt BinaryTool) GetVersion() string {
	return bt.Version
}

// IsCompatible checks if tool is compatible with the platform.
func (bt BinaryTool) IsCompatible(platform Platform) bool {
	_, exists := bt.Sources[platform]
	return exists
}

// Ensure ensures the tool is installed.
func (bt BinaryTool) Ensure(ctx context.Context, platform Platform) error {
	source, exists := bt.Sources[platform]
	if !exists {
		return errors.Errorf("tool %s is not configured for platform %s", bt.Name, platform)
	}

	var install bool
	for dst, src := range source.Links {
		if shouldReinstall(ctx, platform, bt, dst, src) {
			install = true
			break
		}
	}

	if install {
		if err := bt.install(ctx, platform); err != nil {
			return err
		}
	}

	return linkFiles(ctx, platform, bt, lo.Keys(lo.Assign(source.Links)))
}

func (bt BinaryTool) install(ctx context.Context, platform Platform) (retErr error) {
	source, exists := bt.Sources[platform]
	if !exists {
		panic(errors.Errorf("tool %s is not configured for platform %s", bt.Name, platform))
	}

	ctx = logger.With(ctx,
		zap.String("tool", string(bt.Name)),
		zap.String("version", bt.Version),
		zap.String("url", source.URL),
		zap.Stringer("platform", platform))
	log := logger.Get(ctx)
	log.Info("Installing binaries")

	resp, err := http.DefaultClient.Do(lo.Must(http.NewRequestWithContext(ctx, http.MethodGet, source.URL, nil)))
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	hasher, expectedChecksum := hasher(source.Hash)
	reader := io.TeeReader(resp.Body, hasher)
	downloadDir := toolDownloadDir(ctx, platform, bt)
	if err := os.RemoveAll(downloadDir); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(downloadDir, 0o700); err != nil {
		panic(err)
	}
	defer func() {
		if retErr != nil {
			lo.Must0(os.RemoveAll(downloadDir))
		}
	}()

	if err := saveFile(source.URL, reader, downloadDir); err != nil {
		return err
	}

	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
	if actualChecksum != expectedChecksum {
		return errors.Errorf("checksum does not match for tool %s, expected: %s, actual: %s, url: %s", bt.Name,
			expectedChecksum, actualChecksum, source.URL)
	}

	linksDir := toolLinksDir(ctx, platform, bt)
	for dst, src := range source.Links {
		srcPath := filepath.Join(downloadDir, src)

		binChecksum, err := checksum(srcPath)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(linksDir, dst)
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
		if err := os.Chmod(srcPath, 0o700); err != nil {
			return errors.WithStack(err)
		}
		srcLinkPath, err := filepath.Rel(filepath.Dir(dstPathChecksum),
			filepath.Join(toolDownloadDir(ctx, platform, bt), src))
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

// EnsureAll ensures all the tools.
func EnsureAll(ctx context.Context, deps build.DepsFunc) error {
	for _, tool := range toolsMap {
		if !tool.IsCompatible(PlatformLocal) {
			continue
		}
		if err := tool.Ensure(ctx, PlatformLocal); err != nil {
			return err
		}
	}
	return nil
}

// Ensure ensures tool exists for the platform.
func Ensure(ctx context.Context, toolName Name, platform Platform) error {
	tool, err := get(toolName)
	if err != nil {
		return err
	}
	return tool.Ensure(ctx, platform)
}

// VersionDir returns path to the version directory.
func VersionDir(ctx context.Context, platform Platform) string {
	return filepath.Join(platformDir(ctx, platform), envVersion())
}

// Bin returns path to the installed binary.
func Bin(ctx context.Context, binary string, platform Platform) string {
	return lo.Must(filepath.Abs(lo.Must(filepath.EvalSymlinks(
		filepath.Join(VersionDir(ctx, platform), binary)))))
}

func get(name Name) (Tool, error) {
	t, exists := toolsMap[name]
	if !exists {
		return nil, errors.Errorf("tool %s does not exist", name)
	}
	return t, nil
}

func envDir(ctx context.Context) string {
	return filepath.Join(lo.Must(os.UserCacheDir()), build.GetName(ctx))
}

func platformDir(ctx context.Context, platform Platform) string {
	return filepath.Join(envDir(ctx), platform.String())
}

func downloadsDir(ctx context.Context, platform Platform) string {
	return filepath.Join(platformDir(ctx, platform), "downloads")
}

func toolDownloadDir(ctx context.Context, platform Platform, tool Tool) string {
	return filepath.Join(downloadsDir(ctx, platform), string(tool.GetName())+"-"+tool.GetVersion())
}

func toolLinksDir(ctx context.Context, platform Platform, tool Tool) string {
	return filepath.Join(toolDownloadDir(ctx, platform, tool), "_links")
}

func envVersion() string {
	module := module()

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panic("reading build info failed")
	}

	for _, m := range append([]*debug.Module{&bi.Main}, bi.Deps...) {
		if m.Path != module {
			continue
		}
		if m.Replace != nil {
			m = m.Replace
		}

		if m.Version == "(devel)" {
			return "devel"
		}

		return m.Version
	}

	panic("impossible condition: build module not found")
}

func module() string {
	_, file, _, _ := runtime.Caller(0)
	module := strings.Join(strings.Split(file, "/")[:3], "/")
	index := strings.Index(module, "@")
	if index > 0 {
		module = module[:index]
	}
	return module
}

func shouldReinstall(ctx context.Context, platform Platform, tool Tool, dst, src string) bool {
	srcAbsPath, err := filepath.Abs(filepath.Join(toolDownloadDir(ctx, platform, tool), src))
	if err != nil {
		return true
	}

	srcRealPath, err := filepath.EvalSymlinks(srcAbsPath)
	if err != nil {
		return true
	}

	dstAbsPath, err := filepath.Abs(filepath.Join(toolLinksDir(ctx, platform, tool), dst))
	if err != nil {
		return true
	}

	dstRealPath, err := filepath.EvalSymlinks(dstAbsPath)
	if err != nil || dstRealPath != srcRealPath {
		return true
	}

	fInfo, err := os.Stat(dstRealPath)
	if err != nil {
		return true
	}
	if fInfo.Mode()&0o700 == 0 {
		return true
	}

	linkedPath, err := os.Readlink(dstAbsPath)
	if err != nil {
		return true
	}
	linkNameParts := strings.Split(filepath.Base(linkedPath), ":")
	if len(linkNameParts) < 3 {
		return true
	}

	hasher, expectedChecksum := hasher(linkNameParts[len(linkNameParts)-2] + ":" + linkNameParts[len(linkNameParts)-1])
	f, err := os.Open(dstRealPath)
	if err != nil {
		return true
	}
	defer f.Close()

	if _, err := io.Copy(hasher, f); err != nil {
		return true
	}

	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
	return actualChecksum != expectedChecksum
}

func shouldRelinkFile(ctx context.Context, platform Platform, tool Tool, dst string) (bool, error) {
	srcPath := filepath.Join(toolLinksDir(ctx, platform, tool), dst)

	realSrcPath, err := filepath.EvalSymlinks(srcPath)
	if err != nil {
		return false, errors.WithStack(err)
	}

	versionedPath := filepath.Join(VersionDir(ctx, platform), dst)
	realVersionedPath, err := filepath.EvalSymlinks(versionedPath)
	if err != nil {
		return true, nil //nolint:nilerr // this is ok
	}

	return realSrcPath != realVersionedPath, nil
}

func linkFiles(ctx context.Context, platform Platform, tool Tool, binaries []string) error {
	for _, dst := range binaries {
		relink, err := shouldRelinkFile(ctx, platform, tool, dst)
		if err != nil {
			return err
		}

		if !relink {
			continue
		}

		dstVersion := filepath.Join(VersionDir(ctx, platform), dst)
		src, err := filepath.Rel(filepath.Dir(dstVersion), filepath.Join(toolLinksDir(ctx, platform, tool), dst))
		if err != nil {
			return errors.WithStack(err)
		}

		if err := os.Remove(dstVersion); err != nil && !errors.Is(err, os.ErrNotExist) {
			return errors.WithStack(err)
		}

		if err := os.MkdirAll(filepath.Dir(dstVersion), 0o700); err != nil {
			return errors.WithStack(err)
		}

		if err := os.Symlink(src, dstVersion); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func hasher(hashStr string) (hash.Hash, string) {
	parts := strings.SplitN(hashStr, ":", 2)
	if len(parts) != 2 {
		panic(errors.Errorf("incorrect checksum format: %s", hashStr))
	}
	hashAlgorithm := parts[0]
	checksum := parts[1]

	var hasher hash.Hash
	switch hashAlgorithm {
	case "sha256":
		hasher = sha256.New()
	default:
		panic(errors.Errorf("unsupported hashing algorithm: %s", hashAlgorithm))
	}

	return hasher, strings.ToLower(checksum)
}

func saveFile(url string, reader io.Reader, path string) error {
	switch {
	case strings.HasSuffix(url, ".tar.gz") || strings.HasSuffix(url, ".tgz"):
		var err error
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return errors.WithStack(err)
		}
		return untar(reader, path)
	case strings.HasSuffix(url, ".zip"):
		return unzip(reader, path)
	default:
		f, err := os.OpenFile(filepath.Join(path, filepath.Base(url)), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o700)
		if err != nil {
			return errors.WithStack(err)
		}
		defer f.Close()
		_, err = io.Copy(f, reader)
		return errors.WithStack(err)
	}
}

func untar(reader io.Reader, path string) error {
	tr := tar.NewReader(reader)
	for {
		header, err := tr.Next()
		switch {
		case errors.Is(err, io.EOF):
			return nil
		case err != nil:
			return errors.WithStack(err)
		case header == nil:
			continue
		}
		header.Name = path + "/" + header.Name

		// We take mode from header.FileInfo().Mode(), not from header.Mode because they may be in
		// different formats (meaning of bits may be different).
		// header.FileInfo().Mode() returns compatible value.
		mode := header.FileInfo().Mode()

		switch {
		case header.Typeflag == tar.TypeDir:
			if err := os.MkdirAll(header.Name, mode); err != nil && !os.IsExist(err) {
				return errors.WithStack(err)
			}
		case header.Typeflag == tar.TypeReg:
			if err := ensureDir(header.Name); err != nil {
				return err
			}

			f, err := os.OpenFile(header.Name, os.O_CREATE|os.O_WRONLY, mode)
			if err != nil {
				return errors.WithStack(err)
			}
			_, err = io.Copy(f, tr)
			_ = f.Close()
			if err != nil {
				return errors.WithStack(err)
			}
		case header.Typeflag == tar.TypeSymlink:
			if err := ensureDir(header.Name); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, header.Name); err != nil {
				return errors.WithStack(err)
			}
		case header.Typeflag == tar.TypeLink:
			header.Linkname = path + "/" + header.Linkname
			if err := ensureDir(header.Name); err != nil {
				return err
			}
			if err := ensureDir(header.Linkname); err != nil {
				return err
			}
			// linked file may not exist yet, so let's create it - it will be overwritten later
			f, err := os.OpenFile(header.Linkname, os.O_CREATE|os.O_EXCL, mode)
			if err != nil {
				if !os.IsExist(err) {
					return errors.WithStack(err)
				}
			} else {
				_ = f.Close()
			}
			if err := os.Link(header.Linkname, header.Name); err != nil {
				return errors.WithStack(err)
			}
		default:
			return errors.Errorf("unsupported file type: %d", header.Typeflag)
		}
	}
}

func unzip(reader io.Reader, path string) error {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "zipfile")
	if err != nil {
		return errors.WithStack(err)
	}
	defer os.Remove(tempFile.Name()) //nolint: errcheck

	// Copy the contents of the reader to the temporary file
	_, err = io.Copy(tempFile, reader)
	if err != nil {
		return errors.WithStack(err)
	}

	// Open the temporary file for reading
	file, err := os.Open(tempFile.Name())
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	// Get the file information to obtain its size
	fileInfo, err := file.Stat()
	if err != nil {
		return errors.WithStack(err)
	}
	fileSize := fileInfo.Size()

	// Use the file as a ReaderAt to unpack the zip file
	zipReader, err := zip.NewReader(file, fileSize)
	if err != nil {
		return errors.WithStack(err)
	}

	// Process the files in the zip archive
	for _, zf := range zipReader.File {
		// Open each file in the archive
		rc, err := zf.Open()
		if err != nil {
			return errors.WithStack(err)
		}
		defer rc.Close()

		// Construct the destination path for the file
		destPath := filepath.Join(path, zf.Name)

		// skip empty dirs
		if zf.FileInfo().IsDir() {
			continue
		}

		err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		if err != nil {
			return errors.WithStack(err)
		}

		// Create the file in the destination path
		outputFile, err := os.Create(destPath)
		if err != nil {
			return errors.WithStack(err)
		}
		defer outputFile.Close()

		// Copy the file contents
		_, err = io.Copy(outputFile, rc)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func ensureDir(file string) error {
	if err := os.MkdirAll(filepath.Dir(file), 0o700); !os.IsExist(err) {
		return errors.WithStack(err)
	}
	return nil
}

func checksum(file string) (string, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0o600)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", errors.WithStack(err)
	}

	return "sha256:" + hex.EncodeToString(hasher.Sum(nil)), nil
}
