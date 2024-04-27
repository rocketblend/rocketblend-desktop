package container

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
	"github.com/rocketblend/rocketblend/pkg/container"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
	"github.com/rocketblend/rocketblend/pkg/validator"
)

type (
	Options struct {
		Logger          logger.Logger
		Validator       types.Validator
		ApplicationName string
		Development     bool
	}

	Option func(*Options)

	Container struct {
		logger    types.Logger
		validator types.Validator

		applicationDir string

		rbContainer rbtypes.Container
	}
)

func WithLogger(logger types.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithValidator(validator types.Validator) Option {
	return func(o *Options) {
		o.Validator = validator
	}
}

func WithApplicationName(name string) Option {
	return func(o *Options) {
		o.ApplicationName = name
	}
}

func WithDevelopmentMode(development bool) Option {
	return func(o *Options) {
		o.Development = development
	}
}

func New(opts ...Option) (*Container, error) {
	options := &Options{
		Logger:          logger.NoOp(),
		Validator:       validator.New(),
		ApplicationName: types.ApplicationName,
	}

	for _, opt := range opts {
		opt(options)
	}

	rbContainer, err := container.New(
		container.WithLogger(options.Logger),
		container.WithValidator(options.Validator),
		container.WithApplicationName(options.ApplicationName),
		container.WithDevelopmentMode(options.Development),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create rocketblend container: %w", err)
	}

	applicationDir, err := setupApplicationDir(options.ApplicationName, options.Development)
	if err != nil {
		return nil, fmt.Errorf("failed to setup application directory: %w", err)
	}

	options.Logger.Debug("initializing container", map[string]interface{}{
		"path":        applicationDir,
		"development": options.Development,
	})

	return &Container{
		logger:         options.Logger,
		validator:      options.Validator,
		rbContainer:    rbContainer,
		applicationDir: applicationDir,
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
