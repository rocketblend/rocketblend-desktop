package projectservice

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectsettings"
	"github.com/rocketblend/rocketblend-desktop/internal/application/util"
	"github.com/rocketblend/rocketblend/pkg/driver"
	"github.com/rocketblend/rocketblend/pkg/driver/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketfile"
)

type (
	CreateProjectOpts struct {
		DisplayName   string              `json:"name"`
		BlendFileName string              `json:"blendFileName"`
		Path          string              `json:"path"`
		Build         reference.Reference `json:"build"`
	}

	CreateProjectResult struct {
		ID uuid.UUID `json:"id"`
	}
)

func (s *service) Create(ctx context.Context, opts CreateProjectOpts) (*CreateProjectResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	id, err := s.create(ctx, opts)
	if err != nil {
		return nil, err
	}

	s.EmitEvent(ctx, id, CreateEventChannel)

	return &CreateProjectResult{
		ID: id,
	}, nil
}

func (s *service) Render(ctx context.Context, id uuid.UUID) error {
	project, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	driver, err := s.createDriver(project.GetBlendFile())
	if err != nil {
		return err
	}

	if err := driver.Render(ctx); err != nil {
		return err
	}

	return nil
}

func (s *service) Run(ctx context.Context, id uuid.UUID) error {
	project, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	driver, err := s.createDriver(project.GetBlendFile())
	if err != nil {
		return err
	}

	go func() {
		if err := driver.Run(ctx); err != nil {
			s.logger.Error("failed to run project", map[string]interface{}{"error": err})
		}
	}()

	s.EmitEvent(ctx, id, RunEventChannel)

	return nil
}

func (s *service) Explore(ctx context.Context, id uuid.UUID) error {
	project, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	go func() {
		if err := util.Explore(ctx, project.Path); err != nil {
			s.logger.Debug("failed to open project in file explorer", map[string]interface{}{"error": err})
		}
	}()

	s.EmitEvent(ctx, id, ExploreEventChannel)

	return nil
}

func (s *service) EmitEvent(ctx context.Context, id uuid.UUID, channel string) {
	event := NewEvent(id)
	if err := s.dispatcher.EmitEvent(ctx, channel, event); err != nil {
		s.logger.Error("error emitting event", map[string]interface{}{
			"error":   err,
			"event":   event,
			"channel": channel,
		})
	}
}

func (s *service) create(ctx context.Context, opts CreateProjectOpts) (uuid.UUID, error) {
	blendConfig, err := blendconfig.New(
		opts.Path,
		ensureBlendExtension(opts.BlendFileName),
		rocketfile.New(opts.Build),
	)
	if err != nil {
		return uuid.Nil, err
	}

	driver, err := s.createDriver(blendConfig)
	if err != nil {
		return uuid.Nil, err
	}

	// We create a temporary ignore file to avoid adding the project to the index before it is fully created.
	if err := createIgnoreFile(opts.Path); err != nil {
		return uuid.Nil, err
	}

	defer func() {
		err := removeIgnoreFile(opts.Path)
		if err != nil {
			s.logger.Error("failed to remove temporarily project ignore file", map[string]interface{}{"error": err})
		}
	}()

	if err := driver.Create(ctx); err != nil {
		return uuid.Nil, err
	}

	id := uuid.New()
	if err := projectsettings.Save(&projectsettings.ProjectSettings{
		ID:   id,
		Name: opts.DisplayName,
	}, filepath.Join(opts.Path, projectsettings.FileName)); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *service) createDriver(blendConfig *blendconfig.BlendConfig) (driver.Driver, error) {
	return driver.New(
		driver.WithLogger(s.logger),
		driver.WithRocketPackService(s.rocketblendPackageService),
		driver.WithInstallationService(s.rocketblendInstallationService),
		driver.WithBlendFileService(s.rocketblendBlendFileService),
		driver.WithBlendConfig(blendConfig),
	)
}

func createIgnoreFile(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(path, project.IgnoreFileName))
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func removeIgnoreFile(path string) error {
	return os.Remove(filepath.Join(path, project.IgnoreFileName))
}

func ensureBlendExtension(filename string) string {
	if filename == "" {
		return "untitled.blend"
	}

	if !strings.HasSuffix(filename, ".blend") {
		filename += ".blend"
	}

	return filename
}
