package pack

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	"github.com/rocketblend/rocketblend/pkg/reference"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

func (r *Repository) InstallPackage(ctx context.Context, opts *types.InstallPackageOpts) (err error) {
	if err := r.install(ctx, opts.ID, opts.BypassUntrustedDomainCheck); err != nil {
		return fmt.Errorf("error installing package: %w", err)
	}

	return nil
}

func (r *Repository) install(ctx context.Context, id uuid.UUID, bypassUntrustedDomainCheck bool) error {
	item, err := r.get(ctx, id)
	if err != nil {
		return err
	}

	// TODO: Check if the package is already installed.

	results, err := r.rbRepository.GetPackages(ctx, &rbtypes.GetPackagesOpts{
		References: []reference.Reference{item.Reference},
	})
	if err != nil {
		return err
	}

	if len(results.Packs) == 0 {
		return fmt.Errorf("package not found")
	}

	_, err = r.rbRepository.GetInstallations(ctx, &rbtypes.GetInstallationsOpts{
		Dependencies: []*rbtypes.Dependency{
			{
				Reference: item.Reference,
				Type:      results.Packs[item.Reference].Type,
			},
		},
		Fetch:                      true,
		BypassUntrustedDomainCheck: bypassUntrustedDomainCheck,
	})
	if err != nil {
		return err
	}

	return nil
}
