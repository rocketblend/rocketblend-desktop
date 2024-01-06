package factory

import (
	"github.com/rocketblend/rocketblend-desktop/internal/application/config"
	"github.com/rocketblend/rocketblend-desktop/internal/application/packageservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectservice"
)

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
