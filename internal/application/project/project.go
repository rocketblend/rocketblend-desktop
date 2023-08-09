package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rocketblend/rocketblend-desktop/internal/application/projectsettings"
	"github.com/rocketblend/rocketblend/pkg/driver/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketfile"
)

// TODO: store data as in project projectservice. Add GetBlendFile() and GetSettings() methods.

const (
	IgnoreFileName = ".rbdesktopignore"
	ConfigDir      = ".rbdesktop"
)

type (
	Project struct {
		BlendFile *blendconfig.BlendConfig         `json:"blendFile,omitempty"`
		Settings  *projectsettings.ProjectSettings `json:"settings,omitempty"`
		UpdatedAt time.Time                        `json:"updatedAt,omitempty"`
	}
)

func Load(projectPath string) (*Project, error) {
	if ignoreProject(projectPath) {
		return nil, fmt.Errorf("project %s is ignored", projectPath)
	}

	files, err := os.ReadDir(projectPath)
	if err != nil {
		return nil, err
	}

	blendFilePath := ""
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), blendconfig.BlenderFileExtension) {
			blendFilePath = filepath.Join(projectPath, file.Name())
			break
		}
	}

	if blendFilePath == "" {
		return nil, fmt.Errorf("no blend file found in %s", projectPath)
	}

	blendConfig, err := blendconfig.Load(blendFilePath, filepath.Join(projectPath, rocketfile.FileName))
	if err != nil {
		return nil, err
	}

	settings, err := loadOrCreateSettings(filepath.Join(projectPath, ConfigDir, projectsettings.FileName))
	if err != nil {
		return nil, err
	}

	return &Project{
		BlendFile: blendConfig,
		Settings:  settings,
	}, nil
}

func Save(project *Project) error {
	return nil
}

func ignoreProject(projectPath string) bool {
	_, err := os.Stat(filepath.Join(projectPath, IgnoreFileName))
	return !os.IsNotExist(err)
}
