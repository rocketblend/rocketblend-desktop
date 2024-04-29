package pack

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
)

func (r *Repository) GetPackage(ctx context.Context, opts *types.GetPackageOpts) (*types.GetPackageResponse, error) {
	pack, err := r.get(ctx, opts.ID)
	if err != nil {
		return nil, err
	}

	return &types.GetPackageResponse{
		Package: pack,
	}, nil
}

func (r *Repository) get(ctx context.Context, id uuid.UUID) (*types.Package, error) {
	index, err := r.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	pack, err := convertFromIndex(index)
	if err != nil {
		return nil, err
	}

	return pack, nil
}
