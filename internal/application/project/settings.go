package project

import (
	"fmt"
	"path/filepath"

	"github.com/rocketblend/rocketblend-desktop/internal/application/projectsettings"
)

func loadOrCreateSettings(filePath string) (*projectsettings.ProjectSettings, error) {
	settings, err := projectsettings.Load(filePath)
	if err == nil {
		return settings, nil
	}

	// Create a default settings if loading failed
	defaultSettings := &projectsettings.ProjectSettings{
		Name: filepath.Base(filepath.Dir(filepath.Dir(filePath))),
	}

	if err := projectsettings.Save(defaultSettings, filePath); err != nil {
		return nil, fmt.Errorf("failed to save default settings: %s", err)
	}

	return defaultSettings, nil
}
