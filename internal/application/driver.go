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
	"github.com/rocketblend/rocketblend-desktop/internal/application/factory"
	rbruntime "github.com/rocketblend/rocketblend/pkg/runtime"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Driver struct
type Driver struct {
	ctx    context.Context
	logger logger.Logger

	heartbeatInterval time.Duration

	factory      factory.Factory
	events       buffermanager.BufferManager
	cancelTokens sync.Map

	platform rbruntime.Platform

	args []string
}

func NewDriver(factory factory.Factory, events buffermanager.BufferManager, platform rbruntime.Platform, args ...string) (*Driver, error) {
	logger, err := factory.GetLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to get logger: %w", err)
	}

	return &Driver{
		factory:           factory,
		heartbeatInterval: 5000 * time.Millisecond, // 1 second
		events:            events,
		logger:            logger,
		platform:          platform,
		args:              args,
	}, nil
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
