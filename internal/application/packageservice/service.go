package packageservice

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/eventservice"
	pack "github.com/rocketblend/rocketblend-desktop/internal/application/package"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
	"github.com/rocketblend/rocketblend-desktop/internal/application/watcher"

	rocketblendInstallation "github.com/rocketblend/rocketblend/pkg/driver/installation"
	rocketblendPackage "github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/reference"
	rocketblendConfig "github.com/rocketblend/rocketblend/pkg/rocketblend/config"
)

type (
	Service interface {
		Get(ctx context.Context, id uuid.UUID) (*GetPackageResponse, error)
		List(ctx context.Context, opts ...listoption.ListOption) (*ListPackagesResponse, error)

		AppendOperation(ctx context.Context, id uuid.UUID, opid uuid.UUID) error

		Add(ctx context.Context, reference reference.Reference) error
		Install(ctx context.Context, id uuid.UUID) error
		Uninstall(ctx context.Context, id uuid.UUID) error

		Refresh(ctx context.Context) error

		Close() error
	}

	service struct {
		logger logger.Logger

		rocketblendConfigService       rocketblendConfig.Service
		rocketblendPackageService      rocketblendPackage.Service
		rocketblendInstallationService rocketblendInstallation.Service

		store      searchstore.Store
		watcher    watcher.Watcher
		dispatcher eventservice.Service
	}

	Options struct {
		Logger logger.Logger

		RocketblendConfigService       rocketblendConfig.Service
		RocketblendPackageService      rocketblendPackage.Service
		RocketblendInstallationService rocketblendInstallation.Service

		Store      searchstore.Store
		Dispatcher eventservice.Service

		WatcherDebounceDuration time.Duration
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithRocketBlendConfigService(srv rocketblendConfig.Service) Option {
	return func(o *Options) {
		o.RocketblendConfigService = srv
	}
}

func WithRocketBlendPackageService(srv rocketblendPackage.Service) Option {
	return func(o *Options) {
		o.RocketblendPackageService = srv
	}
}

func WithRocketBlendInstallationService(srv rocketblendInstallation.Service) Option {
	return func(o *Options) {
		o.RocketblendInstallationService = srv
	}
}

func WithStore(store searchstore.Store) Option {
	return func(o *Options) {
		o.Store = store
	}
}

func WithDispatcher(dispatcher eventservice.Service) Option {
	return func(o *Options) {
		o.Dispatcher = dispatcher
	}
}

func WithWatcherDebounceDuration(duration time.Duration) Option {
	return func(o *Options) {
		o.WatcherDebounceDuration = duration
	}
}

func New(opts ...Option) (Service, error) {
	options := &Options{
		Logger:                  logger.NoOp(),
		WatcherDebounceDuration: 2 * time.Second,
	}

	for _, o := range opts {
		o(options)
	}

	if options.Store == nil {
		return nil, fmt.Errorf("store service is required")
	}

	if options.Dispatcher == nil {
		return nil, fmt.Errorf("dispatcher service is required")
	}

	if options.RocketblendConfigService == nil {
		return nil, fmt.Errorf("rocketblend config service is required")
	}

	if options.RocketblendPackageService == nil {
		return nil, fmt.Errorf("rocketblend package service is required")
	}

	if options.RocketblendInstallationService == nil {
		return nil, fmt.Errorf("rocketblend installation service is required")
	}

	config, err := options.RocketblendConfigService.Get()
	if err != nil {
		return nil, err
	}

	watcher, err := watcher.New(
		watcher.WithLogger(options.Logger),
		watcher.WithEventDebounceDuration(options.WatcherDebounceDuration),
		watcher.WithPaths(config.PackagesPath),
		watcher.WithIsWatchableFileFunc(func(path string) bool {
			return filepath.Base(path) == rocketblendPackage.FileName
		}),
		watcher.WithUpdateObjectFunc(func(path string) error {
			pack, err := pack.Load(config.PackagesPath, config.InstallationsPath, path, config.Platform)
			if err != nil {
				return fmt.Errorf("failed to load package %s: %w", path, err)
			}

			index, err := convertToIndex(pack)
			if err != nil {
				return err
			}

			// TODO: Pass context from watcher
			return options.Store.Insert(context.Background(), index)
		}),
		watcher.WithRemoveObjectFunc(func(removePath string) error {
			return options.Store.RemoveByReference(context.Background(), path.Clean(removePath))
		}),
	)
	if err != nil {
		return nil, err
	}

	return &service{
		logger:                         options.Logger,
		rocketblendConfigService:       options.RocketblendConfigService,
		rocketblendPackageService:      options.RocketblendPackageService,
		rocketblendInstallationService: options.RocketblendInstallationService,
		store:                          options.Store,
		dispatcher:                     options.Dispatcher,
		watcher:                        watcher,
	}, nil
}

func (s *service) Close() error {
	return s.watcher.Close()
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*GetPackageResponse, error) {
	pack, err := s.get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &GetPackageResponse{
		Package: pack,
	}, nil
}

func (s *service) List(ctx context.Context, opts ...listoption.ListOption) (*ListPackagesResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	opts = append(opts, listoption.WithType(indextype.Package))
	indexes, err := s.store.List(ctx, opts...)
	if err != nil {
		return nil, err
	}

	packs := make([]*pack.Package, 0, len(indexes))
	for _, index := range indexes {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		pack, err := convertFromIndex(index)
		if err != nil {
			return nil, err
		}
		packs = append(packs, pack)
	}

	return &ListPackagesResponse{
		Packages: packs,
	}, nil
}

func (s *service) AppendOperation(ctx context.Context, id uuid.UUID, opid uuid.UUID) error {
	pack, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	pack.Operations = append(pack.Operations, opid.String())

	return s.insert(ctx, pack)
}

func (s *service) get(ctx context.Context, id uuid.UUID) (*pack.Package, error) {
	index, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	pack, err := convertFromIndex(index)
	if err != nil {
		return nil, err
	}

	return pack, nil
}

func (s *service) insert(ctx context.Context, pack *pack.Package) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	index, err := convertToIndex(pack)
	if err != nil {
		return err
	}

	return s.store.Insert(ctx, index)
}

func convertFromIndex(index *searchstore.Index) (*pack.Package, error) {
	var result pack.Package
	if err := json.Unmarshal([]byte(index.Data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func convertToIndex(pack *pack.Package) (*searchstore.Index, error) {
	data, err := json.Marshal(pack)
	if err != nil {
		return nil, err
	}

	return &searchstore.Index{
		ID:         pack.ID,
		Name:       pack.Name,
		Type:       indextype.Package,
		Reference:  path.Clean(pack.Path),
		Category:   strconv.Itoa(int(pack.Type)),
		State:      int(pack.State),
		Operations: pack.Operations,
		Data:       string(data),
	}, nil
}
