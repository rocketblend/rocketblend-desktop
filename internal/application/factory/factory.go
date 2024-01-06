package factory

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rocketblend/rocketblend-desktop/internal/application/build"
	"github.com/rocketblend/rocketblend-desktop/internal/application/config"
	"github.com/rocketblend/rocketblend-desktop/internal/application/eventservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/operationservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/packageservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"

	rocketblendBlendFile "github.com/rocketblend/rocketblend/pkg/driver/blendfile"
	rocketblendInstallation "github.com/rocketblend/rocketblend/pkg/driver/installation"
	rocketblendRocketPack "github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	rocketblendConfig "github.com/rocketblend/rocketblend/pkg/rocketblend/config"
	rocketblendFactory "github.com/rocketblend/rocketblend/pkg/rocketblend/factory"

	"github.com/flowshot-io/x/pkg/logger"
)

type (
	Factory interface {
		GetLogger() (logger.Logger, error)
		GetApplicationConfigService() (config.Service, error)

		GetEventService() (eventservice.Service, error)

		GetSearchStore() (searchstore.Store, error)
		GetProjectService() (projectservice.Service, error)
		GetPackageService() (packageservice.Service, error)
		GetOperationService() (operationservice.Service, error)

		Preload() error
		Close() error

		rocketblendFactory.Factory
	}

	Options struct {
		Logger logger.Logger

		RocketBlendFactory      rocketblendFactory.Factory
		WatcherDebounceDuration time.Duration
	}

	Option func(*Options)

	factory struct {
		appDir string

		logger logger.Logger

		watcherDebounceDuration time.Duration

		rocketblendFactory rocketblendFactory.Factory

		closing      bool
		closingMutex sync.RWMutex

		applicationConfigMutex   sync.RWMutex
		applicationConfigService config.Service

		searchstoreMutex sync.RWMutex
		searchStore      searchstore.Store

		projectMutex   sync.RWMutex
		projectService projectservice.Service

		packageMutex   sync.RWMutex
		packageService packageservice.Service

		operationMutex   sync.RWMutex
		operationService operationservice.Service

		eventMutex   sync.RWMutex
		eventService eventservice.Service
	}
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithRocketBlendFactory(factory rocketblendFactory.Factory) Option {
	return func(o *Options) {
		o.RocketBlendFactory = factory
	}
}

func WithWatcherDebounceDuration(duration time.Duration) Option {
	return func(o *Options) {
		o.WatcherDebounceDuration = duration
	}
}

func New(opts ...Option) (Factory, error) {
	options := &Options{
		Logger:                  logger.NoOp(),
		WatcherDebounceDuration: 2 * time.Second,
	}

	for _, o := range opts {
		o(options)
	}

	if options.RocketBlendFactory == nil {
		factory, err := rocketblendFactory.New()
		if err != nil {
			return nil, fmt.Errorf("failed to create rocketblend factory: %w", err)
		}

		options.RocketBlendFactory = factory
	}

	options.RocketBlendFactory.SetLogger(options.Logger) // TODO: Add WithLogger.

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("cannot find config directory: %v", err)
	}

	appDir := filepath.Join(userConfigDir, build.AppName)
	if build.Version == "dev" {
		appDir = filepath.Join(appDir, "dev")
	}

	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create app directory: %w", err)
	}

	options.Logger.Info("Using app directory", map[string]interface{}{"appDir": appDir})

	return &factory{
		appDir:                  appDir,
		watcherDebounceDuration: options.WatcherDebounceDuration,
		logger:                  options.Logger,
		rocketblendFactory:      options.RocketBlendFactory,
	}, nil
}

func (f *factory) GetLogger() (logger.Logger, error) {
	return f.logger, nil
}

func (f *factory) SetLogger(logger.Logger) error {
	return fmt.Errorf("not implemented: will be removed in future")
}

func (f *factory) Preload() error {
	if err := f.checkClosing(); err != nil {
		return err
	}

	f.logger.Info("preloading factory")

	if _, err := f.GetProjectService(); err != nil {
		return err
	}

	if _, err := f.GetPackageService(); err != nil {
		return err
	}

	return nil
}

