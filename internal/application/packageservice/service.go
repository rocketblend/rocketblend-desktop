package packageservice

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	pack "github.com/rocketblend/rocketblend-desktop/internal/application/package"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoptions"
	"github.com/rocketblend/rocketblend-desktop/internal/application/watcher"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketfile"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/config"
)

type (
	Service interface {
		Get(ctx context.Context, id uuid.UUID) (*GetPackageResponse, error)
		List(ctx context.Context, opts ...listoptions.ListOption) (*ListPackagesResponse, error)
	}

	service struct {
		logger logger.Logger

		config  *config.Service //Swtich to interface
		store   searchstore.Store
		watcher watcher.Watcher
	}

	Options struct {
		Logger                  logger.Logger
		Config                  *config.Service
		Store                   searchstore.Store
		WatcherDebounceDuration time.Duration
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithConfig(config *config.Service) Option {
	return func(o *Options) {
		o.Config = config
	}
}

func WithStore(store searchstore.Store) Option {
	return func(o *Options) {
		o.Store = store
	}
}

func WithWatcherDebounceDuration(duration time.Duration) Option {
	return func(o *Options) {
		o.WatcherDebounceDuration = duration
	}
}

func New(opts ...Option) (Service, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	if options.Store == nil {
		return nil, fmt.Errorf("store service is required")
	}

	if options.Config == nil {
		return nil, fmt.Errorf("config service is required")
	}

	config, err := options.Config.Get()
	if err != nil {
		return nil, err
	}

	watcher, err := watcher.New(
		watcher.WithLogger(options.Logger),
		watcher.WithEventDebounceDuration(options.WatcherDebounceDuration),
		watcher.WithPaths(config.PackagesPath),
		watcher.WithIsWatchableFileFunc(func(path string) bool {
			return filepath.Base(path) == rocketfile.FileName
		}),
		watcher.WithUpdateObjectFunc(func(path string) error {
			project, err := pack.Load(config.PackagesPath, config.InstallationsPath, path)
			if err != nil {
				return fmt.Errorf("failed to load project %s: %w", path, err)
			}

			data, err := json.Marshal(project)
			if err != nil {
				return err
			}

			return options.Store.Insert(&searchstore.Index{
				ID:   project.ID,
				Type: searchstore.Package,
				Path: path,
				Data: string(data),
			})
		}),
		watcher.WithRemoveObjectFunc(func(path string) error {
			return options.Store.RemoveByPath(path)
		}),
	)
	if err != nil {
		return nil, err
	}

	return &service{
		logger:  options.Logger,
		config:  options.Config,
		store:   options.Store,
		watcher: watcher,
	}, nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*GetPackageResponse, error) {
	index, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}

	pack, err := s.convert(index)
	if err != nil {
		return nil, err
	}

	return &GetPackageResponse{
		Package: pack,
	}, nil
}

func (s *service) List(ctx context.Context, opts ...listoptions.ListOption) (*ListPackagesResponse, error) {
	indexes, err := s.store.List(opts...)
	if err != nil {
		return nil, err
	}

	response := &ListPackagesResponse{
		Packages: make([]*pack.Package, 0),
	}

	for _, index := range indexes {
		pack, err := s.convert(index)
		if err != nil {
			return nil, err
		}

		response.Packages = append(response.Packages, pack)
	}

	return response, nil
}

func (s *service) convert(index *searchstore.Index) (*pack.Package, error) {
	var result pack.Package
	if err := json.Unmarshal([]byte(index.Data), &result); err == nil {
		return nil, err
	}

	return &result, nil
}