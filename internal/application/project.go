package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
	"github.com/rocketblend/rocketblend/pkg/reference"
)

type (
	Selected struct {
		ID uuid.UUID `json:"id"`
	}

	GetSelectedOpts struct {
		ID uuid.UUID `json:"id"`
	}

	UpdateSelectedOpts struct {
		ID uuid.UUID `json:"id"`
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
)

// GetProject gets a project by id
func (d *Driver) GetProject(id uuid.UUID) (*projectservice.GetProjectResponse, error) {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	project, err := projectService.Get(ctx, id)
	if err != nil {
		d.logger.Error("failed to find project by id", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return project, nil
}

// ListProjects lists all projects
func (d *Driver) ListProjects(query string) (*projectservice.ListProjectsResponse, error) {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	response, err := projectService.List(ctx, listoption.WithQuery(query))
	if err != nil {
		d.logger.Error("failed to find all projects", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	d.logger.Debug("found projects", map[string]interface{}{"projects": len(response.Projects)})

	return response, nil
}

// AddProjectPackage adds a package to a project
func (d *Driver) AddProjectPackage(opts AddProjectPackageOpts) error {
	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.AddPackage(d.ctx, &projectservice.AddProjectPackageOpts{
		ID:        opts.ID,
		Reference: opts.Reference,
	}); err != nil {
		d.logger.Error("failed to add package to project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("package added to project", map[string]interface{}{"id": opts.ID, "reference": opts.Reference})

	return nil
}

// RemoveProjectPackage removes a package from a project
func (d *Driver) RemoveProjectPackage(opts RemoveProjectPackageOpts) error {
	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.RemovePackage(d.ctx, &projectservice.RemoveProjectPackageOpts{
		ID:        opts.ID,
		Reference: opts.Reference,
	}); err != nil {
		d.logger.Error("failed to remove package from project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("package removed from project", map[string]interface{}{"id": opts.ID, "reference": opts.Reference})

	return nil
}

// UpdateProject updates a project
func (d *Driver) UpdateProject(opts UpdateProjectOpts) error {
	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Update(d.ctx, &projectservice.UpdateProjectOpts{
		ID:   opts.ID,
		Name: opts.Name,
	}); err != nil {
		d.logger.Error("failed to update project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("project updated", map[string]interface{}{"id": opts.ID})

	return nil
}

// DeleteProject deletes a project
func (d *Driver) DeleteProject(id uuid.UUID) error {
	d.logger.Debug("deleting project", map[string]interface{}{"id": id})
	return nil
}

// RunProject runs a project
func (d *Driver) RunProject(id uuid.UUID) error {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Run(ctx, id); err != nil {
		d.logger.Error("failed to run project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("project started", map[string]interface{}{"id": id})
	return nil
}

// RenderProject renders a project
func (d *Driver) RenderProject(id uuid.UUID) error {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Render(ctx, id); err != nil {
		d.logger.Error("failed to render project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("project rendered", map[string]interface{}{"id": id})
	return nil
}

// ExploreProject explores a project
func (d *Driver) ExploreProject(id uuid.UUID) error {
	ctx := context.Background()

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := projectService.Explore(ctx, id); err != nil {
		d.logger.Error("failed to explore project", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("project explored", map[string]interface{}{"id": id})
	return nil
}