func (f *factory) Close() error {
	if err := f.setClosing(); err != nil {
		return err
	}

	f.logger.Info("closing factory")

	f.projectMutex.Lock()
	if f.projectService != nil {
		if err := f.projectService.Close(); err != nil {
			f.projectMutex.Unlock()
			return err
		}
		f.projectService = nil
	}
	f.projectMutex.Unlock()

	f.packageMutex.Lock()
	if f.packageService != nil {
		if err := f.packageService.Close(); err != nil {
			f.packageMutex.Unlock()
			return err
		}
		f.packageService = nil
	}
	f.packageMutex.Unlock()

	f.operationMutex.Lock()
	if f.operationService != nil {
		f.operationService = nil
	}
	f.operationMutex.Unlock()

	f.searchstoreMutex.Lock()
	if f.searchStore != nil {
		if err := f.searchStore.Close(); err != nil {
			f.searchstoreMutex.Unlock()
			return err
		}
		f.searchStore = nil
	}
	f.searchstoreMutex.Unlock()

	f.eventMutex.Lock()
	if f.eventService != nil {
		if err := f.eventService.Close(); err != nil {
			f.eventMutex.Unlock()
			return err
		}
		f.eventService = nil
	}
	f.eventMutex.Unlock()

	f.applicationConfigMutex.Lock()
	if f.applicationConfigService != nil {
		f.applicationConfigService = nil
	}
	f.applicationConfigMutex.Unlock()

	f.unsetClosing()

	f.logger.Info("factory cleared")

	return nil
}

func (f *factory) GetApplicationConfigService() (config.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	f.applicationConfigMutex.RLock()
	if f.applicationConfigService != nil {
		defer f.applicationConfigMutex.RUnlock()
		return f.applicationConfigService, nil
	}
	f.applicationConfigMutex.RUnlock()

	f.applicationConfigMutex.Lock()
	defer f.applicationConfigMutex.Unlock()
	configService, err := config.New(f.appDir)
	if err != nil {
		return nil, err
	}

	f.applicationConfigService = configService

	return f.applicationConfigService, nil
}

func (f *factory) GetEventService() (eventservice.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	f.eventMutex.RLock()
	if f.eventService != nil {
		defer f.eventMutex.RUnlock()
		return f.eventService, nil
	}
	f.eventMutex.RUnlock()

	f.eventMutex.Lock()
	defer f.eventMutex.Unlock()
	eventService, err := eventservice.New(
		eventservice.WithLogger(f.logger),
	)
	if err != nil {
		return nil, err
	}

	f.eventService = eventService

	return f.eventService, nil
}

func (f *factory) GetSearchStore() (searchstore.Store, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	f.searchstoreMutex.RLock()
	if f.searchStore != nil {
		defer f.searchstoreMutex.RUnlock()
		return f.searchStore, nil
	}
	f.searchstoreMutex.RUnlock()

	event, err := f.GetEventService()
	if err != nil {
		return nil, err
	}

	f.searchstoreMutex.Lock()
	defer f.searchstoreMutex.Unlock()
	store, err := searchstore.New(
		searchstore.WithLogger(f.logger),
		searchstore.WithDispatcherService(event),
	)
	if err != nil {
		return nil, err
	}

	f.searchStore = store

	return f.searchStore, nil
}

