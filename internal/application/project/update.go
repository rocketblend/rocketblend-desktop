package project

import (
	"context"

	"github.com/rocketblend/rocketblend-desktop/internal/application/events"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
)

func (r *Repository) UpdateProject(ctx context.Context, opts *types.UpdateProjectOpts) error {
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

	if err := r.saveDetail(project.Path, detail, false, false); err != nil {
		return err
	}

	r.emitEvent(ctx, project.ID, events.ProjectUpdateChannel)

	return nil
}
