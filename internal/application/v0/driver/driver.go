package driver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/buffermanager"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/configurator"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/dispatcher"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/operator"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/project"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/tracker"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
	rbcontainer "github.com/rocketblend/rocketblend/pkg/container"
	rbruntime "github.com/rocketblend/rocketblend/pkg/runtime"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
	rbvalidator "github.com/rocketblend/rocketblend/pkg/validator"
)

type (
	Options struct {
		Logger    types.Logger
		Validator types.Validator
		Store     types.Store

		// Dispatcher types.Dispatcher
		// Tracker    types.Tracker
		// Operator   types.Operator

		// RBConfigurator rbtypes.Configurator

		// Configurator types.Configurator

		// Portfolio    types.Portfolio
		// Catalog      types.Catalog

		ApplicationName string
		Development     bool

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
		store      types.Store

		operator     types.Operator
		configurator types.Configurator
		portfolio    types.Portfolio
		catalog      types.Catalog

		rbConfigurator rbtypes.Configurator

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
		Logger:          logger.NoOp(),
		Validator:       rbvalidator.New(),
		ApplicationName: types.ApplicationName,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Store == nil {
		return nil, errors.New("store is required")
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

	operator, err := operator.New(
		operator.WithLogger(options.Logger),
		operator.WithDispatcher(dispatcher),
		operator.WithStore(options.Store),
	)
	if err != nil {
		return nil, err
	}

	applicationDir, err := setupApplicationDir(options.ApplicationName, options.Development)
	if err != nil {
		return nil, fmt.Errorf("failed to setup application directory: %w", err)
	}

	configurator, err := configurator.New(
		configurator.WithLogger(options.Logger),
		configurator.WithValidator(options.Validator),
		configurator.WithLocation(applicationDir),
	)
	if err != nil {
		return nil, err
	}

	container, err := rbcontainer.New(
		rbcontainer.WithLogger(options.Logger),
		rbcontainer.WithValidator(options.Validator),
		rbcontainer.WithDevelopmentMode(options.Development),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create rocketblend container: %w", err)
	}

	rbConfigurator, err := container.GetConfigurator()
	if err != nil {
		return nil, fmt.Errorf("failed to get rocketblend configurator: %w", err)
	}

	rbRepository, err := container.GetRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to get rocketblend repository: %w", err)
	}

	portfolio, err := project.New(
		project.WithLogger(options.Logger),
		project.WithValidator(options.Validator),
		project.WithConfigurator(configurator),
		project.WithRocketBlendConfigurator(rbConfigurator),
		project.WithStore(options.Store),
	)
	if err != nil {
		return nil, err
	}

	return &Driver{
		logger:            options.Logger,
		validator:         options.Validator,
		store:             options.Store,
		dispatcher:        dispatcher,
		tracker:           tracker,
		operator:          operator,
		configurator:      configurator,
		rbConfigurator:    rbConfigurator,
		portfolio:         portfolio,
		catalog:           options.Catalog,
		events:            options.Writer,
		platform:          options.Platform,
		args:              options.Args,
		heartbeatInterval: 5000 * time.Millisecond, // 5 seconds
	}, nil
}

func setupApplicationDir(name string, development bool) (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot find config directory: %v", err)
	}

	appDir := filepath.Join(userConfigDir, name)
	if development {
		appDir = filepath.Join(appDir, "dev")
	}

	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create app directory: %w", err)
	}

	return appDir, nil
}
