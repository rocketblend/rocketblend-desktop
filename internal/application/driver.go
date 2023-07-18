package application

import (
	"context"
	"fmt"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectsearcher"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectwatcher"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Driver struct
type Driver struct {
	ctx context.Context

	logger logger.Logger

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
		logger:          logger,
		projectSearcher: projectsearcher,
		projectWatcher:  projectWatcher,
	}, nil
}

// Greet returns a greeting for the given name
func (a *Driver) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (d *Driver) FindAllProjects() []*project.Project {
	return nil

	// projects, err := d.projectSearcher.FindAll()
	// if err != nil {
	// 	d.logger.Error("Failed to find all projects", map[string]interface{}{"error": err.Error()})
	// 	return nil
	// }

	// return projects
}

func (d *Driver) FindProjectByPath(projectPath string) *project.Project {
	return nil

	// project, err := d.projectSearcher.FindByPath(projectPath)
	// if err != nil {
	// 	d.logger.Error("Failed to find project by path", map[string]interface{}{"error": err.Error()})
	// 	return nil
	// }

	// return project
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
