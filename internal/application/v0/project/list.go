package project

import (
	"context"

	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

func (r *repository) ListProjects(ctx context.Context, opts ...listoption.ListOption) (*types.ListProjectsResponse, error) {
	opts = append(opts, listoption.WithType(indextype.Project))
	indexes, err := r.store.List(ctx, opts...)
	if err != nil {
		return nil, err
	}

	projects := make([]*types.Project, 0, len(indexes))
	for _, index := range indexes {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		pack, err := convertFromIndex(index)
		if err != nil {
			return nil, err
		}

		projects = append(projects, pack)
	}

	r.logger.Debug("Found projects", map[string]interface{}{
		"projects": len(projects),
		"indexes":  len(indexes),
	})

	return &types.ListProjectsResponse{
		Projects: projects,
	}, nil
}
