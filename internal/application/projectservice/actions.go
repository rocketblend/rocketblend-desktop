package projectservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend/pkg/driver"
)

func (s *service) Render(ctx context.Context, id uuid.UUID) error {
	project, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	driver, err := s.createDriver(project)
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

	driver, err := s.createDriver(project)
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
		if err := openInFileExplorer(ctx, project.Path); err != nil {
			s.logger.Debug("failed to open project in file explorer", map[string]interface{}{"error": err})
		}
	}()

	s.EmitEvent(ctx, id, RunEventChannel)

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

func (s *service) createDriver(project *project.Project) (driver.Driver, error) {
	return driver.New(
		driver.WithLogger(s.logger),
		driver.WithRocketPackService(s.rocketblendPackageService),
		driver.WithInstallationService(s.rocketblendInstallationService),
		driver.WithBlendFileService(s.rocketblendBlendFileService),
		driver.WithBlendConfig(project.GetBlendFile()),
	)
}
