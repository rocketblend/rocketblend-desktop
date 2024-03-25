package application

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/buffermanager"
	"github.com/rocketblend/rocketblend-desktop/internal/application/config"
	"github.com/rocketblend/rocketblend-desktop/internal/application/factory"
	pack "github.com/rocketblend/rocketblend-desktop/internal/application/package"
	"github.com/rocketblend/rocketblend-desktop/internal/application/packageservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	rbruntime "github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	rocketblendConfig "github.com/rocketblend/rocketblend/pkg/rocketblend/config"
)

// Driver struct
type Driver struct {
	ctx    context.Context
	logger logger.Logger

	heartbeatInterval time.Duration

	factory      factory.Factory
	events       buffermanager.BufferManager
	cancelTokens sync.Map

	args []string
}

func NewDriver(factory factory.Factory, events buffermanager.BufferManager, args ...string) (*Driver, error) {
	logger, err := factory.GetLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to get logger: %w", err)
	}

	return &Driver{
		factory:           factory,
		heartbeatInterval: 5000 * time.Millisecond, // 1 second
		events:            events,
		logger:            logger,
		args:              args,
	}, nil
}

func (d *Driver) GetPlatform() (*rbruntime.Platform, error) {
	configService, err := d.factory.GetConfigService()
	if err != nil {
		d.logger.Error("failed to get config service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	config, err := configService.Get()
	if err != nil {
		d.logger.Error("failed to get config", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return &config.Platform, nil
}

func (d *Driver) GetApplicationConfig() (*config.Config, error) {
	configService, err := d.factory.GetApplicationConfigService()
	if err != nil {
		d.logger.Error("failed to get application config service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	config, err := configService.Get()
	if err != nil {
		d.logger.Error("failed to get config", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return config, nil
}

func (d *Driver) GetRocketBlendConfig() (*rocketblendConfig.Config, error) {
	configService, err := d.factory.GetConfigService()
	if err != nil {
		d.logger.Error("failed to get rocketblend config service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	config, err := configService.Get()
	if err != nil {
		d.logger.Error("failed to get config", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return config, nil
}

// GetProject gets a project by id
func (d *Driver) GetProject(id uuid.UUID) (*projectservice.GetProjectResponse, error) {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	project, err := projectService.Get(ctx, id)
	if err != nil {
		d.logger.Error("failed to find project by id", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return project, nil
}

// ListProjects lists all projects
func (d *Driver) ListProjects(query string) (*projectservice.ListProjectsResponse, error) {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	response, err := projectService.List(ctx, listoption.WithQuery(query))
	if err != nil {
		d.logger.Error("failed to find all projects", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	d.logger.Debug("found projects", map[string]interface{}{"projects": len(response.Projects)})

	return response, nil
}

// CreateProject creates a new project
func (d *Driver) CreateProject(request *projectservice.CreateProjectRequest) error {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Create(ctx, request); err != nil {
		d.logger.Error("failed to create project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("project created")
	return nil
}

// UpdateProject updates a project
func (d *Driver) UpdateProject(request *projectservice.UpdateProjectRequest) error {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Update(ctx, request); err != nil {
		d.logger.Error("failed to update project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("project updated", map[string]interface{}{"id": request.ID})

	return nil
}

// DeleteProject deletes a project
func (d *Driver) DeleteProject(id uuid.UUID) error {
	d.logger.Debug("deleting project", map[string]interface{}{"id": id})
	return nil
}

// RunProject runs a project
func (d *Driver) RunProject(id uuid.UUID) error {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Run(ctx, id); err != nil {
		d.logger.Error("failed to run project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("project started", map[string]interface{}{"id": id})
	return nil
}

// RenderProject renders a project
func (d *Driver) RenderProject(id uuid.UUID) error {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Render(ctx, id); err != nil {
		d.logger.Error("failed to render project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("project rendered", map[string]interface{}{"id": id})
	return nil
}

// ExploreProject explores a project
func (d *Driver) ExploreProject(id uuid.UUID) error {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Explore(ctx, id); err != nil {
		d.logger.Error("failed to explore project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("project explored", map[string]interface{}{"id": id})
	return nil
}

func (d *Driver) GetPackage(id uuid.UUID) (*packageservice.GetPackageResponse, error) {
	ctx := context.Background()

	packageService, err := d.factory.GetPackageService()
	if err != nil {
		d.logger.Error("failed to get package service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	pack, err := packageService.Get(ctx, id)
	if err != nil {
		d.logger.Error("failed to find package by id", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return pack, err
}

func (d *Driver) ListPackages(query string, packageType pack.PackageType, installed bool) (*packageservice.ListPackagesResponse, error) {
	ctx := context.Background()

	packageService, err := d.factory.GetPackageService()
	if err != nil {
		d.logger.Error("failed to get package service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	category := ""
	if packageType != pack.Unknown {
		category = strconv.Itoa(int(packageType))
	}

	var state *int = nil
	if installed {
		stateInt := int(pack.Installed)
		state = &stateInt
	}

	response, err := packageService.List(ctx, []listoption.ListOption{
		listoption.WithQuery(query),
		listoption.WithCategory(category),
		listoption.WithState(state),
	}...)
	if err != nil {
		d.logger.Error("failed to find all packages", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	d.logger.Debug("found packages", map[string]interface{}{"packages": len(response.Packages)})

	return response, err
}

func (d *Driver) AddPackage(referenceStr string) error {
	packageService, err := d.factory.GetPackageService()
	if err != nil {
		d.logger.Error("failed to get package service", map[string]interface{}{"error": err.Error()})
		return err
	}

	ref, err := reference.Parse(referenceStr)
	if err != nil {
		d.logger.Error("failed to parse reference", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := packageService.Add(d.ctx, ref); err != nil {
		d.logger.Error("failed to add package", map[string]interface{}{"error": err.Error()})
		return err
	}

	return nil
}

func (d *Driver) RefreshPackages() error {
	packageService, err := d.factory.GetPackageService()
	if err != nil {
		d.logger.Error("failed to get package service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := packageService.Refresh(d.ctx); err != nil {
		d.logger.Error("failed to refresh packages", map[string]interface{}{"error": err.Error()})
		return err
	}

	return nil
}

func (d *Driver) UninstallPackage(id uuid.UUID) error {
	packageService, err := d.factory.GetPackageService()
	if err != nil {
		d.logger.Error("failed to get package service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := packageService.Uninstall(d.ctx, id); err != nil {
		d.logger.Error("failed to uninstall package", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("package uninstalled", map[string]interface{}{"id": id})
	return nil
}

func (d *Driver) LongRunningRequestWithCancellation(cid uuid.UUID) error {
	_, err := d.runWithCancellation(cid, func(ctx context.Context) (interface{}, error) {
		// Simulate a long-running operation
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				d.logger.Debug("long running request canceled", map[string]interface{}{"cid": cid})
				return nil, ctx.Err()
			default:
				time.Sleep(2 * time.Second)
			}
		}

		return struct{}{}, nil
	})
	if err != nil {
		return err
	}

	d.logger.Debug("long running request completed", map[string]interface{}{"cid": cid})
	return nil
}

// Quit quits the application
func (d *Driver) Quit() {
	d.logger.Debug("quitting application")

	if err := d.addApplicationMetrics(); err != nil {
		d.logger.Error("failed to add application metrics", map[string]interface{}{"error": err.Error()})
		return
	}

	runtime.Quit(d.ctx)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (d *Driver) startup(ctx context.Context) {
	d.logger.Debug("starting application")

	d.ctx = ctx

	// Start listening to log events
	go d.listenToLogEvents()

	// Preloads all the data
	if err := d.factory.Preload(); err != nil {
		d.logger.Error("failed to preload", map[string]interface{}{"error": err.Error()})
		return
	}

	// Setup driver event handlers (backend)
	if err := d.setupDriverEventHandlers(); err != nil {
		d.logger.Error("failed to setup driver event handlers", map[string]interface{}{"error": err.Error()})
		return
	}

	// Setup runtime event handlers (frontend)
	if err := d.setupRuntimeEventHandlers(); err != nil {
		d.logger.Error("failed to setup runtime event handlers", map[string]interface{}{"error": err.Error()})
		return
	}
}

// shutdown is called when the app is shutting down
func (d *Driver) shutdown(ctx context.Context) {
	d.logger.Debug("shutting down application")

	// Close the event stream
	d.events.Close()

	// Close the factory watchers
	if err := d.factory.Close(); err != nil {
		d.logger.Error("failed to close factory", map[string]interface{}{"error": err.Error()})
	}

	d.logger.Debug("application shutdown")
}

// onDomReady is called when the DOM is ready
func (d *Driver) onDomReady(ctx context.Context) {
	d.logger.Debug("DOM is ready")

	// Wait for main layout to be ready.
	runtime.EventsOnce(ctx, "ready", func(optionalData ...interface{}) {
		d.onLayoutReady(ctx)
	})
}

// onLayoutReady is called when the layout is ready
func (d *Driver) onLayoutReady(ctx context.Context) {
	d.logger.Debug("main layout is ready")

	if err := d.addApplicationMetrics(); err != nil {
		d.logger.Error("failed to add application metrics", map[string]interface{}{"error": err.Error()})
		return
	}

	d.eventEmitLaunchArgs(ctx, LaunchEvent{
		Args: os.Args[1:],
	})
}

// onSecondInstanceLaunch is called when the user opens a second instance of the application
func (d *Driver) onSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
	secondInstanceArgs := secondInstanceData.Args

	d.logger.Info("user opened second instance", map[string]interface{}{
		"args":             strings.Join(secondInstanceData.Args, ","),
		"workingDirectory": secondInstanceData.WorkingDirectory,
	})

	runtime.WindowUnminimise(d.ctx)
	runtime.Show(d.ctx)

	d.eventEmitLaunchArgs(d.ctx, LaunchEvent{
		Args: secondInstanceArgs,
	})
}

// runWithCancellation is a helper function that allows to have request cancellation.
// Wails doesn't support context cancellation yet, so we have to do it ourselves.
// TODO: This can be simplified massivly. Can create a background context with cancel and store it against an ID, rather then true/false. Also don't need heartbeat.
func (d *Driver) runWithCancellation(cid uuid.UUID, requestFunc func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	d.cancelTokens.Store(cid.String(), false)
	defer d.cancelTokens.Delete(cid.String())

	ctx, cancel := context.WithCancel(d.ctx)
	defer cancel()

	resultChan := make(chan interface{})
	errChan := make(chan error)

	// Run the request function in its own goroutine
	go func() {
		defer close(resultChan)
		defer close(errChan)

		result, err := requestFunc(ctx)
		if err != nil {
			d.logger.Error("request function failed", map[string]interface{}{"error": err.Error(), "cid": cid})
			errChan <- err
			return
		}
		resultChan <- result
	}()

	// Start a ticker for heartbeats
	heartbeatTicker := time.NewTicker(d.heartbeatInterval)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-heartbeatTicker.C:
			d.logger.Debug("request heartbeat", map[string]interface{}{"cid": cid})
			runtime.EventsEmit(ctx, "requestHeartBeat", cid.String())

			cancelValue, ok := d.cancelTokens.Load(cid.String())
			if ok && cancelValue.(bool) {
				d.logger.Debug("request cancelled", map[string]interface{}{"cid": cid})
				return nil, errors.New("request cancelled")
			}
		case result := <-resultChan:
			return result, nil
		case err := <-errChan:
			return nil, err
		case <-ctx.Done():
			// Context cancelled (e.g., application-level cancellation)
			return nil, ctx.Err()
		}
	}
}
