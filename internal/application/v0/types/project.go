package types

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
	"github.com/rocketblend/rocketblend/pkg/reference"
)

type (
	Project struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		Tags     []string  `json:"tags"`
		Path     string    `json:"path"`
		FileName string    `json:"fileName"`

		Build  reference.Reference   `json:"build"`
		Addons []reference.Reference `json:"addons"`

		SplashPath    string    `json:"splashPath"`
		ThumbnailPath string    `json:"thumbnailPath"`
		Version       string    `json:"version"`
		UpdatedAt     time.Time `json:"updatedAt"`
	}

	GetProjectResponse struct {
		Project *Project `json:"project,omitempty"`
	}

	ListProjectsResponse struct {
		Projects []*Project `json:"projects,omitempty"`
	}

	CreateProjectOpts struct {
		DisplayName   string              `json:"displayName"`
		BlendFileName string              `json:"blendFileName"`
		Path          string              `json:"path"`
		Build         reference.Reference `json:"build"`
	}

	CreateProjectResult struct {
		ID uuid.UUID
	}

	UpdateProjectOpts struct {
		ID            uuid.UUID `json:"id"`
		Name          *string   `json:"name"`
		Tags          *[]string `json:"tags"`
		ThumbnailPath *string   `json:"thumbnailPath"`
		SplashPath    *string   `json:"splashPath"`
	}

	AddProjectPackageOpts struct {
		ID        uuid.UUID           `json:"id"`
		Reference reference.Reference `json:"reference"`
	}

	RemoveProjectPackageOpts struct {
		ID        uuid.UUID           `json:"id"`
		Reference reference.Reference `json:"reference"`
	}

	RunProjectOpts struct {
		ID uuid.UUID `json:"id"`
	}

	ProjectRepository interface {
		Get(ctx context.Context, id uuid.UUID) (*GetProjectResponse, error)
		List(ctx context.Context, opts ...listoption.ListOption) (*ListProjectsResponse, error)

		Create(ctx context.Context, opts *CreateProjectOpts) (*CreateProjectResult, error)
		Update(ctx context.Context, opts *UpdateProjectOpts) error
		AddPackage(ctx context.Context, opts *AddProjectPackageOpts) error
		RemovePackage(ctx context.Context, opts *RemoveProjectPackageOpts) error

		//Render(ctx context.Context, id uuid.UUID) error
		Run(ctx context.Context, opts *RunProjectOpts) error

		Refresh(ctx context.Context) error

		Close() error
	}
)
