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

func (d *Driver) GetProject(opts GetPackageOpts) (*GetProjectResult, error) {
	ctx := context.Background()

	project, err := d.portfolio.GetProject(ctx, &types.GetProjectOpts{
		ID: opts.ID,
	})
	if err != nil {
		d.logger.Error("failed to find project by id", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return &GetProjectResult{
		Project: project.Project,
	}, nil
}

func (d *Driver) ListProjects(opts ListProjectsOpts) (*ListProjectsResult, error) {
	ctx := context.Background()

	response, err := d.portfolio.ListProjects(ctx, listoption.WithQuery(opts.Query))
	if err != nil {
		d.logger.Error("failed to list projects", map[string]interface{}{
			"error": err.Error(),
			"query": opts.Query,
		})
		return nil, err
	}

	d.logger.Debug("found projects", map[string]interface{}{
		"total": len(response.Projects),
	})

	return &ListProjectsResult{
		Projects: response.Projects,
	}, nil
}

func (d *Driver) CreateProject(opts CreateProjectOpts) (*CreateProjectResult, error) {
	projectPath, err := d.getProjectPath()
	if err != nil {
		return nil, err
	}

	defaultBuild, err := d.getDefaultBuild()
	if err != nil {
		return nil, err
	}

	fileName := helpers.DisplayNameToFilename(opts.Name)
	opid, err := d.operator.Create(d.ctx, func(ctx context.Context, opid uuid.UUID) (interface{}, error) {
		result, err := d.portfolio.CreateProject(ctx, &types.CreateProjectOpts{
			DisplayName:   opts.Name,
			BlendFileName: fileName + rbtypes.BlendFileExtension,
			Path:          filepath.Join(projectPath, fileName),
			Build:         defaultBuild,
		})
		if err != nil {
			d.logger.Error("failed to create project", map[string]interface{}{
				"error": err.Error(),
				"opid":  opid,
			})
			return nil, err
		}

		d.logger.Debug("project created", map[string]interface{}{
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

func (d *Driver) AddProjectPackage(opts AddProjectPackageOpts) error {
	if err := d.portfolio.AddProjectPackage(d.ctx, &types.AddProjectPackageOpts{
		ID:        opts.ID,
		Reference: opts.Reference,
	}); err != nil {
		d.logger.Error("failed to add package to project", map[string]interface{}{
			"error":     err.Error(),
			"id":        opts.ID,
			"reference": opts.Reference,
		})
		return err
	}

	d.logger.Debug("package added to project", map[string]interface{}{
		"id":        opts.ID,
		"reference": opts.Reference,
	})

	return nil
}

func (d *Driver) RemoveProjectPackage(opts RemoveProjectPackageOpts) error {
	if err := d.portfolio.RemoveProjectPackage(d.ctx, &types.RemoveProjectPackageOpts{
		ID:        opts.ID,
		Reference: opts.Reference,
	}); err != nil {
		d.logger.Error("failed to remove package from project", map[string]interface{}{
			"error":     err.Error(),
			"id":        opts.ID,
			"reference": opts.Reference,
		})
		return err
	}

	d.logger.Debug("package removed from project", map[string]interface{}{
		"id":        opts.ID,
		"reference": opts.Reference,
	})

	return nil
}

// UpdateProject updates a project
func (d *Driver) UpdateProject(opts UpdateProjectOpts) error {
	if err := d.portfolio.UpdateProject(d.ctx, &types.UpdateProjectOpts{
		ID:   opts.ID,
		Name: opts.Name,
	}); err != nil {
		d.logger.Error("failed to update project", map[string]interface{}{
			"error": err.Error(),
			"id":    opts.ID,
			"name":  opts.Name,
		})
		return err
	}

	d.logger.Debug("project updated", map[string]interface{}{
		"id": opts.ID,
	})

	return nil
}

func (d *Driver) DeleteProject(opts DeleteProjectOpts) error {
	d.logger.Debug("deleting project", map[string]interface{}{"id": opts.ID})
	return types.ErrNotImplement
}

func (d *Driver) RunProject(opts RunProjectOpts) error {
	ctx := context.Background()

	if err := d.portfolio.RunProject(ctx, &types.RunProjectOpts{
		ID: opts.ID,
	}); err != nil {
		d.logger.Error("failed to run project", map[string]interface{}{
			"error": err.Error(),
			"id":    opts.ID,
		})
		return err
	}

	d.logger.Debug("project started", map[string]interface{}{
		"id": opts.ID,
	})

	return nil
}

func (d *Driver) RenderProject(opts RenderProjectOpts) error {
	d.logger.Debug("project rendered", map[string]interface{}{"id": opts.ID})
	return types.ErrNotImplement
}

func (d *Driver) getDefaultBuild() (reference.Reference, error) {
	config, err := d.rbConfigurator.Get()
	if err != nil {
		return "", err
	}

	return config.DefaultBuild, nil
}

func (d *Driver) getProjectPath() (string, error) {
	config, err := d.configurator.Get()
	if err != nil {
		return "", err
	}

	if len(config.Project.Paths) == 0 {
		return "", errors.New("no project path configured")
	}

	// TODO: support multiple project paths
	return config.Project.Paths[0], nil
}
