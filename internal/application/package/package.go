package pack

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/util"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

const packageFolderName = "packages" // TODO: This isn't currently guaranteed to be the name of the folder!

type (
	Package struct {
		ID               uuid.UUID             `json:"id,omitempty"`
		Type             Type                  `json:"type,omitempty"`
		Reference        reference.Reference   `json:"reference,omitempty"`
		Name             string                `json:"name,omitempty"`
		Path             string                `json:"path,omitempty"`
		InstallationPath string                `json:"installationPath,omitempty"`
		Version          semver.Version        `json:"version,omitempty"`
		Dependencies     []reference.Reference `json:"addons,omitempty"`
		Sources          rocketpack.Sources
		UpdatedAt        time.Time `json:"updatedAt,omitempty"`
	}
)

func Load(packagePath string, installationRootPath string) (*Package, error) {
	pack, err := rocketpack.Load(packagePath)
	if err != nil {
		return nil, fmt.Errorf("error loading package: %w", err)
	}

	reference, err := pathToReference(packagePath)
	if err != nil {
		return nil, fmt.Errorf("error getting path reference for package: %w", err)
	}

	modTime, err := util.GetDirModTime(packagePath)
	if err != nil {
		return nil, fmt.Errorf("error getting package modification time: %w", err)
	}

	packType := Addon
	name := pack.Addon.Name
	version := pack.Addon.Version
	sources := make(rocketpack.Sources)
	sources[runtime.Undefined] = &rocketpack.Source{
		Resource: pack.Addon.Source.Resource,
		URI:      pack.Addon.Source.URI,
	}

	if pack.IsBuild() {
		packType = Build
		version = pack.Build.Version
		sources = pack.Build.Sources
		name = reference.String()
	}

	return &Package{
		ID:               uuid.New(),
		Type:             packType,
		Name:             name,
		Reference:        reference,
		Path:             packagePath,
		InstallationPath: filepath.Join(installationRootPath, reference.String()),
		Sources:          sources,
		Dependencies:     pack.GetDependencies(),
		Version:          *version,
		UpdatedAt:        modTime,
	}, nil
}

func pathToReference(path string) (reference.Reference, error) {
	strippedPath := strings.ToLower(stripPathToFolder(path, packageFolderName))
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
