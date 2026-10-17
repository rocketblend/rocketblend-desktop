package application

import (
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend-desktop/internal/helpers"
)

type (
	FileFilter struct {
		// Filter information EG: "Image Files (*.jpg, *.png)"
		DisplayName string `json:"displayName"`

		// semicolon separated list of extensions, EG: "*.jpg;*.png"
		Pattern string `json:"pattern"`
	}

	OpenDialogOptions struct {
		DefaultDirectory string       `json:"defaultDirectory,omitempty" `
		DefaultFilename  string       `json:"defaultFilename,omitempty"`
		Title            string       `json:"title,omitempty"`
		Filters          []FileFilter `json:"filters,omitempty"`
	}

	SaveDialogOptions struct {
		DefaultDirectory string       `json:"defaultDirectory,omitempty"`
		DefaultFilename  string       `json:"defaultFilename,omitempty"`
		Title            string       `json:"title,omitempty"`
		Filters          []FileFilter `json:"filters,omitempty"`
	}

	OpenExplorerOptions struct {
		Path string `json:"path"`
	}
)

func (s *System) OpenDirectoryDialog(opts OpenDialogOptions) (string, error) {
	return s.ui.OpenDirectoryDialog(opts)
}

func (s *System) SaveFileDialog(opts SaveDialogOptions) (string, error) {
	return s.ui.SaveFileDialog(opts)
}

func (s *System) OpenFileDialog(opts OpenDialogOptions) (string, error) {
	return s.ui.OpenFileDialog(opts)
}

func (s *System) OpenExplorer(opts OpenExplorerOptions) error {
	path, err := determinePath(opts.Path)
	if err != nil {
		return err
	}

	if err := helpers.Explore(s.ctx, path); err != nil {
		return err
	}

	return nil
}

func determinePath(path string) (string, error) {
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return filepath.Dir(path), nil
	} else if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	return path, nil
}
