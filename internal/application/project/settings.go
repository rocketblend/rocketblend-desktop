package project

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectsettings"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func loadOrCreateSettings(filePath string) (*projectsettings.ProjectSettings, error) {
	settings, err := projectsettings.Load(filePath)
	if err == nil {
		return settings, nil
	}

	// Create a default settings if loading failed
	defaultSettings := &projectsettings.ProjectSettings{
		ID:   uuid.New(),
		Name: filenameToDisplayName(filepath.Dir(filePath)),
	}

	if err := projectsettings.Save(defaultSettings, filePath); err != nil {
		return nil, fmt.Errorf("failed to save default settings: %s", err)
	}

	return defaultSettings, nil
}

// Converts a filename into a more human readable display name
func filenameToDisplayName(filename string) string {
	name := filepath.Base(filename)
	ext := filepath.Ext(name)

	// Remove the extension
	name = name[0 : len(name)-len(ext)]

	// Replace underscores or hyphens with spaces
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")

	// Capitalize the first letter of each word
	title := cases.Title(language.English)
	name = title.String(name)

	return name
}
