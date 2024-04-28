package driver

import (
	"github.com/rocketblend/rocketblend-desktop/internal/application/build"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Feature struct {
		Addon     bool `json:"addon"`
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
		WatchPath string  `json:"watchPath"`
		Feature   Feature `json:"feature"`
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
		WatchPath: watchPath,
		Feature: Feature{
			Addon:     aConfig.Feature.Addon,
			Developer: aConfig.Feature.Developer,
		},
	}, nil
}

func (d *Driver) UpdatePreferences(opts UpdatePreferencesOpts) error {
	config, err := d.configurator.Get()
	if err != nil {
		return err
	}

	config.Project.Paths = []string{opts.WatchPath}
	config.Feature.Addon = opts.Feature.Addon
	config.Feature.Developer = opts.Feature.Developer

	if err := d.configurator.Save(config); err != nil {
		return err
	}

	if err := d.portfolio.Refresh(d.ctx); err != nil {
		return err
	}

	return nil
}

func (d *Driver) GetDetails() (*Details, error) {
	aConfigPath, _, err := d.getApplicationConfig()
	if err != nil {
		return nil, err
	}

	rbConfigPath, rbConfig, err := d.getRocketBlendConfig()
	if err != nil {
		return nil, err
	}

	return &Details{
		Version:               build.Version,
		Platform:              d.platform.String(), // TODO: Convert to wails enum.
		InstallationPath:      rbConfig.InstallationsPath,
		PackagePath:           rbConfig.PackagesPath,
		ApplicationConfigPath: aConfigPath,
		RocketBlendConfigPath: rbConfigPath,
	}, nil
}

func (d *Driver) getApplicationConfig() (string, *types.Config, error) {
	// TODO: create new struct that just includes the path and the config.
	config, err := d.configurator.Get()
	if err != nil {
		return "", nil, err
	}

	return d.configurator.Path(), config, nil
}

func (d *Driver) getRocketBlendConfig() (string, *rbtypes.Config, error) {
	// TODO: create new struct that just includes the path and the config.
	config, err := d.rbConfigurator.Get()
	if err != nil {
		return "", nil, err
	}

	return d.rbConfigurator.Path(), config, nil
}
