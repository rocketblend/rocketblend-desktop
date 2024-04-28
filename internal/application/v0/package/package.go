package pack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
	"github.com/rocketblend/rocketblend-desktop/internal/application/watcher"
)

type (
	Repository struct {
		logger    types.Logger
		validator types.Validator

		rbRepository   types.RBRepository
		rbConfigurator types.RBConfigurator

		store      types.Store
		watcher    types.Watcher
		dispatcher types.Dispatcher
	}

	Options struct {
		Logger    types.Logger
		Validator types.Validator

		rbRepository   types.RBRepository
		rbConfigurator types.RBConfigurator

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

func WithRocketBlendRepository(repository types.RBRepository) Option {
	return func(o *Options) {
		o.rbRepository = repository
	}
}

func WithRocketBlendConfigurator(configurator types.RBConfigurator) Option {
	return func(o *Options) {
		o.rbConfigurator = configurator
	}
}

func WithStore(store types.Store) Option {
	return func(o *Options) {
		o.Store = store
	}
}

func WithDispatcher(dispatcher types.Dispatcher) Option {
	return func(o *Options) {
		o.Dispatcher = dispatcher
	}
}

func WithWatcherDebounceDuration(duration time.Duration) Option {
	return func(o *Options) {
		o.WatcherDebounceDuration = duration
	}
}

func New(opts ...Option) (*Repository, error) {
	options := &Options{
		Logger:                  logger.NoOp(),
		WatcherDebounceDuration: 2 * time.Second,
	}

	for _, o := range opts {
		o(options)
	}

	if options.Validator == nil {
		return nil, errors.New("validator is required")
	}

	if options.Store == nil {
		return nil, errors.New("store is required")
	}

	if options.Dispatcher == nil {
		return nil, errors.New("dispatcher is required")
	}

	if options.rbRepository == nil {
		return nil, fmt.Errorf("rocketblend repository is required")
	}

	if options.rbConfigurator == nil {
		return nil, fmt.Errorf("rocketblend configurator is required")
	}

	config, err := options.rbConfigurator.Get()
	if err != nil {
		return nil, err
	}

	watcher, err := watcher.New(
		watcher.WithLogger(options.Logger),
		watcher.WithEventDebounceDuration(options.WatcherDebounceDuration),
		watcher.WithPaths(config.PackagesPath),
		watcher.WithIsWatchableFileFunc(func(path string) bool {
			return filepath.Base(path) == types.PackageFileName
		}),
		watcher.WithUpdateObjectFunc(func(path string) error {
			pack, err := load(options.rbConfigurator, options.Validator, path)
			if err != nil {
				return fmt.Errorf("failed to load package %s: %w", path, err)
			}

			index, err := convertToIndex(pack)
			if err != nil {
				return err
			}

			return options.Store.Insert(context.Background(), index)
		}),
		watcher.WithRemoveObjectFunc(func(removePath string) error {
			return options.Store.RemoveByReference(context.Background(), path.Clean(removePath))
		}),
	)
	if err != nil {
		return nil, err
	}

	return &Repository{
		logger:         options.Logger,
		validator:      options.Validator,
		rbRepository:   options.rbRepository,
		rbConfigurator: options.rbConfigurator,
		store:          options.Store,
		dispatcher:     options.Dispatcher,
		watcher:        watcher,
	}, nil
}

func (r *Repository) Close() error {
	return r.watcher.Close()
}

func convertFromIndex(index *types.Index) (*types.Package, error) {
	var result types.Package
	if err := json.Unmarshal([]byte(index.Data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func convertToIndex(pack *types.Package) (*types.Index, error) {
	data, err := json.Marshal(pack)
	if err != nil {
		return nil, err
	}

	return &types.Index{
		ID:        pack.ID,
		Name:      pack.Name,
		Type:      indextype.Package,
		Reference: path.Clean(pack.Path),
		Category:  string(pack.Type),
		//State:      int(pack.State),
		Operations: pack.Operations,
		Data:       string(data),
	}, nil
}
