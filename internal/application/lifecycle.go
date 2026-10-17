package application

import (
	"context"
	"strings"

	"github.com/rocketblend/rocketblend-desktop/internal/buffer"
)

type (
	lifecycle struct {
		*dependencies

		ctx    context.Context
		ui     UI
		events buffer.BufferManager // TODO: Improve this
		args   []string
	}
)

func newLifecycle(ctx context.Context, deps *dependencies, ui UI, events buffer.BufferManager, args []string) *lifecycle {
	return &lifecycle{
		dependencies: deps,
		ctx:          ctx,
		ui:           ui,
		events:       events,
		args:         args,
	}
}

// startup is called when the app starts.
func (l *lifecycle) startup() {
	l.logger.Debug("starting application")

	// Start listening to log events
	go l.listenToLogEvents()
}

// shutdown is called when the app is shutting down
func (l *lifecycle) shutdown() {
	l.logger.Debug("shutting down application")

	// Close the event stream
	l.events.Close()

	l.logger.Debug("application shutdown")
}

// onDomReady is called when the DOM is ready
func (l *lifecycle) onDomReady() {
	l.logger.Debug("DOM is ready")

	// Wait for main layout to be ready.
	l.ui.Once("ready", func(optionalData ...interface{}) {
		l.onLayoutReady()
	})
}

// onLayoutReady is called when the layout is ready
func (l *lifecycle) onLayoutReady() {
	l.logger.Debug("main layout is ready")

	if err := l.addApplicationMetrics(); err != nil {
		l.logger.Error("failed to add application metrics", map[string]interface{}{"error": err.Error()})
	}

	// Setup driver event handlers (backend)
	if err := l.setupDriverEventHandlers(); err != nil {
		l.logger.Error("failed to setup driver event handlers", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := l.autoRefreshPackages(); err != nil {
		l.logger.Error("failed to auto refresh packages", map[string]interface{}{"error": err.Error()})
	}

	args := []string{}
	if len(l.args) > 1 {
		args = l.args[1:]
	}

	l.eventEmitLaunchArgs(LaunchEvent{
		Args: args,
	})
}

// onSecondInstanceLaunch is called when the user opens a second instance of the application
func (l *lifecycle) onSecondInstanceLaunch(args []string, workingDirectory string) {
	l.logger.Info("user opened second instance", map[string]interface{}{
		"args":             strings.Join(args, ","),
		"workingDirectory": workingDirectory,
	})

	l.ui.ShowWindow()

	l.eventEmitLaunchArgs(LaunchEvent{
		Args: args,
	})
}

// autoRefreshPackages refreshes the packages, if auto pull is enabled.
func (l *lifecycle) autoRefreshPackages() error {
	config, err := l.configurator.Get()
	if err != nil {
		return err
	}

	if !config.Package.AutoPull {
		return nil
	}

	if err := l.catalog.RefreshPackages(l.ctx); err != nil {
		return err
	}

	return nil
}
