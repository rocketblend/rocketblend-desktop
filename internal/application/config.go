package application

import (
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
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

func (s *System) GetPreferences() (*Preferences, error) {
	_, aConfig, err := s.getApplicationConfig()
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

func (s *System) UpdatePreferences(opts UpdatePreferencesOpts) error {
	config, err := s.configurator.Get()
	if err != nil {
		return err
	}

	config.Project.Paths = []string{opts.WatchPath}
	config.Feature.Addon = opts.Feature.Addon
	config.Feature.Developer = opts.Feature.Developer

	if err := s.configurator.Save(config); err != nil {
		return err
	}

	if err := s.portfolio.Refresh(s.ctx); err != nil {
		return err
	}

	return nil
}

func (s *System) GetDetails() (*Details, error) {
	aConfigPath, _, err := s.getApplicationConfig()
	if err != nil {
		return nil, err
	}

	rbConfigPath, rbConfig, err := s.getRocketBlendConfig()
	if err != nil {
		return nil, err
	}

	return &Details{
		Version:               s.version,
		Platform:              s.platform.String(), // TODO: Convert to wails enum.
		InstallationPath:      rbConfig.InstallationsPath,
		PackagePath:           rbConfig.PackagesPath,
		ApplicationConfigPath: aConfigPath,
		RocketBlendConfigPath: rbConfigPath,
	}, nil
}

func (s *System) getApplicationConfig() (string, *types.Config, error) {
	// TODO: create new struct that just includes the path and the config.
	config, err := s.configurator.Get()
	if err != nil {
		return "", nil, err
	}

	return s.configurator.Path(), config, nil
}

func (s *System) getRocketBlendConfig() (string, *rbtypes.Config, error) {
	// TODO: create new struct that just includes the path and the config.
	config, err := s.rbConfigurator.Get()
	if err != nil {
		return "", nil, err
	}

	return s.rbConfigurator.Path(), config, nil
}
