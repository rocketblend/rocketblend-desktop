package pack

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend-desktop/internal/application/enums"
	"github.com/rocketblend/rocketblend/pkg/downloader"
	"github.com/rocketblend/rocketblend/pkg/repository"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

func determineState(installationPath string, source *rbtypes.Source) (enums.PackageState, error) {
	if source == nil {
		return enums.PackageStateInstalled, nil
	}

	installed, err := isInstalled(installationPath, source.Resource)
	if err != nil {
		return "", err
	}

	if installed {
		return enums.PackageStateInstalled, nil
	}

	partial, err := isPartial(installationPath, source.URI)
	if err != nil {
		return "", err
	}

	if partial {
		return enums.PackageStateDownloading, nil
	}

	active, err := isActive(installationPath)
	if err != nil {
		return "", err
	}

	if active {
		return enums.PackageStateDownloading, nil
	}

	return enums.PackageStateAvailable, nil
}

func isInstalled(installationPath string, resource string) (bool, error) {
	resourcePath := filepath.Join(installationPath, resource)
	installed, err := checkFileExistence(resourcePath)
	if err != nil {
		return false, err
	}

	return installed, nil
}

func isPartial(installationPath string, uri *rbtypes.URI) (bool, error) {
	if uri == nil {
		return false, errors.New("cannot check partial state without URI")
	}

	partialResourcePath := filepath.Join(installationPath, filepath.Base(uri.Path)+downloader.TempFileExtension)
	partial, err := checkFileExistence(partialResourcePath)
	if err != nil {
		return false, err
	}

	return partial, nil
}

func isActive(installationPath string) (bool, error) {
	lockFilePath := filepath.Join(installationPath, repository.LockFileName)
	locked, err := checkFileExistence(lockFilePath)
	if err != nil {
		return false, err
	}

	return locked, nil
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
