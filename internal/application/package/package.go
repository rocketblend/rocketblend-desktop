package pack

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/util"
	"github.com/rocketblend/rocketblend/pkg/downloader"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

const packageFolderName = "packages" // TODO: This isn't currently guaranteed to be the name of the folder!

type (
	Package struct {
		ID           uuid.UUID             `json:"id,omitempty"`
		Type         string                `json:"type,omitempty"`
		Reference    reference.Reference   `json:"reference,omitempty"`
		Name         string                `json:"name,omitempty"`
		Path         string                `json:"path,omitempty"`
		Version      *semver.Version       `json:"version,omitempty"`
		Dependencies []reference.Reference `json:"addons,omitempty"`
		Resource     string                `json:"resource,omitempty"`
		URI          *downloader.URI       `json:"uri,omitempty"`
		UpdatedAt    time.Time             `json:"updatedAt,omitempty"`
	}
)

func Load(packagePath string, platform runtime.Platform) (*Build, *Addon, error) {
	pack, err := rocketpack.Load(packagePath)
	if err != nil {
		return nil, nil, err
	}

	if pack.IsAddon() {
		addon, err := NewAddon(packagePath, pack)
		return nil, addon, err
	}

	build, err := NewBuild(packagePath, platform, pack)
	return build, nil, err
}

func LoadPack(packagePath string, platform runtime.Platform) (*Package, error) {
	pack, err := rocketpack.Load(packagePath)
	if err != nil {
		return nil, err
	}

	reference, err := pathToReference(packagePath)
	if err != nil {
		return nil, err
	}

	name, err := reference.GetRepo()
	if err != nil {
		return nil, err
	}

	modTime, err := util.GetDirModTime(packagePath)
	if err != nil {
		return nil, err
	}

	version := pack.Addon.Version
	resource := pack.Addon.Source.Resource
	uri := pack.Addon.Source.URI
	packType := "addon"

	if pack.IsBuild() {
		version = pack.Build.Version
		resource = pack.Build.Sources[platform].Resource
		uri = pack.Build.Sources[platform].URI
		packType = "build"
	}

	return &Package{
		ID:           uuid.New(),
		Type:         packType,
		Name:         name,
		Reference:    reference,
		Path:         packagePath,
		Dependencies: pack.GetDependencies(),
		Version:      version,
		Resource:     resource,
		URI:          uri,
		UpdatedAt:    modTime,
	}, nil
}

func pathToReference(path string) (reference.Reference, error) {
	strippedPath := stripPathToFolder(path, packageFolderName)
	dir := filepath.Dir(strippedPath)
	if dir == "." {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	return reference.Parse(dir)
}

func stripPathToFolder(path, folderName string) string {
	normPath := filepath.ToSlash(strings.ToLower(path))
	normFolderName := strings.ToLower(folderName)

	index := strings.Index(normPath, normFolderName)
	if index == -1 {
		return path
	}

	return path[index+len(normFolderName):]
}
