package projectservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend/pkg/driver"
)

func (s *service) Render(ctx context.Context, id uuid.UUID) error {
	project, err := s.store.GetProject(id)
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
	project, err := s.store.GetProject(id)
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

	return nil
}

func (s *service) Explore(ctx context.Context, id uuid.UUID) error {
	project, err := s.store.GetProject(id)
	if err != nil {
		return err
	}

	go func() {
		if err := openInFileExplorer(ctx, project.Path); err != nil {
			s.logger.Debug("failed to open project in file explorer", map[string]interface{}{"error": err})
		}
	}()

	return nil
}

func (s *service) createDriver(project *project.Project) (driver.Driver, error) {
	rocketPackService, err := s.factory.GetRocketPackService()
	if err != nil {
		return nil, err
	}

	installationService, err := s.factory.GetInstallationService()
	if err != nil {
		return nil, err
	}

	blendFileService, err := s.factory.GetBlendFileService()
	if err != nil {
		return nil, err
	}

	return driver.New(
		driver.WithLogger(s.logger),
		driver.WithRocketPackService(rocketPackService),
		driver.WithInstallationService(installationService),
		driver.WithBlendFileService(blendFileService),
		driver.WithBlendConfig(project.GetBlendFile()),
	)
}
