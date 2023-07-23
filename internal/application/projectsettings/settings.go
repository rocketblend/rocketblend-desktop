package projectsettings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/driver/helpers"
	"sigs.k8s.io/yaml"
)

const (
	FileName = "settings.yaml"
)

type (
	ThumbnailSettings struct {
		Width      int    `json:"width,omitempty"`
		Height     int    `json:"height,omitempty"`
		StartFrame int    `json:"startFrame,omitempty"`
		EndFrame   int    `json:"endFrame,omitempty"`
		RenderType string `json:"renderType,omitempty"`
	}

	ProjectSettings struct {
		Name              string             `json:"name,omitempty"`
		Tags              []string           `json:"tags,omitempty"`
		ThumbnailSettings *ThumbnailSettings `json:"thumbnailSettings,omitempty"`
		ThumbnailFilePath string             `json:"thumbnailFilePath,omitempty"`
	}
)

func Load(filePath string) (*ProjectSettings, error) {
	if err := validateFilePath(filePath); err != nil {
		return nil, fmt.Errorf("failed to validate file path: %s", err)
	}

	if err := helpers.FileExists(filePath); err != nil {
		return nil, fmt.Errorf("failed to find file: %s", err)
	}

	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	var settings ProjectSettings
	if err := yaml.Unmarshal(f, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rocketfile: %s", err)
	}

	if err := Validate(&settings); err != nil {
		return nil, fmt.Errorf("failed to validate rocketfile: %s", err)
	}

	return &settings, nil
}

func Save(settings *ProjectSettings, filePath string) error {
	if err := Validate(settings); err != nil {
		return fmt.Errorf("failed to validate project settings: %s", err)
	}

	if err := validateFilePath(filePath); err != nil {
		return fmt.Errorf("failed to validate file path: %s", err)
	}

	data, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal project settings: %s", err)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %s", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write project settings: %s", err)
	}

	return nil
}

func Validate(settings *ProjectSettings) error {
	return nil
}

func validateFilePath(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if filepath.Base(filePath) != FileName {
		return fmt.Errorf("invalid file name (must be '%s'): %s", FileName, filepath.Base(filePath))
	}

	return nil
}