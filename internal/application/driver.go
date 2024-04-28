package application

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
	"github.com/rocketblend/rocketblend-desktop/internal/buffer"
	rbruntime "github.com/rocketblend/rocketblend/pkg/runtime"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

type (
	dependecies struct {
		logger    types.Logger
		validator types.Validator

		dispatcher types.Dispatcher
		tracker    types.Tracker
		operator   types.Operator

		portfolio types.Portfolio
		catalog   types.Catalog

		configurator   types.Configurator
		rbConfigurator rbtypes.Configurator
	}

	Options struct {
		Container types.Container

		Writer   buffer.BufferManager // TODO: Improve this
		Platform rbruntime.Platform
		Version  string
		Args     []string
	}

	Option func(*Options)

	Driver struct {
		logger    types.Logger
		validator types.Validator

		dispatcher types.Dispatcher
		tracker    types.Tracker
		operator   types.Operator

		portfolio types.Portfolio
		catalog   types.Catalog

		configurator   types.Configurator
		rbConfigurator rbtypes.Configurator

		ctx               context.Context
		version           string
		heartbeatInterval time.Duration // TODO: Remove this
		events            buffer.BufferManager
		cancelTokens      sync.Map
		platform          rbruntime.Platform
		args              []string
	}
)

func WithContainer(container types.Container) Option {
	return func(o *Options) {
		o.Container = container
	}
}

func WithWriter(writer buffer.BufferManager) Option {
	return func(o *Options) {
		o.Writer = writer
	}
}

func WithPlatform(platform rbruntime.Platform) Option {
	return func(o *Options) {
		o.Platform = platform
	}
}

func WithVersion(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

func WithArgs(args ...string) Option {
	return func(o *Options) {
		o.Args = args
	}
}

func NewDriver(opts ...Option) (*Driver, error) {
	options := &Options{}

	for _, opt := range opts {
		opt(options)
	}

	if options.Container == nil {
		return nil, errors.New("container is required")
	}

	dependencies, err := setupDependencies(options.Container)
	if err != nil {
		return nil, fmt.Errorf("failed to setup dependencies: %w", err)
	}

	return &Driver{
		logger:            dependencies.logger,
		validator:         dependencies.validator,
		dispatcher:        dependencies.dispatcher,
		tracker:           dependencies.tracker,
		operator:          dependencies.operator,
		portfolio:         dependencies.portfolio,
		catalog:           dependencies.catalog,
		configurator:      dependencies.configurator,
		rbConfigurator:    dependencies.rbConfigurator,
		events:            options.Writer,
		platform:          options.Platform,
		version:           options.Version,
		args:              options.Args,
		heartbeatInterval: 5000 * time.Millisecond, // 5 seconds
	}, nil
}

func setupDependencies(container types.Container) (*dependecies, error) {
	logger, err := container.GetLogger()
	if err != nil {
		return nil, err
	}

	validator, err := container.GetValidator()
	if err != nil {
		return nil, err
	}

	dispatcher, err := container.GetDispatcher()
	if err != nil {
		return nil, err
	}

	tracker, err := container.GetTracker()
	if err != nil {
		return nil, err
	}

	operator, err := container.GetOperator()
	if err != nil {
		return nil, err
	}

	configurator, err := container.GetConfigurator()
	if err != nil {
		return nil, err
	}

	portfolio, err := container.GetPortfolio()
	if err != nil {
		return nil, err
	}

	catalog, err := container.GetCatalog()
	if err != nil {
		return nil, err
	}

	rbConfigurator, err := container.GetRBConfigurator()
	if err != nil {
		return nil, err
	}

	return &dependecies{
		logger:         logger,
		validator:      validator,
		dispatcher:     dispatcher,
		tracker:        tracker,
		operator:       operator,
		configurator:   configurator,
		rbConfigurator: rbConfigurator,
		portfolio:      portfolio,
		catalog:        catalog,
	}, nil
}
