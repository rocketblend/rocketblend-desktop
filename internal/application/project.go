package application

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/listoption"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	"github.com/rocketblend/rocketblend-desktop/internal/helpers"
	"github.com/rocketblend/rocketblend/pkg/reference"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Project struct {
		logger types.Logger

		operator  types.Operator
		portfolio types.Portfolio

		configurator   types.Configurator
		rbConfigurator rbtypes.Configurator

		ctx context.Context
	}

	GetProjectOpts struct {
		ID uuid.UUID `json:"id"`
	}

	GetProjectResult struct {
		Project *types.Project `json:"project"`
	}

	ListProjectsOpts struct {
		Query string `json:"query"`
	}

	ListProjectsResult struct {
		Projects []*types.Project `json:"projects"`
	}

	CreateProjectOpts struct {
		Name string `json:"name"`
	}

	CreateProjectResult struct {
		OperationID uuid.UUID `json:"operationID"`
	}

	AddProjectPackageOpts struct {
		ID        uuid.UUID           `json:"id"`
		Reference reference.Reference `json:"reference"`
	}

	RemoveProjectPackageOpts struct {
		ID        uuid.UUID           `json:"id"`
		Reference reference.Reference `json:"reference"`
	}

	UpdateProjectOpts struct {
		ID   uuid.UUID `json:"id"`
		Name *string   `json:"name,omitempty"`
	}

	DeleteProjectOpts struct {
		ID uuid.UUID `json:"id"`
	}

	RunProjectOpts struct {
		ID uuid.UUID `json:"id"`
	}

	RenderProjectOpts struct {
		ID uuid.UUID `json:"id"`
	}
)

func newProject(ctx context.Context, deps *dependencies) *Project {
	return &Project{
		logger:         deps.logger,
		operator:       deps.operator,
		portfolio:      deps.portfolio,
		configurator:   deps.configurator,
		rbConfigurator: deps.rbConfigurator,
		ctx:            ctx,
	}
}

func (s *Project) GetProject(opts GetPackageOpts) (*GetProjectResult, error) {
	project, err := s.portfolio.GetProject(s.ctx, &types.GetProjectOpts{
		ID: opts.ID,
	})
	if err != nil {
		s.logger.Error("failed to find project by id", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return &GetProjectResult{
		Project: project.Project,
	}, nil
}

func (s *Project) ListProjects(opts ListProjectsOpts) (*ListProjectsResult, error) {
	response, err := s.portfolio.ListProjects(s.ctx, listoption.WithQuery(opts.Query))
	if err != nil {
		s.logger.Error("failed to list projects", map[string]interface{}{
			"error": err.Error(),
			"query": opts.Query,
		})
		return nil, err
	}

	s.logger.Debug("found projects", map[string]interface{}{
		"total": len(response.Projects),
	})

	return &ListProjectsResult{
		Projects: response.Projects,
	}, nil
}

func (s *Project) CreateProject(opts CreateProjectOpts) (*CreateProjectResult, error) {
	projectPath, err := s.getProjectPath()
	if err != nil {
		return nil, err
	}

	defaultBuild, err := s.getDefaultBuild()
	if err != nil {
		return nil, err
	}

	fileName := helpers.DisplayNameToFilename(opts.Name)
	opid, err := s.operator.Create(s.ctx, func(ctx context.Context, opid uuid.UUID) (interface{}, error) {
		result, err := s.portfolio.CreateProject(ctx, &types.CreateProjectOpts{
			DisplayName:   opts.Name,
			BlendFileName: fileName + rbtypes.BlendFileExtension,
			Path:          filepath.Join(projectPath, fileName),
			Build:         defaultBuild,
		})
		if err != nil {
			s.logger.Error("failed to create project", map[string]interface{}{
				"error": err.Error(),
				"opid":  opid,
			})
			return nil, err
		}

		s.logger.Debug("project created", map[string]interface{}{
			"id":   result.ID,
			"opid": opid,
		})

		return result, nil
	})
	if err != nil {
		return nil, err
	}

	return &CreateProjectResult{
		OperationID: opid,
	}, nil
}

func (s *Project) AddProjectPackage(opts AddProjectPackageOpts) error {
	if err := s.portfolio.AddProjectPackage(s.ctx, &types.AddProjectPackageOpts{
		ID:        opts.ID,
		Reference: opts.Reference,
	}); err != nil {
		s.logger.Error("failed to add package to project", map[string]interface{}{
			"error":     err.Error(),
			"id":        opts.ID,
			"reference": opts.Reference,
		})
		return err
	}

	s.logger.Debug("package added to project", map[string]interface{}{
		"id":        opts.ID,
		"reference": opts.Reference,
	})

	return nil
}

func (s *Project) RemoveProjectPackage(opts RemoveProjectPackageOpts) error {
	if err := s.portfolio.RemoveProjectPackage(s.ctx, &types.RemoveProjectPackageOpts{
		ID:        opts.ID,
		Reference: opts.Reference,
	}); err != nil {
		s.logger.Error("failed to remove package from project", map[string]interface{}{
			"error":     err.Error(),
			"id":        opts.ID,
			"reference": opts.Reference,
		})
		return err
	}

	s.logger.Debug("package removed from project", map[string]interface{}{
		"id":        opts.ID,
		"reference": opts.Reference,
	})

	return nil
}

// UpdateProject updates a project
func (s *Project) UpdateProject(opts UpdateProjectOpts) error {
	if err := s.portfolio.UpdateProject(s.ctx, &types.UpdateProjectOpts{
		ID:   opts.ID,
		Name: opts.Name,
	}); err != nil {
		s.logger.Error("failed to update project", map[string]interface{}{
			"error": err.Error(),
			"id":    opts.ID,
			"name":  opts.Name,
		})
		return err
	}

	s.logger.Debug("project updated", map[string]interface{}{
		"id": opts.ID,
	})

	return nil
}

func (s *Project) DeleteProject(opts DeleteProjectOpts) error {
	s.logger.Debug("deleting project", map[string]interface{}{"id": opts.ID})
	return types.ErrNotImplement
}

func (s *Project) RunProject(opts RunProjectOpts) error {
	// Blender is started with this context and must outlive the app.
	ctx := context.WithoutCancel(s.ctx)

	if err := s.portfolio.RunProject(ctx, &types.RunProjectOpts{
		ID: opts.ID,
	}); err != nil {
		s.logger.Error("failed to run project", map[string]interface{}{
			"error": err.Error(),
			"id":    opts.ID,
		})
		return err
	}

	s.logger.Debug("project started", map[string]interface{}{
		"id": opts.ID,
	})

	return nil
}

func (s *Project) RenderProject(opts RenderProjectOpts) error {
	s.logger.Debug("project rendered", map[string]interface{}{"id": opts.ID})
	return types.ErrNotImplement
}

func (s *Project) getDefaultBuild() (reference.Reference, error) {
	config, err := s.rbConfigurator.Get()
	if err != nil {
		return "", err
	}

	return config.DefaultBuild, nil
}

func (s *Project) getProjectPath() (string, error) {
	config, err := s.configurator.Get()
	if err != nil {
		return "", err
	}

	if len(config.Project.Paths) == 0 {
		return "", errors.New("no project path configured")
	}

	// TODO: support multiple project paths
	return config.Project.Paths[0], nil
}
