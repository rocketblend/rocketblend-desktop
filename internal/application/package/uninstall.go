package pack

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	"github.com/rocketblend/rocketblend/pkg/reference"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

func (r *Repository) UninstallPackage(ctx context.Context, opts *types.UninstallPackageOpts) error {
	if err := r.uninstall(ctx, opts.ID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) uninstall(ctx context.Context, id uuid.UUID) error {
	pack, err := r.get(ctx, id)
	if err != nil {
		return err
	}

	if err := r.rbRepository.RemoveInstallations(ctx, &rbtypes.RemoveInstallationsOpts{
		References: []reference.Reference{pack.Reference},
	}); err != nil {
		return err
	}

	return nil
}
