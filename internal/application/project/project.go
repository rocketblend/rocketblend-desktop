package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectsettings"
	"github.com/rocketblend/rocketblend-desktop/internal/application/util"
	"github.com/rocketblend/rocketblend/pkg/driver/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketfile"
)

const IgnoreFileName = ".rocketignore"

type (
	Project struct {
		ID            uuid.UUID             `json:"id"`
		Name          string                `json:"name,omitempty"`
		Tags          []string              `json:"tags,omitempty"`
		Path          string                `json:"path,omitempty"`
		FileName      string                `json:"fileName,omitempty"`
		Build         reference.Reference   `json:"build,omitempty"`
		Addons        []reference.Reference `json:"addons,omitempty"`
		SplashPath    string                `json:"splashPath,omitempty"`
		ThumbnailPath string                `json:"thumbnailPath,omitempty"`
		Version       string                `json:"version,omitempty"`
		UpdatedAt     time.Time             `json:"updatedAt,omitempty"`
	}
)

func (p *Project) GetBlendFile() *blendconfig.BlendConfig {
	return &blendconfig.BlendConfig{
		ProjectPath:   p.Path,
		BlendFileName: p.FileName,
		RocketFile: rocketfile.New(
			p.Build,
			p.Addons...,
		),
	}
}

func (p *Project) GetSettings() *projectsettings.ProjectSettings {
	return &projectsettings.ProjectSettings{
		ID:   p.ID,
		Name: p.Name,
		Tags: p.Tags,
	}
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

	blendFile, err := blendconfig.Load(blendFilePath, filepath.Join(projectPath, rocketfile.FileName))
	if err != nil {
		return nil, err
	}

	settings, err := loadOrCreateSettings(filepath.Join(projectPath, projectsettings.FileName))
	if err != nil {
		return nil, err
	}

	modTime, err := util.GetModTime(projectPath)
	if err != nil {
		return nil, err
	}

	thumbnailPath := ""
	if settings.ThumbnailPath != "" {
		if filepath.IsAbs(settings.ThumbnailPath) {
			return nil, fmt.Errorf("thumbnail file path must be relative: %s", settings.ThumbnailPath)
		}

		thumbnailPath = filepath.ToSlash(filepath.Join(projectPath, settings.ThumbnailPath))
	}

	splashPath := ""
	if settings.SplashPath != "" {
		if filepath.IsAbs(settings.SplashPath) {
			return nil, fmt.Errorf("splash file path must be relative: %s", settings.SplashPath)
		}

		splashPath = filepath.ToSlash(filepath.Join(projectPath, settings.SplashPath))
	}

	return &Project{
		ID:            settings.ID,
		Name:          settings.Name,
		Tags:          settings.Tags,
		Path:          blendFile.ProjectPath,
		FileName:      blendFile.BlendFileName, //TODO: Use full path.
		ThumbnailPath: thumbnailPath,
		SplashPath:    splashPath,
		Build:         blendFile.RocketFile.GetBuild(),
		Addons:        blendFile.RocketFile.GetAddons(),
		Version:       blendFile.RocketFile.GetVersion(),
		UpdatedAt:     modTime,
	}, nil
}

func Save(project *Project) error {
	return nil
}

func ignoreProject(projectPath string) bool {
	_, err := os.Stat(filepath.Join(projectPath, IgnoreFileName))
	return !os.IsNotExist(err)
}
