package application

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/buffermanager"
	"github.com/rocketblend/rocketblend-desktop/internal/application/config"
	"github.com/rocketblend/rocketblend-desktop/internal/application/factory"
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
		d.logger.Error("Failed to get config service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	config, err := configService.Get()
	if err != nil {
		d.logger.Error("Failed to get config", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return &config.Platform, nil
}

func (d *Driver) GetApplicationConfig() (*config.Config, error) {
	configService, err := d.factory.GetApplicationConfigService()
	if err != nil {
		d.logger.Error("Failed to get application config service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	config, err := configService.Get()
	if err != nil {
		d.logger.Error("Failed to get config", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return config, nil
}

func (d *Driver) GetRocketBlendConfig() (*rocketblendConfig.Config, error) {
	configService, err := d.factory.GetConfigService()
	if err != nil {
		d.logger.Error("Failed to get rocketblend config service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	config, err := configService.Get()
	if err != nil {
		d.logger.Error("Failed to get config", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return config, nil
}

// GetProject gets a project by id
func (d *Driver) GetProject(id uuid.UUID) (*projectservice.GetProjectResponse, error) {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("Failed to get project service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	project, err := projectService.Get(ctx, id)
	if err != nil {
		d.logger.Error("Failed to find project by id", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return project, nil
}

// ListProjects lists all projects
func (d *Driver) ListProjects(query string) (*projectservice.ListProjectsResponse, error) {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("Failed to get project service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	response, err := projectService.List(ctx, listoption.WithQuery(query))
	if err != nil {
		d.logger.Error("Failed to find all projects", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	d.logger.Debug("Found projects", map[string]interface{}{"projects": len(response.Projects)})

	return response, nil
}

// CreateProject creates a new project
func (d *Driver) CreateProject(request *projectservice.CreateProjectRequest) error {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("Failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Create(ctx, request); err != nil {
		d.logger.Error("Failed to create project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("Project created")
	return nil
}

// UpdateProject updates a project
func (d *Driver) UpdateProject(request *projectservice.UpdateProjectRequest) error {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("Failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Update(ctx, request); err != nil {
		d.logger.Error("Failed to update project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("Project updated", map[string]interface{}{"id": request.ID})

	return nil
}

// DeleteProject deletes a project
func (d *Driver) DeleteProject(id uuid.UUID) error {
	d.logger.Debug("Deleting project", map[string]interface{}{"id": id})
	return nil
}

// RunProject runs a project
func (d *Driver) RunProject(id uuid.UUID) error {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("Failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Run(ctx, id); err != nil {
		d.logger.Error("Failed to run project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("Project started", map[string]interface{}{"id": id})
	return nil
}

// RenderProject renders a project
func (d *Driver) RenderProject(id uuid.UUID) error {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("Failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Render(ctx, id); err != nil {
		d.logger.Error("Failed to render project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("Project rendered", map[string]interface{}{"id": id})
	return nil
}

// ExploreProject explores a project
func (d *Driver) ExploreProject(id uuid.UUID) error {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("Failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Explore(ctx, id); err != nil {
		d.logger.Error("Failed to explore project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("Project explored", map[string]interface{}{"id": id})
	return nil
}

func (d *Driver) GetPackage(id uuid.UUID) (*packageservice.GetPackageResponse, error) {
	ctx := context.Background()

	packageService, err := d.factory.GetPackageService()
	if err != nil {
		d.logger.Error("Failed to get package service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	pack, err := packageService.Get(ctx, id)
	if err != nil {
		d.logger.Error("Failed to find package by id", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return pack, err
}

func (d *Driver) ListPackages(query string, category string, installed bool) (*packageservice.ListPackagesResponse, error) {
	ctx := context.Background()

	packageService, err := d.factory.GetPackageService()
	if err != nil {
		d.logger.Error("Failed to get package service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	response, err := packageService.List(ctx, []listoption.ListOption{
		listoption.WithQuery(query),
		listoption.WithCategory(category),
		listoption.WithReady(installed),
	}...)
	if err != nil {
		d.logger.Error("Failed to find all packages", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	d.logger.Debug("Found packages", map[string]interface{}{"packages": len(response.Packages)})

	return response, err
}

func (d *Driver) AddPackageWithOperation(opid uuid.UUID, referenceStr string) error {
	_, err := d.withCancellation(opid, func(ctx context.Context) (interface{}, error) {
		packageService, err := d.factory.GetPackageService()
		if err != nil {
			d.logger.Error("Failed to get package service", map[string]interface{}{"error": err.Error(), "opid": opid})
			return nil, err
		}

		ref, err := reference.Parse(referenceStr)
		if err != nil {
			d.logger.Error("Failed to parse reference", map[string]interface{}{"error": err.Error(), "opid": opid})
			return nil, err
		}

		if err := packageService.Add(ctx, ref); err != nil {
			d.logger.Error("Failed to add package", map[string]interface{}{"error": err.Error(), "opid": opid})
			return nil, err
		}

		return struct{}{}, nil
	})

	if err != nil {
		return err
	}

	d.logger.Debug("Package added", map[string]interface{}{"reference": referenceStr, "opid": opid})
	return nil
}

func (d *Driver) RefreshPackagesWithOperation(opid uuid.UUID) error {
	_, err := d.withCancellation(opid, func(ctx context.Context) (interface{}, error) {
		packageService, err := d.factory.GetPackageService()
		if err != nil {
			d.logger.Error("Failed to get package service", map[string]interface{}{"error": err.Error(), "opid": opid})
			return nil, err
		}

		if err := packageService.Refresh(ctx); err != nil {
			d.logger.Error("Failed to refresh packages", map[string]interface{}{"error": err.Error(), "opid": opid})
			return nil, err
		}

		return struct{}{}, nil
	})
	if err != nil {
		return err
	}

	d.logger.Debug("Packages refreshed!", map[string]interface{}{"opid": opid})
	return nil
}

func (d *Driver) InstallPackageWithOperation(opid uuid.UUID, id uuid.UUID) error {
	_, err := d.withCancellation(opid, func(ctx context.Context) (interface{}, error) {
		packageService, err := d.factory.GetPackageService()
		if err != nil {
			d.logger.Error("Failed to get package service", map[string]interface{}{"error": err.Error(), "opid": opid})
			return nil, err
		}

		if err := packageService.Install(ctx, id); err != nil {
			d.logger.Error("Failed to install package", map[string]interface{}{"error": err.Error(), "opid": opid})
			return nil, err
		}

		d.logger.Debug("Package installed", map[string]interface{}{"id": id})
		return struct{}{}, nil
	})
	if err != nil {
		return err
	}

	d.logger.Debug("Package installation completed!", map[string]interface{}{"id": id, "opid": opid})
	return nil
}

func (d *Driver) UninstallPackageWithOperation(opid uuid.UUID, id uuid.UUID) error {
	_, err := d.withCancellation(opid, func(ctx context.Context) (interface{}, error) {
		packageService, err := d.factory.GetPackageService()
		if err != nil {
			d.logger.Error("Failed to get package service", map[string]interface{}{"error": err.Error()})
			return nil, err
		}

		if err := packageService.Uninstall(ctx, id); err != nil {
			d.logger.Error("Failed to uninstall package", map[string]interface{}{"error": err.Error()})
			return nil, err
		}

		d.logger.Debug("Package uninstalled", map[string]interface{}{"id": id})
		return struct{}{}, nil
	})

	if err != nil {
		return err
	}

	d.logger.Debug("Package uninstallation completed!", map[string]interface{}{"id": id, "opid": opid})
	return nil
}

func (d *Driver) LongRunningWithOperation(opid uuid.UUID) error {
	_, err := d.withCancellation(opid, func(ctx context.Context) (interface{}, error) {
		// Simulate a long-running operation
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				d.logger.Debug("Long running operation canceled", map[string]interface{}{"opid": opid})
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

	d.logger.Debug("Long running operation completed!", map[string]interface{}{"opID": opid})
	return nil
}

// Quit quits the application
func (d *Driver) Quit() {
	d.logger.Debug("Quitting application")
	runtime.Quit(d.ctx)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (d *Driver) startup(ctx context.Context) {
	d.logger.Debug("Starting application")

	d.ctx = ctx

	// Start listening to log events
	go d.listenToLogEvents()

	runtime.EventsOn(ctx, "cancelOperation", func(optionalData ...interface{}) {
		if len(optionalData) > 0 {
			if opIDStr, ok := optionalData[0].(string); ok {
				opID, err := uuid.Parse(opIDStr)
				if err != nil {
					d.logger.Error("Invalid operation ID format for cancellation", map[string]interface{}{"error": err.Error()})
					return
				}
				d.cancelTokens.Store(opID.String(), true)
				d.logger.Debug("Cancellation requested", map[string]interface{}{"opID": opID})
			} else {
				d.logger.Error("Invalid data type for operation ID", map[string]interface{}{"type": fmt.Sprintf("%T", optionalData[0])})
			}
		} else {
			d.logger.Error("No operation ID provided for cancellation")
		}
	})

	// Preloads all the data
	if err := d.factory.Preload(); err != nil {
		d.logger.Error("Failed to preload", map[string]interface{}{"error": err.Error()})
		return
	}
}

// shutdown is called when the app is shutting down
func (d *Driver) shutdown(ctx context.Context) {
	d.logger.Debug("Shutting down application")

	// Close the event stream
	d.events.Close()

	// Close the factory watchers
	if err := d.factory.Close(); err != nil {
		d.logger.Error("Failed to close factory", map[string]interface{}{"error": err.Error()})
	}

	d.logger.Debug("Application shutdown")
}

// onDomReady is called when the DOM is ready
func (d *Driver) onDomReady(ctx context.Context) {
	d.logger.Debug("DOM is ready")

	// Wait for main layout to be ready.
	runtime.EventsOnce(ctx, "ready", func(optionalData ...interface{}) {
		d.logger.Debug("Main layout is ready")
		d.eventEmitLaunchArgs(ctx, LaunchEvent{
			Args: os.Args[1:],
		})
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

func (d *Driver) listenToLogEvents() {
	for {
		select {
		case <-d.ctx.Done():
			return
		default:
			data, ok := d.events.GetNextData()
			if ok {
				if logEvent, isLogEvent := data.(LogEvent); isLogEvent {
					runtime.EventsEmit(d.ctx, "logStream", logEvent)
				}
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}

func (d *Driver) eventEmitLaunchArgs(ctx context.Context, event LaunchEvent) {
	d.logger.Debug("emitting launchArgs event", map[string]interface{}{
		"event": event,
	})

	runtime.EventsEmit(ctx, "launchArgs", event)
}

// withCancellation is a helper function that runs an operation with a cancellation token.
// Wails doesn't support context cancellation yet, so we have to do it ourselves.
func (d *Driver) withCancellation(opID uuid.UUID, operation func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	d.cancelTokens.Store(opID.String(), false)
	defer d.cancelTokens.Delete(opID.String())

	ctx, cancel := context.WithCancel(d.ctx)
	defer cancel()

	resultChan := make(chan interface{})
	errChan := make(chan error)

	// Run the operation in its own goroutine
	go func() {
		defer close(resultChan)
		defer close(errChan)

		result, err := operation(ctx)
		if err != nil {
			d.logger.Error("Operation failed", map[string]interface{}{"error": err.Error(), "opID": opID})
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
			d.logger.Debug("Operation heartbeat", map[string]interface{}{"opID": opID})
			runtime.EventsEmit(ctx, "operationHeartBeat", opID.String())

			cancelValue, ok := d.cancelTokens.Load(opID.String())
			if ok && cancelValue.(bool) {
				d.logger.Debug("Operation cancelled", map[string]interface{}{"opID": opID})
				return nil, errors.New("operation cancelled")
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
