package project

import (
	"context"

	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

func (r *repository) UpdateProject(ctx context.Context, opts *types.UpdateProjectOpts) error {
	project, err := r.get(ctx, opts.ID)
	if err != nil {
		return err
	}

	detail := project.Detail()
	if opts.Name != nil {
		detail.Name = *opts.Name
	}

	if opts.Tags != nil {
		detail.Tags = *opts.Tags
	}

	if opts.ThumbnailPath != nil {
		detail.ThumbnailPath = *opts.ThumbnailPath
	}

	if opts.SplashPath != nil {
		detail.SplashPath = *opts.SplashPath
	}

	if err := r.saveDetail(project.Path, detail); err != nil {
		return err
	}

	r.emitEvent(ctx, project.ID, UpdateEventChannel)

	return nil
}
