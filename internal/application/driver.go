package application

import (
	"context"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/packageservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/config"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/factory"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Driver struct
type Driver struct {
	ctx context.Context

	logger logger.Logger

	configService *config.Service

	projectService projectservice.Service
	packageService packageservice.Service
}

// NewApp creates a new App application struct
func NewDriver() (*Driver, error) {
	logger := logger.New(
		logger.WithLogLevel("debug"),
	)

	factory, err := factory.New()
	if err != nil {
		return nil, err
	}

	if err := factory.SetLogger(logger); err != nil {
		return nil, err
	}

	storeService, err := searchstore.New(
		searchstore.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	configService, err := factory.GetConfigService()
	if err != nil {
		return nil, err
	}

	projectservice, err := projectservice.New(
		projectservice.WithLogger(logger),
		projectservice.WithFactory(factory),
		projectservice.WithStore(storeService),
	)
	if err != nil {
		return nil, err
	}

	packageservice, err := packageservice.New(
		packageservice.WithLogger(logger),
		packageservice.WithStore(storeService),
		packageservice.WithConfig(configService),
	)
	if err != nil {
		return nil, err
	}

	return &Driver{
		logger: logger,

		configService: configService,

		projectService: projectservice,
		packageService: packageservice,
	}, nil
}

// GetProject gets a project by id
func (d *Driver) GetProject(id uuid.UUID) *projectservice.GetProjectResponse {
	ctx := context.Background()

	project, err := d.projectService.Get(ctx, id)
	if err != nil {
		d.logger.Error("Failed to find project by id", map[string]interface{}{"error": err.Error()})
		return nil
	}

	return project
}

// ListProjects lists all projects
func (d *Driver) ListProjects(query string) *projectservice.ListProjectsResponse {
	ctx := context.Background()

	response, err := d.projectService.List(ctx, listoption.WithQuery(query))
	if err != nil {
		d.logger.Error("Failed to find all projects", map[string]interface{}{"error": err.Error()})
		return nil
	}

	d.logger.Debug("Found projects", map[string]interface{}{"projects": len(response.Projects)})

	return response
}

// CreateProject creates a new project
func (d *Driver) CreateProject(request *projectservice.CreateProjectRequest) {
	ctx := context.Background()

	if err := d.projectService.Create(ctx, request); err != nil {
		d.logger.Error("Failed to create project", map[string]interface{}{"error": err.Error()})
		return
	}

	d.logger.Debug("Project created")
}

// UpdateProject updates a project
func (d *Driver) UpdateProject(request *projectservice.UpdateProjectRequest) {
	ctx := context.Background()

	if err := d.projectService.Update(ctx, request); err != nil {
		d.logger.Error("Failed to update project", map[string]interface{}{"error": err.Error()})
		return
	}

	d.logger.Debug("Project updated", map[string]interface{}{"id": request.ID})
}

// DeleteProject deletes a project
func (d *Driver) DeleteProject(id uuid.UUID) {
	d.logger.Debug("Deleting project", map[string]interface{}{"id": id})
}

// RunProject runs a project
func (d *Driver) RunProject(id uuid.UUID) {
	ctx := context.Background()

	if err := d.projectService.Run(ctx, id); err != nil {
		d.logger.Error("Failed to run project", map[string]interface{}{"error": err.Error()})
		return
	}

	d.logger.Debug("Project started", map[string]interface{}{"id": id})
}

// RenderProject renders a project
func (d *Driver) RenderProject(id uuid.UUID) {
	ctx := context.Background()

	if err := d.projectService.Render(ctx, id); err != nil {
		d.logger.Error("Failed to render project", map[string]interface{}{"error": err.Error()})
		return
	}

	d.logger.Debug("Project rendered", map[string]interface{}{"id": id})
}

// ExploreProject explores a project
func (d *Driver) ExploreProject(id uuid.UUID) {
	ctx := context.Background()

	if err := d.projectService.Explore(ctx, id); err != nil {
		d.logger.Error("Failed to explore project", map[string]interface{}{"error": err.Error()})
		return
	}

	d.logger.Debug("Project explored", map[string]interface{}{"id": id})
}

func (d *Driver) GetPackage(id uuid.UUID) *packageservice.GetPackageResponse {
	ctx := context.Background()

	pack, err := d.packageService.Get(ctx, id)
	if err != nil {
		d.logger.Error("Failed to find package by id", map[string]interface{}{"error": err.Error()})
		return nil
	}

	return pack
}

func (d *Driver) ListPackages(query string) *packageservice.ListPackagesResponse {
	ctx := context.Background()

	response, err := d.packageService.List(ctx, listoption.WithQuery(query))
	if err != nil {
		d.logger.Error("Failed to find all packages", map[string]interface{}{"error": err.Error()})
		return nil
	}

	d.logger.Debug("Found packages", map[string]interface{}{"packages": len(response.Packages)})

	return response
}

// Quit quits the application
func (d *Driver) Quit() {
	if d.projectService != nil {
		d.projectService.Close()
	}

	if d.packageService != nil {
		d.packageService.Close()
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
