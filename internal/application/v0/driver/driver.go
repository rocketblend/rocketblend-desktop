package driver

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/buffermanager"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/dispatcher"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/tracker"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
	rbruntime "github.com/rocketblend/rocketblend/pkg/runtime"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Options struct {
		Logger    types.Logger
		Validator types.Validator

		// Dispatcher types.Dispatcher
		// Tracker    types.Tracker
		// Operator   types.Operator

		// RBConfigurator rbtypes.Configurator

		// Configurator types.Configurator
		// Store        types.Store
		// Portfolio    types.Portfolio
		// Catalog      types.Catalog

		Writer   buffermanager.BufferManager // TODO: Improve this
		Platform rbruntime.Platform
		Args     []string
	}

	Option func(*Options)

	Driver struct {
		logger    types.Logger
		validator types.Validator

		dispatcher types.Dispatcher
		tracker    types.Tracker
		operator   types.Operator

		rbConfigurator rbtypes.Configurator

		configurator types.Configurator
		store        types.Store
		portfolio    types.Portfolio
		catalog      types.Catalog

		ctx               context.Context
		heartbeatInterval time.Duration // TODO: Remove this
		events            buffermanager.BufferManager
		cancelTokens      sync.Map
		platform          rbruntime.Platform
		args              []string
	}
)

func New(opts ...Option) (*Driver, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Validator == nil {
		return nil, errors.New("validator is required")
	}

	dispatcher, err := dispatcher.New(
		dispatcher.WithLogger(options.Logger),
	)
	if err != nil {
		return nil, err
	}

	tracker, err := tracker.New(
		tracker.WithLogger(options.Logger),
	)
	if err != nil {
		return nil, err
	}

	return &Driver{
		logger:            options.Logger,
		validator:         options.Validator,
		dispatcher:        dispatcher,
		tracker:           tracker,
		operator:          options.Operator,
		rbConfigurator:    options.RBConfigurator,
		configurator:      options.Configurator,
		store:             options.Store,
		portfolio:         options.Portfolio,
		catalog:           options.Catalog,
		events:            options.Writer,
		platform:          options.Platform,
		args:              options.Args,
		heartbeatInterval: 5000 * time.Millisecond, // 5 seconds
	}, nil
}
