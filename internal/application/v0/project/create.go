package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

func (r *repository) Create(ctx context.Context, opts *types.CreateProjectOpts) (*types.CreateProjectResult, error) {
	id, err := r.create(ctx, opts)
	if err != nil {
		return nil, err
	}

	r.emitEvent(ctx, id, CreateEventChannel)

	return &types.CreateProjectResult{
		ID: id,
	}, nil
}

func (r *repository) create(ctx context.Context, opts *types.CreateProjectOpts) (uuid.UUID, error) {
	// blendConfig, err := blendconfig.New(
	// 	opts.Path,
	// 	ensureBlendExtension(opts.BlendFileName),
	// 	rocketfile.New(opts.Build),
	// )
	// if err != nil {
	// 	return uuid.Nil, err
	// }

	// driver, err := s.createDriver(blendConfig)
	// if err != nil {
	// 	return uuid.Nil, err
	// }

	// // We create a temporary ignore file to avoid adding the project to the index before it is fully created.
	// if err := createIgnoreFile(opts.Path); err != nil {
	// 	return uuid.Nil, err
	// }

	// defer func() {
	// 	err := removeIgnoreFile(opts.Path)
	// 	if err != nil {
	// 		s.logger.Error("failed to remove temporarily project ignore file", map[string]interface{}{"error": err})
	// 	}
	// }()

	// if err := driver.Create(ctx); err != nil {
	// 	return uuid.Nil, err
	// }

	// id := uuid.New()
	// if err := projectsettings.Save(&projectsettings.ProjectSettings{
	// 	ID:   id,
	// 	Name: opts.DisplayName,
	// }, filepath.Join(opts.Path, projectsettings.FileName)); err != nil {
	// 	return uuid.Nil, err
	// }

	// return id, nil

	return uuid.New(), nil
}
