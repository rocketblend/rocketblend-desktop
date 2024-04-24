package types

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/semver"
	"github.com/rocketblend/rocketblend/pkg/types"
)

const PackageFileName = types.PackageFileName

type (
	Definition struct {
		types.Package
	}

	// TODO: Just embed Definition in Package.

	Package struct {
		ID   uuid.UUID         `json:"id"`
		Type types.PackageType `json:"type"`
		//Type             PackageType         `json:"type"`
		//State            PackageState        `json:"state"`
		Reference        reference.Reference `json:"reference"`
		Name             string              `json:"name"`
		Author           string              `json:"author"`
		Tag              string              `json:"tag"`
		Path             string              `json:"path"`
		InstallationPath string              `json:"installationPath"`
		Operations       []string            `json:"operations"`
		Platform         types.Platform      `json:"platform"`
		URI              *types.URI          `json:"uri"`
		Version          *semver.Version     `json:"version"`
		Verified         bool                `json:"verified"`
		UpdatedAt        time.Time           `json:"updatedAt"`
	}

	GetPackageOpts struct {
		ID uuid.UUID `json:"id"`
	}

	GetPackageResponse struct {
		Package *Package `json:"package,omitempty"`
	}

	ListPackagesResponse struct {
		Packages []*Package `json:"packages,omitempty"`
	}

	AddPackageOpts struct {
		Reference reference.Reference `json:"reference"`
	}

	InstallPackageOpts struct {
		ID uuid.UUID `json:"id"`
	}

	UninstallPackageOpts struct {
		ID uuid.UUID `json:"id"`
	}

	PackageRepoistory interface {
		GetPackage(ctx context.Context, opts *GetPackageOpts) (*GetPackageResponse, error)
		ListPackages(ctx context.Context, opts ...listoption.ListOption) (*ListPackagesResponse, error) // TODO: Change opts struct.

		//AppendOperation(ctx context.Context, id uuid.UUID, opid uuid.UUID) error

		AddPackage(ctx context.Context, opts *AddPackageOpts) error
		InstallPackage(ctx context.Context, opts *InstallPackageOpts) error
		UninstallPackage(ctx context.Context, opts *UninstallPackageOpts) error

		RefreshPackages(ctx context.Context) error

		Close() error
	}
)
