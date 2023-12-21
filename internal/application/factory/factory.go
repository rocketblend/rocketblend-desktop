package factory

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/rocketblend/rocketblend-desktop/internal/application/build"
	"github.com/rocketblend/rocketblend-desktop/internal/application/config"
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

		GetSearchStore() (searchstore.Store, error)
		GetProjectService() (projectservice.Service, error)
		GetPackageService() (packageservice.Service, error)

		Preload() error
		Close() error

		rocketblendFactory.Factory
	}

	Options struct {
		Logger             logger.Logger
		RocketBlendFactory rocketblendFactory.Factory
	}

	Option func(*Options)

	factory struct {
		appDir string

		logger logger.Logger

		rocketblendFactory rocketblendFactory.Factory

		applicationConfigMutex   sync.Mutex
		applicationConfigService config.Service

		searchstoreMutex sync.Mutex
		searchStore      searchstore.Store

		projectMutex   sync.Mutex
		projectService projectservice.Service

		packageMutex   sync.Mutex
		packageService packageservice.Service
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

func New(opts ...Option) (Factory, error) {
	options := &Options{
		Logger: logger.NoOp(),
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
		appDir:             appDir,
		logger:             options.Logger,
		rocketblendFactory: options.RocketBlendFactory,
	}, nil
}

func (f *factory) GetLogger() (logger.Logger, error) {
	return f.logger, nil
}

func (f *factory) SetLogger(logger.Logger) error {
	return fmt.Errorf("not implemented: will be removed in future")
}

func (f *factory) Preload() error {
	if _, err := f.GetProjectService(); err != nil {
		return err
	}

	if _, err := f.GetPackageService(); err != nil {
		return err
	}

	return nil
}

func (f *factory) Close() error {
	if f.packageService != nil {
		if err := f.packageService.Close(); err != nil {
			return err
		}
	}

	if f.projectService != nil {
		if err := f.projectService.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (f *factory) GetApplicationConfigService() (config.Service, error) {
	f.applicationConfigMutex.Lock()
	defer f.applicationConfigMutex.Unlock()

	if f.applicationConfigService != nil {
		return f.applicationConfigService, nil
	}

	configService, err := config.New(f.appDir)
	if err != nil {
		return nil, err
	}

	f.applicationConfigService = configService

	return f.applicationConfigService, nil
}

func (f *factory) GetSearchStore() (searchstore.Store, error) {
	f.searchstoreMutex.Lock()
	defer f.searchstoreMutex.Unlock()

	if f.searchStore != nil {
		return f.searchStore, nil
	}

	store, err := searchstore.New(
		searchstore.WithLogger(f.logger),
	)
	if err != nil {
		return nil, err
	}

	f.searchStore = store

	return f.searchStore, nil
}

func (f *factory) GetProjectService() (projectservice.Service, error) {
	f.projectMutex.Lock()
	defer f.projectMutex.Unlock()

	if f.projectService != nil {
		return f.projectService, nil
	}

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

	store, err := f.GetSearchStore()
	if err != nil {
		return nil, err
	}

	projectService, err := projectservice.New(
		projectservice.WithLogger(f.logger),
		projectservice.WithApplicationConfigService(applicationConfigService),
		projectservice.WithRocketBlendBlendFileService(rocketblendBlendFileService),
		projectservice.WithRocketBlendInstallationService(rocketblendInstallationService),
		projectservice.WithRocketBlendPackageService(rocketblendPackageService),
		projectservice.WithStore(store),
	)
	if err != nil {
		return nil, err
	}

	f.projectService = projectService

	return f.projectService, nil
}

func (f *factory) GetPackageService() (packageservice.Service, error) {
	f.packageMutex.Lock()
	defer f.packageMutex.Unlock()

	if f.packageService != nil {
		return f.packageService, nil
	}

	rocketblendConfigService, err := f.GetConfigService()
	if err != nil {
		return nil, err
	}

	store, err := f.GetSearchStore()
	if err != nil {
		return nil, err
	}

	packageService, err := packageservice.New(
		packageservice.WithLogger(f.logger),
		packageservice.WithRocketBlendConfigService(rocketblendConfigService),
		packageservice.WithStore(store),
	)
	if err != nil {
		return nil, err
	}

	f.packageService = packageService

	return f.packageService, nil
}

func (f *factory) GetConfigService() (rocketblendConfig.Service, error) {
	return f.rocketblendFactory.GetConfigService()
}

func (f *factory) GetRocketPackService() (rocketblendRocketPack.Service, error) {
	return f.rocketblendFactory.GetRocketPackService()
}

func (f *factory) GetInstallationService() (rocketblendInstallation.Service, error) {
	return f.rocketblendFactory.GetInstallationService()
}

func (f *factory) GetBlendFileService() (rocketblendBlendFile.Service, error) {
	return f.rocketblendFactory.GetBlendFileService()
}
