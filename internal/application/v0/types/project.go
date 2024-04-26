package types

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
	"github.com/rocketblend/rocketblend/pkg/reference"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

const IgnoreFileName = ".rocketignore"

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

func (p *Project) Profile() *rbtypes.Profile {
	dependencies := make([]*rbtypes.Dependency, 0, len(p.Addons)+1)
	dependencies = append(dependencies, &rbtypes.Dependency{
		Reference: p.Build,
		Type:      rbtypes.PackageBuild,
	})

	for _, addon := range p.Addons {
		dependencies = append(dependencies, &rbtypes.Dependency{
			Reference: addon,
			Type:      rbtypes.PackageAddon,
		})
	}

	return &rbtypes.Profile{
		Dependencies: dependencies,
	}
}

func (p *Project) Detail() *Detail {
	return &Detail{
		ID:            p.ID,
		Name:          p.Name,
		Tags:          p.Tags,
		ThumbnailPath: p.ThumbnailPath,
		SplashPath:    p.SplashPath,
	}
}

func (p *Project) HasDependency(dep reference.Reference) bool {
	if p.Build == dep {
		return true
	}

	for _, addon := range p.Addons {
		if addon == dep {
			return true
		}
	}

	return false
}
