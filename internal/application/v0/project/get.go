package project

import (
	"context"
	"encoding/json"
	"path"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

func (r *repository) get(ctx context.Context, id uuid.UUID) (*types.Project, error) {
	index, err := r.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	project, err := convertFromIndex(index)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func convertFromIndex(index *types.Index) (*types.Project, error) {
	var result types.Project
	if err := json.Unmarshal([]byte(index.Data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func convertToIndex(project *types.Project) (*types.Index, error) {
	data, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	resources := []string{}
	if project.ThumbnailPath != "" {
		resources = append(resources, filepath.ToSlash(project.ThumbnailPath))
		resources = append(resources, filepath.ToSlash(project.SplashPath))
	}

	return &types.Index{
		ID:        project.ID,
		Name:      project.Name,
		Type:      indextype.Project,
		Reference: path.Clean(project.Path),
		Resources: resources,
		Data:      string(data),
	}, nil
}
