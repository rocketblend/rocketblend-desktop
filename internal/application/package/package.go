package pack

import (
	"fmt"
	"path"
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

func Load(packageRootPath string, installationRootPath string, packagePath string) (*Package, error) {
	pack, err := rocketpack.Load(packagePath)
	if err != nil {
		return nil, fmt.Errorf("error loading package: %w", err)
	}

	reference, err := filePathToReference(packageRootPath, packagePath)
	if err != nil {
		return nil, fmt.Errorf("error getting path reference for package: %w", err)
	}

	modTime, err := util.GetModTime(packagePath)
	if err != nil {
		return nil, fmt.Errorf("error getting package modification time: %w", err)
	}

	packType := Unknown
	name := filepath.Base(packagePath)
	version := semver.Version{}
	sources := make(rocketpack.Sources)

	if pack.IsAddon() {
		packType = Addon
		name = pack.Addon.Name

		if pack.Addon.Version != nil {
			version = *pack.Addon.Version
		}

		if pack.Addon.Source != nil {
			sources[runtime.Undefined] = &rocketpack.Source{
				Resource: pack.Addon.Source.Resource,
				URI:      pack.Addon.Source.URI,
			}
		}
	}

	if pack.IsBuild() {
		packType = Build
		name = reference.String()

		if pack.Build.Version != nil {
			version = *pack.Build.Version
		}

		if pack.Build.Sources != nil {
			sources = pack.Build.Sources
		}
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
		Version:          version,
		UpdatedAt:        modTime,
	}, nil
}

func filePathToReference(packageRootPath string, filePath string) (reference.Reference, error) {
	strippedFilePath := stripPathToFolder(filePath, filepath.Base(packageRootPath))
	return reference.Parse(path.Dir(path.Clean(strings.TrimPrefix(strippedFilePath, "/"))))
}

func stripPathToFolder(path, folderName string) string {
	normPath := filepath.ToSlash(strings.ToLower(path))
	normFolderName := strings.ToLower(folderName)

	index := strings.Index(normPath, normFolderName)
	if index == -1 {
		return normPath
	}

	return normPath[index+len(normFolderName):]
}
