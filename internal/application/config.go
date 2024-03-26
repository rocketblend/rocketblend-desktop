package application

import (
	"github.com/rocketblend/rocketblend-desktop/internal/application/build"
	"github.com/rocketblend/rocketblend-desktop/internal/application/config"
	rbconfig "github.com/rocketblend/rocketblend/pkg/rocketblend/config"
)

type (
	Feature struct {
		Addon     bool `json:"addon"`
		Terminal  bool `json:"terminal"`
		Developer bool `json:"developer"`
	}

	Details struct {
		Version               string `json:"version"`
		Platform              string `json:"platform"`
		InstallationPath      string `json:"installationPath"`
		PackagePath           string `json:"packagePath"`
		ApplicationConfigPath string `json:"applicationConfigPath"`
		RocketBlendConfigPath string `json:"rocketblendConfigPath"`
	}

	Preferences struct {
		WatchPaths string  `json:"watchPaths"`
		Feature    Feature `json:"feature"`
	}

	UpdatePreferencesOpts Preferences
)

func (d *Driver) GetPreferences() (*Preferences, error) {
	_, aConfig, err := d.getApplicationConfig()
	if err != nil {
		return nil, err
	}

	watchPath := ""
	if len(aConfig.Project.Paths) > 0 {
		watchPath = aConfig.Project.Paths[0]
	}

	return &Preferences{
		WatchPaths: watchPath,
		Feature: Feature{
			Addon:     aConfig.Feature.Addon,
			Terminal:  aConfig.Feature.Terminal,
			Developer: aConfig.Feature.Developer,
		},
	}, nil
}

func (d *Driver) UpdatePreferences(opts UpdatePreferencesOpts) error {
	aConfigService, err := d.factory.GetApplicationConfigService()
	if err != nil {
		return err
	}

	aConfig, err := aConfigService.Get()
	if err != nil {
		return err
	}

	aConfig.Project.Paths = []string{opts.WatchPaths}
	aConfig.Feature.Addon = opts.Feature.Addon
	aConfig.Feature.Terminal = opts.Feature.Terminal
	aConfig.Feature.Developer = opts.Feature.Developer

	if err := aConfigService.Save(aConfig); err != nil {
		return err
	}

	if err := d.refresh(); err != nil {
		return err
	}

	return nil
}

func (d *Driver) GetDetails() (*Details, error) {
	aConfigPath, _, err := d.getApplicationConfig()
	if err != nil {
		return nil, err
	}

	rbConfig, err := d.getRocketBlendConfig()
	if err != nil {
		return nil, err
	}

	return &Details{
		Version:               build.Version,
		Platform:              rbConfig.Platform.String(),
		InstallationPath:      rbConfig.InstallationsPath,
		PackagePath:           rbConfig.PackagesPath,
		ApplicationConfigPath: aConfigPath,
	}, nil
}

func (d *Driver) getApplicationConfig() (string, *config.Config, error) {
	configService, err := d.factory.GetApplicationConfigService()
	if err != nil {
		d.logger.Error("failed to get application config service", map[string]interface{}{"error": err.Error()})
		return "", nil, err
	}

	config, err := configService.Get()
	if err != nil {
		d.logger.Error("failed to get config", map[string]interface{}{"error": err.Error()})
		return "", nil, err
	}

	return configService.Path(), config, nil
}

func (d *Driver) getRocketBlendConfig() (*rbconfig.Config, error) {
	configService, err := d.factory.GetConfigService()
	if err != nil {
		d.logger.Error("failed to get rocketblend config service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	config, err := configService.Get()
	if err != nil {
		d.logger.Error("failed to get config", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return config, nil
}

func (d *Driver) refresh() error {
	projectService, err := d.factory.GetProjectService()
	if err != nil {
		return err
	}

	return projectService.Refresh(d.ctx)
}