func (f *factory) GetProjectService() (projectservice.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	f.projectMutex.RLock()
	if f.projectService != nil {
		defer f.projectMutex.RUnlock()
		return f.projectService, nil
	}
	f.projectMutex.RUnlock()

	applicationConfigService, err := f.GetApplicationConfigService()
	if err != nil {
		return nil, err
	}

	rocketblendBlendFileService, err := f.GetBlendFileService()
	if err != nil {
		return nil, err
	}

	rocketblendInstallationService, err := f.GetInstallationService()
	if err != nil {
		return nil, err
	}

	rocketblendPackageService, err := f.GetRocketPackService()
	if err != nil {
		return nil, err
	}

	dispatcher, err := f.GetEventService()
	if err != nil {
		return nil, err
	}

	store, err := f.GetSearchStore()
	if err != nil {
		return nil, err
	}

	f.projectMutex.Lock()
	defer f.projectMutex.Unlock()
	projectService, err := projectservice.New(
		projectservice.WithLogger(f.logger),
		projectservice.WithWatcherDebounceDuration(f.watcherDebounceDuration),
		projectservice.WithApplicationConfigService(applicationConfigService),
		projectservice.WithRocketBlendBlendFileService(rocketblendBlendFileService),
		projectservice.WithRocketBlendInstallationService(rocketblendInstallationService),
		projectservice.WithRocketBlendPackageService(rocketblendPackageService),
		projectservice.WithStore(store),
		projectservice.WithDispatcher(dispatcher),
	)
	if err != nil {
		return nil, err
	}

	f.projectService = projectService

	return f.projectService, nil
}

func (f *factory) GetPackageService() (packageservice.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	f.packageMutex.RLock()
	if f.packageService != nil {
		defer f.packageMutex.RUnlock()
		return f.packageService, nil
	}
	f.packageMutex.RUnlock()

	rocketblendConfigService, err := f.GetConfigService()
	if err != nil {
		return nil, err
	}

	rocketblendPackageService, err := f.GetRocketPackService()
	if err != nil {
		return nil, err
	}

	rocketblendInstallationService, err := f.GetInstallationService()
	if err != nil {
		return nil, err
	}

	dispatcher, err := f.GetEventService()
	if err != nil {
		return nil, err
	}

	store, err := f.GetSearchStore()
	if err != nil {
		return nil, err
	}

	f.packageMutex.Lock()
	defer f.packageMutex.Unlock()
	packageService, err := packageservice.New(
		packageservice.WithLogger(f.logger),
		packageservice.WithWatcherDebounceDuration(f.watcherDebounceDuration),
		packageservice.WithRocketBlendPackageService(rocketblendPackageService),
		packageservice.WithRocketBlendInstallationService(rocketblendInstallationService),
		packageservice.WithRocketBlendConfigService(rocketblendConfigService),
		packageservice.WithStore(store),
		packageservice.WithDispatcher(dispatcher),
	)
	if err != nil {
		return nil, err
	}

	f.packageService = packageService

	return f.packageService, nil
}

func (f *factory) GetOperationService() (operationservice.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	f.operationMutex.RLock()
	if f.operationService != nil {
		defer f.operationMutex.RUnlock()
		return f.operationService, nil
	}
	f.operationMutex.RUnlock()

	dispatcher, err := f.GetEventService()
	if err != nil {
		return nil, err
	}

	store, err := f.GetSearchStore()
	if err != nil {
		return nil, err
	}

	f.operationMutex.Lock()
	defer f.operationMutex.Unlock()
	operationService, err := operationservice.New(
		operationservice.WithLogger(f.logger),
		operationservice.WithStore(store),
		operationservice.WithDispatcher(dispatcher),
	)
	if err != nil {
		return nil, err
	}

	f.operationService = operationService

	return f.operationService, nil
}

func (f *factory) GetConfigService() (rocketblendConfig.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	return f.rocketblendFactory.GetConfigService()
}

func (f *factory) GetRocketPackService() (rocketblendRocketPack.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	return f.rocketblendFactory.GetRocketPackService()
}

func (f *factory) GetInstallationService() (rocketblendInstallation.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	return f.rocketblendFactory.GetInstallationService()
}

func (f *factory) GetBlendFileService() (rocketblendBlendFile.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	return f.rocketblendFactory.GetBlendFileService()
}

func (f *factory) checkClosing() error {
	f.closingMutex.RLock()
	defer f.closingMutex.RUnlock()

	if f.closing {
		return errors.New("factory is closing or closed")
	}

	return nil
}

func (f *factory) setClosing() error {
	f.closingMutex.Lock()
	defer f.closingMutex.Unlock()

	if f.closing {
		return errors.New("factory is already closing")
	}

	f.closing = true
	return nil
}

func (f *factory) unsetClosing() {
	f.closingMutex.Lock()
	defer f.closingMutex.Unlock()

	f.closing = false
}
