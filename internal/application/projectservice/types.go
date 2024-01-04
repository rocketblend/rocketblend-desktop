package projectservice

import (
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
)

type (
	GetProjectResponse struct {
		Project *project.Project `json:"project,omitempty"`
	}

	ListProjectsResponse struct {
		Projects []*project.Project `json:"projects,omitempty"`
	}

	CreateProjectRequest struct {
		Name     string                `json:"name,omitempty"`
		Tags     []string              `json:"tags,omitempty"`
		Path     string                `json:"path,omitempty"`
		FileName string                `json:"fileName,omitempty"`
		Build    reference.Reference   `json:"build,omitempty"`
		Addons   []reference.Reference `json:"addons,omitempty"`
	}

	UpdateProjectRequest struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
)
