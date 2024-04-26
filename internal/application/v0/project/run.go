package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

func (r *repository) Run(ctx context.Context, opts *types.RunProjectOpts) error {
	if err := r.run(ctx, opts.ID); err != nil {
		return err
	}

	return nil
}

func (r *repository) run(ctx context.Context, id uuid.UUID) error {
	// project, err := r.get(ctx, id)
	// if err != nil {
	// 	return err
	// }

	// driver, err := r.createDriver(project.BlendFile())
	// if err != nil {
	// 	return err
	// }

	// go func() {
	// 	if err := driver.Run(ctx); err != nil {
	// 		r.logger.Error("failed to run project", map[string]interface{}{"error": err})
	// 	}
	// }()

	// r.emitEvent(ctx, id, RunEventChannel)

	return nil
}
