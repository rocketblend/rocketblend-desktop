package application

import (
	"path/filepath"

	"github.com/rocketblend/rocketblend-desktop/internal/helpers"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

	OpenExplorerOptions struct {
		Path string `json:"path"`
	}
)

func (d *Driver) OpenDirectoryDialog(opts OpenDialogOptions) (string, error) {
	path, err := runtime.OpenDirectoryDialog(d.ctx, convertOpenDialogOptions(opts))
	if err != nil {
		return "", err
	}

	return path, nil
}

func (d *Driver) OpenFileDialog(opts OpenDialogOptions) (string, error) {
	path, err := runtime.OpenFileDialog(d.ctx, convertOpenDialogOptions(opts))
	if err != nil {
		return "", err
	}

	return path, nil
}

func (d *Driver) OpenExplorer(opts OpenExplorerOptions) error {
	if err := helpers.Explore(d.ctx, filepath.Dir(opts.Path)); err != nil {
		return err
	}

	return nil
}

func convertOpenDialogOptions(opts OpenDialogOptions) runtime.OpenDialogOptions {
	return runtime.OpenDialogOptions{
		DefaultDirectory: opts.DefaultDirectory,
		DefaultFilename:  opts.DefaultFilename,
		Title:            opts.Title,
		Filters:          convertFileFilters(opts.Filters),
	}
}

func convertFileFilters(filters []FileFilter) []runtime.FileFilter {
	result := make([]runtime.FileFilter, 0, len(filters))
	for _, filter := range filters {
		result = append(result, runtime.FileFilter{
			DisplayName: filter.DisplayName,
			Pattern:     filter.Pattern,
		})
	}

	return result
}
