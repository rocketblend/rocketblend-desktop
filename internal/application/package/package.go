package pack

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/util"
	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/repository"
	"github.com/rocketblend/rocketblend/pkg/runtime"
	"github.com/rocketblend/rocketblend/pkg/semver"
	"github.com/rocketblend/rocketblend/pkg/types"
)

const TempFileExtension = ".tmp" //TODO: Get via rocketblend downloader.

type (
	Package struct {
		ID               uuid.UUID           `json:"id,omitempty"`
		Type             PackageType         `json:"type"`
		State            PackageState        `json:"state"`
		Reference        reference.Reference `json:"reference,omitempty"`
		Name             string              `json:"name,omitempty"`
		Author           string              `json:"author,omitempty"`
		Tag              string              `json:"tag,omitempty"`
		Path             string              `json:"path,omitempty"`
		InstallationPath string              `json:"installationPath,omitempty"`
		Operations       []string            `json:"operations,omitempty"`
		Platform         runtime.Platform    `json:"platform,omitempty"`
		Source           *types.Source       `json:"source,omitempty"`
		Version          *semver.Version     `json:"version,omitempty"`
		Verified         bool                `json:"verified,omitempty"`
		UpdatedAt        time.Time           `json:"updatedAt,omitempty"`
	}
)

func Load(validator types.Validator, packageRootPath string, installationRootPath string, packagePath string, platform runtime.Platform) (*Package, error) {
	pack, err := helpers.Load[types.Package](validator, packagePath)
	if err != nil {
		return nil, err
	}

	reference, err := convertPathToReference(packageRootPath, packagePath)
	if err != nil {
		return nil, fmt.Errorf("error converting package path to reference: %w", err)
	}

	modTime, err := util.GetModTime(packagePath)
	if err != nil {
		return nil, fmt.Errorf("error getting package modification time: %w", err)
	}

	var source *types.Source = nil
	if !pack.Bundled() {
		source = pack.Source(types.Platform(platform.String()))
	}

	installationPath := filepath.Join(installationRootPath, reference.String())
	state, err := determinePackageState(installationPath, source)
	if err != nil {
		return nil, fmt.Errorf("error determining package state: %w", err)
	}

	id, err := util.StringToUUID(reference.String())
	if err != nil {
		return nil, fmt.Errorf("error generating package id: %w", err)
	}

	return &Package{
		ID:               id,
		Type:             convertPackageType(pack.Type),
		State:            state,
		Name:             extractPackageName(reference),
		Tag:              extractPackageTag(reference),
		Author:           extractPackageAuthor(reference),
		Reference:        reference,
		Path:             packagePath,
		InstallationPath: installationPath,
		Platform:         platform,
		Source:           source,
		Verified:         isPackageVerified(reference),
		Version:          pack.Version,
		UpdatedAt:        modTime,
	}, nil
}

func convertPackageType(packageType types.PackageType) PackageType {
	switch packageType {
	case types.PackageAddon:
		return Addon
	case types.PackageBuild:
		return Build
	default:
		return Unknown
	}
}

func determinePackageState(installationPath string, source *types.Source) (PackageState, error) {
	if source == nil {
		return Installed, nil
	}

	return verifyInstallationState(installationPath, source)
}

func verifyInstallationState(installationPath string, source *types.Source) (PackageState, error) {
	if source == nil {
		return 0, fmt.Errorf("error verifying installation state: source is nil")
	}

	resourcePath := filepath.Join(installationPath, source.Resource)
	if installed, err := checkFileExistence(resourcePath); err != nil {
		return 0, fmt.Errorf("error checking if package resource '%s' is installed: %w", resourcePath, err)
	} else if installed {
		return Installed, nil
	}

	return verifyPartialDownloadState(installationPath, source)
}

func verifyPartialDownloadState(installationPath string, source *types.Source) (PackageState, error) {
	if source == nil || source.URI == nil {
		return 0, fmt.Errorf("error verifying partial download state: source URI is nil")
	}

	partialResourcePath := filepath.Join(installationPath, filepath.Base(source.URI.Path)+TempFileExtension)
	if partial, err := checkFileExistence(partialResourcePath); err != nil {
		return 0, fmt.Errorf("error checking if package resource is partially downloaded: %w", err)
	} else if partial {
		return checkLockFile(installationPath)
	}

	return Available, nil
}

func checkLockFile(installationPath string) (PackageState, error) {
	lockFilePath := filepath.Join(installationPath, repository.LockFileName)
	if locked, err := checkFileExistence(lockFilePath); err != nil {
		return 0, fmt.Errorf("error checking if package is locked: %w", err)
	} else if locked {
		return Downloading, nil
	}

	return Cancelled, nil
}

func convertPathToReference(packageRootPath string, filePath string) (reference.Reference, error) {
	strippedFilePath := trimPathFromFolder(filePath, filepath.Base(packageRootPath))
	return reference.Parse(path.Dir(path.Clean(strings.TrimPrefix(strippedFilePath, "/"))))
}

func trimPathFromFolder(path, folderName string) string {
	normPath := filepath.ToSlash(strings.ToLower(path))
	normFolderName := strings.ToLower(folderName)

	index := strings.Index(normPath, normFolderName)
	if index == -1 {
		return normPath
	}

	return normPath[index+len(normFolderName):]
}

func checkFileExistence(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// TODO: Move these to reference package
func getPathSegment(ref reference.Reference, n int) string {
	parts := strings.Split(string(ref), "/")
	if len(parts) < n {
		return ""
	}

	return parts[len(parts)-n]
}

func extractPackageName(ref reference.Reference) string {
	return getPathSegment(ref, 2)
}

func extractPackageTag(ref reference.Reference) string {
	return getPathSegment(ref, 1)
}

func extractPackageAuthor(ref reference.Reference) string {
	author, err := ref.GetRepo()
	if err != nil {
		return ""
	}

	return author
}

// TODO: Move safe list into rocketblend config.
func isPackageVerified(ref reference.Reference) bool {
	repo, err := ref.GetRepo()
	if err != nil {
		return false
	}

	return repo == "github.com/rocketblend/official-library"
}
