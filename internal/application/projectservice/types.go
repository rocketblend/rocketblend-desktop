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

	UpdateProjectOpts struct {
		ID            uuid.UUID
		Name          *string
		Tags          *[]string
		ThumbnailPath *string
		SplashPath    *string
	}

	AddProjectPackageOpts struct {
		ID        uuid.UUID
		Reference reference.Reference
	}

	RemoveProjectPackageOpts struct {
		ID        uuid.UUID
		Reference reference.Reference
	}
)
