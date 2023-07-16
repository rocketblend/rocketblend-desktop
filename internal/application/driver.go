package application

import (
	"context"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectsearcher"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectwatcher"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Driver struct
type Driver struct {
	ctx context.Context

	projectSearcher projectsearcher.Searcher
	projectWatcher  projectwatcher.Watcher
}

// NewApp creates a new App application struct
func NewDriver() (*Driver, error) {
	logger := logger.New()

	projectWatcher, err := projectwatcher.New(
		projectwatcher.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	projectsearcher, err := projectsearcher.New(
		projectsearcher.WithLogger(logger),
		projectsearcher.WithWatcher(projectWatcher),
	)
	if err != nil {
		return nil, err
	}

	// TODO: Move this to a config file
	watchPaths := []string{
		"D:\\Creative\\Blender\\Projects\\Testing\\RocketBlend",
	}

	for _, path := range watchPaths {
		if err := projectWatcher.AddWatchPath(path); err != nil {
			return nil, err
		}
	}

	return &Driver{
		projectSearcher: projectsearcher,
		projectWatcher:  projectWatcher,
	}, nil
}

func (d *Driver) FindAllProjects() ([]*project.Project, error) {
	return d.projectSearcher.FindAll()
}

func (d *Driver) FindProjectByPath(projectPath string) (*project.Project, error) {
	return d.projectSearcher.FindByPath(projectPath)
}

// Quit quits the application
func (d *Driver) Quit() {
	if d.projectWatcher != nil {
		d.projectWatcher.Close()
	}

	runtime.Quit(d.ctx)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (d *Driver) startup(ctx context.Context) {
	d.ctx = ctx
}

// shutdown is called when the app is shutting down
func (d *Driver) shutdown(ctx context.Context) {}
