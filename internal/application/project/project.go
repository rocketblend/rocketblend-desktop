package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/copystructure"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectsettings"
	"github.com/rocketblend/rocketblend/pkg/driver/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketfile"
)

const (
	IgnoreFileName = ".rbdesktopignore"
	ConfigDir      = ".rbdesktop"
)

type (
	Project struct {
		Key       string                           `json:"key"`
		BlendFile *blendconfig.BlendConfig         `json:"blendFile"`
		Settings  *projectsettings.ProjectSettings `json:"settings"`
		UpdatedAt time.Time                        `json:"updatedAt"`
	}
)

// Copy performs a deep copy of the Project.
func (p *Project) Copy() *Project {
	copiedProject, _ := copystructure.Copy(p)
	return copiedProject.(*Project)
}

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
		Key:       projectPath,
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
