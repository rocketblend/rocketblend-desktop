package application

import (
	"context"
	"fmt"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectstore"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectstore/listoptions"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Driver struct
type Driver struct {
	ctx context.Context

	logger logger.Logger

	projectService projectservice.Service
	projectStore   projectstore.Store
}

// NewApp creates a new App application struct
func NewDriver() (*Driver, error) {
	logger := logger.New(
		logger.WithLogLevel("debug"),
	)

	projectStore, err := projectstore.New(
		projectstore.WithLogger(logger),
		projectstore.WithWatcher(),
		projectstore.WithDebounceDuration(2*time.Second),
		// TODO: Move this to a config file
		projectstore.WithPaths("D:\\Creative\\Blender\\Projects\\Testing\\RocketBlend"),
	)
	if err != nil {
		return nil, err
	}

	projectservice, err := projectservice.New(
		projectservice.WithLogger(logger),
		projectservice.WithStore(projectStore),
	)
	if err != nil {
		return nil, err
	}

	return &Driver{
		logger:         logger,
		projectService: projectservice,
		projectStore:   projectStore,
	}, nil
}

// Greet returns a greeting for the given name
func (a *Driver) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// FindAllProjects finds all projects
func (d *Driver) FindAllProjects(query string) []*projectservice.Project {
	projects, err := d.projectService.FindAll(listoptions.WithQuery(query))
	if err != nil {
		d.logger.Error("Failed to find all projects", map[string]interface{}{"error": err.Error()})
		return nil
	}

	d.logger.Debug("Found projects", map[string]interface{}{"projects": len(projects)})

	return projects
}

// FindProjectByID finds a project by its id
func (d *Driver) FindProjectByID(id uuid.UUID) *projectservice.Project {
	project, err := d.projectService.FindByID(id)
	if err != nil {
		d.logger.Error("Failed to find project by id", map[string]interface{}{"error": err.Error()})
		return nil
	}

	return project
}

// Quit quits the application
func (d *Driver) Quit() {
	if d.projectStore != nil {
		d.projectStore.Close()
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
