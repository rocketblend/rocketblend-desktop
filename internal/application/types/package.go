package types

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/enums"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/listoption"
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

	Progress struct {
		CurrentBytes   int64   `json:"currentBytes"`
		TotalBytes     int64   `json:"totalBytes"`
		BytesPerSecond float64 `json:"bytesPerSecond"`
	}

	Package struct {
		ID               uuid.UUID           `json:"id"`
		Type             enums.PackageType   `json:"type"`
		State            enums.PackageState  `json:"state"`
		Reference        reference.Reference `json:"reference"`
		Name             string              `json:"name"`
		Author           string              `json:"author"`
		Tag              string              `json:"tag"`
		Path             string              `json:"path"`
		Verified         bool                `json:"verified"`
		InstallationPath string              `json:"installationPath"`
		Operations       []string            `json:"operations"`
		Platform         types.Platform      `json:"platform"`
		URI              *types.URI          `json:"uri"`
		Version          *semver.Version     `json:"version"`
		Progress         *Progress           `json:"progress,omitempty"`
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
		Update    bool                `json:"update"`
	}

	InstallPackageOpts struct {
		ID uuid.UUID `json:"id"`
	}

	UninstallPackageOpts struct {
		ID uuid.UUID `json:"id"`
	}

	AddPackageOperationOpts struct {
		ID          uuid.UUID `json:"id"`
		OperationID uuid.UUID `json:"operationID"`
	}

	RemovePackageOperationOpts struct {
		ID          uuid.UUID `json:"id"`
		OperationID uuid.UUID `json:"operationID"`
	}

	Catalog interface {
		GetPackage(ctx context.Context, opts *GetPackageOpts) (*GetPackageResponse, error)
		ListPackages(ctx context.Context, opts ...listoption.ListOption) (*ListPackagesResponse, error) // TODO: Change opts struct.

		AddPackageOperation(ctx context.Context, opts *AddPackageOperationOpts) error
		RemovePackageOperation(ctx context.Context, opts *RemovePackageOperationOpts) error

		AddPackage(ctx context.Context, opts *AddPackageOpts) error

		InstallPackage(ctx context.Context, opts *InstallPackageOpts) error
		UninstallPackage(ctx context.Context, opts *UninstallPackageOpts) error

		RefreshPackages(ctx context.Context) error

		Close() error
	}
)
