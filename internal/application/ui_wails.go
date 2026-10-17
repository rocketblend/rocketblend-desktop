package application

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type wailsUI struct {
	ctx context.Context
}

var _ UI = (*wailsUI)(nil)

func newWailsUI() *wailsUI {
	return &wailsUI{}
}

// startup must be called from OnStartup before any other method.
func (w *wailsUI) startup(ctx context.Context) {
	w.ctx = ctx
}

func (w *wailsUI) Emit(event string, data interface{}) {
	runtime.EventsEmit(w.ctx, event, data)
}

func (w *wailsUI) On(event string, handler func(data ...interface{})) func() {
	return runtime.EventsOn(w.ctx, event, handler)
}

func (w *wailsUI) Once(event string, handler func(data ...interface{})) func() {
	return runtime.EventsOnce(w.ctx, event, handler)
}

func (w *wailsUI) OpenFileDialog(opts OpenDialogOptions) (string, error) {
	return runtime.OpenFileDialog(w.ctx, convertOpenDialogOptions(opts))
}

func (w *wailsUI) OpenDirectoryDialog(opts OpenDialogOptions) (string, error) {
	return runtime.OpenDirectoryDialog(w.ctx, convertOpenDialogOptions(opts))
}

func (w *wailsUI) SaveFileDialog(opts SaveDialogOptions) (string, error) {
	return runtime.SaveFileDialog(w.ctx, convertSaveDialogOptions(opts))
}

func (w *wailsUI) ShowWindow() {
	runtime.WindowUnminimise(w.ctx)
	runtime.Show(w.ctx)
}

func (w *wailsUI) Quit() {
	runtime.Quit(w.ctx)
}

func convertSaveDialogOptions(opts SaveDialogOptions) runtime.SaveDialogOptions {
	return runtime.SaveDialogOptions{
		DefaultDirectory:     opts.DefaultDirectory,
		DefaultFilename:      opts.DefaultFilename,
		Title:                opts.Title,
		Filters:              convertFileFilters(opts.Filters),
		CanCreateDirectories: true,
	}
}

func convertOpenDialogOptions(opts OpenDialogOptions) runtime.OpenDialogOptions {
	return runtime.OpenDialogOptions{
		DefaultDirectory:     opts.DefaultDirectory,
		DefaultFilename:      opts.DefaultFilename,
		CanCreateDirectories: true,
		Title:                opts.Title,
		Filters:              convertFileFilters(opts.Filters),
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
