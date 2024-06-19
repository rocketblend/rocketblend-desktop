package types

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/enums"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/listoption"
	"github.com/rocketblend/rocketblend/pkg/reference"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

const IgnoreFileName = ".rocketignore"

type (
	Media struct {
		FilePath  string `json:"filePath"`
		URL       string `json:"url"`
		Splash    bool   `json:"splash"`
		Thumbnail bool   `json:"thumbnail"`
		// Height    int    `json:"height"`
		// Width     int    `json:"width"`
	}

	Dependency struct {
		Reference reference.Reference `json:"reference"`
		Type      enums.PackageType   `json:"type"`
	}

	Project struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
		Tags []string  `json:"tags"`
		Path string    `json:"path"`

		MediaPath string `json:"mediaPath"`
		FileName  string `json:"fileName"`

		Dependencies []*Dependency `json:"dependencies"`
		Media        []*Media      `json:"media"`

		Strict bool `json:"strict"`

		Version   string    `json:"version"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	GetProjectOpts struct {
		ID uuid.UUID `json:"id"`
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
		ID   uuid.UUID `json:"id"`
		Name *string   `json:"name"`
		Tags *[]string `json:"tags"`
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

	Portfolio interface {
		GetProject(ctx context.Context, opts *GetProjectOpts) (*GetProjectResponse, error)
		ListProjects(ctx context.Context, opts ...listoption.ListOption) (*ListProjectsResponse, error)

		CreateProject(ctx context.Context, opts *CreateProjectOpts) (*CreateProjectResult, error)
		UpdateProject(ctx context.Context, opts *UpdateProjectOpts) error
		AddProjectPackage(ctx context.Context, opts *AddProjectPackageOpts) error
		RemoveProjectPackage(ctx context.Context, opts *RemoveProjectPackageOpts) error

		//RenderProject(ctx context.Context, id uuid.UUID) error
		RunProject(ctx context.Context, opts *RunProjectOpts) error

		Refresh(ctx context.Context) error

		Close() error
	}
)

func (p *Project) Profile() *rbtypes.Profile {
	dependency := make([]*rbtypes.Dependency, 0, len(p.Dependencies))
	for _, d := range p.Dependencies {
		dependency = append(dependency, &rbtypes.Dependency{
			Reference: d.Reference,
			Type:      rbtypes.PackageType(d.Type),
		})
	}

	return &rbtypes.Profile{
		Dependencies: dependency,
		Strict:       p.Strict,
	}
}

func (p *Project) Detail() *Detail {
	return &Detail{
		ID:        p.ID,
		Name:      p.Name,
		Tags:      p.Tags,
		MediaPath: p.MediaPath,
	}
}

func (p *Project) HasDependency(dep reference.Reference) bool {
	for _, d := range p.Dependencies {
		if d.Reference == dep {
			return true
		}
	}

	return false
}
