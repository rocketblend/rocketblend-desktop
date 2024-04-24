package pack

import (
	"context"

	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

func (r *repository) ListPackages(ctx context.Context, opts ...listoption.ListOption) (*types.ListPackagesResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	opts = append(opts, listoption.WithType(indextype.Package))
	indexes, err := r.store.List(ctx, opts...)
	if err != nil {
		return nil, err
	}

	packs := make([]*types.Package, 0, len(indexes))
	for _, index := range indexes {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		pack, err := convertFromIndex(index)
		if err != nil {
			return nil, err
		}
		packs = append(packs, pack)
	}

	return &types.ListPackagesResponse{
		Packages: packs,
	}, nil
}
