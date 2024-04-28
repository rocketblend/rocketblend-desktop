package pack

import (
	"context"

	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

// func (r *Repository) AppendOperation(ctx context.Context, id uuid.UUID, opid uuid.UUID) error {
// 	pack, err := r.get(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	pack.Operations = append(pack.Operations, opid.String())

// 	return r.insert(ctx, pack)
// }

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
