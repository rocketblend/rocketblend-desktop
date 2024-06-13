package pack

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend-desktop/internal/application/enums"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	"github.com/rocketblend/rocketblend-desktop/internal/helpers"
	rbhelpers "github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/repository"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

func load(configurator types.RBConfigurator, validator types.Validator, path string) (*types.Package, error) {
	definition, err := rbhelpers.Load[types.Definition](validator, path)
	if err != nil {
		return nil, err
	}

	config, err := configurator.Get()
	if err != nil {
		return nil, err
	}

	reference, err := convertPathToReference(config.PackagesPath, path)
	if err != nil {
		return nil, fmt.Errorf("failed to convert path to reference: %w", err)
	}

	modTime, err := helpers.GetModTime(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get mod time: %w", err)
	}

	id, err := helpers.StringToUUID(reference.String())
	if err != nil {
		return nil, fmt.Errorf("failed to convert reference to UUID: %w", err)
	}

	source := definition.Source(rbtypes.Platform(config.Platform.String()))
	installationPath := filepath.Join(config.InstallationsPath, reference.String())
	state, err := determineState(installationPath, source)
	if err != nil {
		return nil, fmt.Errorf("failed to determine state: %w", err)
	}

	var progress *types.Progress
	if state == enums.PackageStateDownloading || state == enums.PackageStateIncomplete {
		progressFilePath := filepath.Join(installationPath, repository.DownloadProgressFileName)
		progress, err = loadDownloadProgress(progressFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load download progress: %w", err)
		}
	}

	platform := rbtypes.PlatformAny
	if source != nil {
		platform = source.Platform
	}

	var uri *rbtypes.URI
	if source != nil {
		uri = source.URI
	}

	return &types.Package{
		ID:               id,
		Type:             enums.PackageType(definition.Type),
		State:            state,
		Name:             extractPackageName(reference),
		Tag:              extractPackageTag(reference),
		Author:           extractPackageAuthor(reference),
		Reference:        reference,
		Path:             path,
		InstallationPath: installationPath,
		Platform:         platform,
		URI:              uri,
		Verified:         isPackageVerified(reference),
		Version:          definition.Version,
		Progress:         progress,
		UpdatedAt:        modTime,
	}, nil
}

func loadDownloadProgress(path string) (*types.Progress, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result rbtypes.Progress
	if err := json.Unmarshal(f, &result); err != nil {
		return nil, err
	}

	return &types.Progress{
		CurrentBytes:   result.Current,
		TotalBytes:     result.Total,
		BytesPerSecond: result.Speed,
	}, nil
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

func convertPathToReference(packageRootPath string, filePath string) (reference.Reference, error) {
	strippedFilePath := trimPathFromFolder(filePath, filepath.Base(packageRootPath))
	return reference.Parse(path.Dir(path.Clean(strings.TrimPrefix(strippedFilePath, "/"))))
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
