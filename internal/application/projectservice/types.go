package projectservice

import (
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
)

type (
	GetProjectResponse struct {
		Project *project.Project `json:"project,omitempty"`
	}

	ListProjectsResponse struct {
		Projects []*project.Project `json:"projects,omitempty"`
	}

	UpdateProjectRequest struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
)
