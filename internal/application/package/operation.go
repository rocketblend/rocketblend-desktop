package pack

import (
	"context"

	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
)

func (r *Repository) AddPackageOperation(ctx context.Context, opts *types.AddPackageOperationOpts) error {
	pack, err := r.get(ctx, opts.ID)
	if err != nil {
		return err
	}

	pack.Operations = append(pack.Operations, opts.OperationID.String())

	return r.insert(ctx, pack)
}

func (r *Repository) RemovePackageOperation(ctx context.Context, opts *types.RemovePackageOperationOpts) error {
	pack, err := r.get(ctx, opts.ID)
	if err != nil {
		return err
	}

	for i, operationID := range pack.Operations {
		if operationID == opts.OperationID.String() {
			pack.Operations = append(pack.Operations[:i], pack.Operations[i+1:]...)
			break
		}
	}

	return r.insert(ctx, pack)
}

func (r *Repository) insert(ctx context.Context, pack *types.Package) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	index, err := convertToIndex(pack)
	if err != nil {
		return err
	}

	return r.store.Insert(ctx, index)
}
