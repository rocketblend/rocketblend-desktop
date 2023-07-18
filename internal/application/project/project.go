package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rocketblend/rocketblend/pkg/driver/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketfile"
)

const (
	IgnoreFileName = ".rockdeskignore"
	SettingsFolder = ".rocketdesk"
)

type (
	Project struct {
		Key           string                 `json:"key"`
		Name          string                 `json:"name"`
		BlendFileName string                 `json:"blendFileName"`
		Config        *rocketfile.RocketFile `json:"rocketFile"`
		LastUpdated   time.Time              `json:"lastUpdated"`
	}
)

func Find(projectPath string) (*Project, error) {
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

	// Rocketfiles are always named rocketfile.yaml.
	rocketFilePath := filepath.Join(projectPath, rocketfile.FileName)

	// Here the name of the project is assumed to be the base of the projectPath.
	projectName := filepath.Base(projectPath)

	// Will validate the existence of the blend file and the rocket file.
	project, err := Load(projectName, blendFilePath, rocketFilePath)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func Load(projectName string, blendFilePath string, rocketFilePath string) (*Project, error) {
	blendConfig, err := blendconfig.Load(blendFilePath, rocketFilePath)
	if err != nil {
		return nil, err
	}

	blendFileStat, err := os.Stat(blendFilePath)
	if err != nil {
		return nil, err
	}

	rocketFileStat, err := os.Stat(rocketFilePath)
	if err != nil {
		return nil, err
	}

	lastUpdated := blendFileStat.ModTime()
	if rocketFileStat.ModTime().After(blendFileStat.ModTime()) {
		lastUpdated = rocketFileStat.ModTime()
	}

	// Create a new project instance with loaded BlendConfig.
	return &Project{
		Key:           filepath.Dir(blendFilePath),
		Name:          projectName,
		BlendFileName: blendConfig.BlendFileName,
		Config:        blendConfig.RocketFile,
		LastUpdated:   lastUpdated,
	}, nil
}
