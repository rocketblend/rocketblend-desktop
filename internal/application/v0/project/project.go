package project

import (
	"context"
	"errors"
	"path"
	"path/filepath"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/config"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
	"github.com/rocketblend/rocketblend-desktop/internal/application/watcher"
)

type (
	repository struct {
		logger    types.Logger
		validator types.Validator

		configurator config.Service

		rbRepository types.RBRepository
		rbDriver     types.RBDriver
		blender      types.Blender

		store      types.Store
		watcher    types.Watcher
		dispatcher types.Dispatcher
	}

	Options struct {
		Logger    types.Logger
		Validator types.Validator

		Configurator config.Service

		RBRepository types.RBRepository
		RBDriver     types.RBDriver
		Blender      types.Blender

		Store      types.Store
		Dispatcher types.Dispatcher

		WatcherDebounceDuration time.Duration
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithValidator(validator types.Validator) Option {
	return func(o *Options) {
		o.Validator = validator
	}
}

func WithWatcherDebounceDuration(duration time.Duration) Option {
	return func(o *Options) {
		o.WatcherDebounceDuration = duration
	}
}

func WithConfigurator(configurator config.Service) Option {
	return func(o *Options) {
		o.Configurator = configurator
	}
}

func WithRocketBlendRepository(repository types.RBRepository) Option {
	return func(o *Options) {
		o.RBRepository = repository
	}
}

func WithRocketBlendDriver(driver types.RBDriver) Option {
	return func(o *Options) {
		o.RBDriver = driver
	}
}

func WithBlender(blender types.Blender) Option {
	return func(o *Options) {
		o.Blender = blender
	}
}

func New(opts ...Option) (*repository, error) {
	options := &Options{
		Logger:                  logger.NoOp(),
		WatcherDebounceDuration: 2 * time.Second,
	}

	for _, o := range opts {
		o(options)
	}

	if options.Configurator == nil {
		return nil, errors.New("application configurator is required")
	}

	if options.Store == nil {
		return nil, errors.New("store is required")
	}

	if options.Dispatcher == nil {
		return nil, errors.New("dispatcher is required")
	}

	if options.RBRepository == nil {
		return nil, errors.New("rocketblend repository is required")
	}

	if options.RBDriver == nil {
		return nil, errors.New("rocketblend driver is required")
	}

	if options.Blender == nil {
		return nil, errors.New("blender is required")
	}

	config, err := options.Configurator.Get()
	if err != nil {
		return nil, err
	}

	watcher, err := watcher.New(
		watcher.WithLogger(options.Logger),
		watcher.WithEventDebounceDuration(options.WatcherDebounceDuration),
		watcher.WithPaths(config.Project.Paths...),
		watcher.WithIsWatchableFileFunc(func(path string) bool {
			for _, ext := range config.Project.Watcher.FileExtensions {
				if filepath.Ext(path) == ext {
					return true
				}
			}

			return false
		}),
		watcher.WithResolveObjectPathFunc(func(path string) string {
			return filepath.Dir(path)
		}),
		watcher.WithUpdateObjectFunc(func(path string) error {
			// project, err := project.Load(path)
			// if err != nil {
			// 	return fmt.Errorf("failed to load project %s: %w", path, err)
			// }

			// index, err := convertToIndex(project)
			// if err != nil {
			// 	return err
			// }

			// options.Logger.Debug("updating project index", map[string]interface{}{
			// 	"id":        index.ID,
			// 	"reference": index.Reference,
			// })

			// // TODO: Pass context from watcher
			// return options.Store.Insert(context.Background(), index)

			return nil
		}),
		watcher.WithRemoveObjectFunc(func(removePath string) error {
			return options.Store.RemoveByReference(context.Background(), path.Clean(removePath))
		}),
	)
	if err != nil {
		return nil, err
	}

	return &repository{
		logger:       options.Logger,
		configurator: options.Configurator,
		validator:    options.Validator,
		rbRepository: options.RBRepository,
		rbDriver:     options.RBDriver,
		blender:      options.Blender,
		store:        options.Store,
		dispatcher:   options.Dispatcher,
		watcher:      watcher,
	}, nil
}

func (r *repository) Close() error {
	return r.watcher.Close()
}
