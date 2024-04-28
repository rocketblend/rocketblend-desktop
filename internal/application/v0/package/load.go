package pack

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
	"github.com/rocketblend/rocketblend-desktop/internal/helpers"
	rbhelpers "github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/reference"
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

	return &types.Package{
		ID:   id,
		Type: definition.Type,
		//State:            state,
		Name:             extractPackageName(reference),
		Tag:              extractPackageTag(reference),
		Author:           extractPackageAuthor(reference),
		Reference:        reference,
		Path:             path,
		InstallationPath: filepath.Join(config.InstallationsPath, reference.String()),
		Platform:         source.Platform,
		URI:              source.URI,
		Verified:         isPackageVerified(reference),
		Version:          definition.Version,
		UpdatedAt:        modTime,
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
