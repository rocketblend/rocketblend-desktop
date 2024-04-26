package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

func (r *repository) RemovePackage(ctx context.Context, opts *types.RemoveProjectPackageOpts) error {
	if err := r.removePackage(ctx, opts.ID); err != nil {
		return err
	}

	return nil
}

func (r *repository) removePackage(ctx context.Context, id uuid.UUID) error {
	// project, err := r.get(ctx, id)
	// if err != nil {
	// 	return err
	// }

	// if !project.HasDependency(opts.Reference) {
	// 	return errors.New("package does not exist on project")
	// }

	// driver, err := s.createDriver(project.BlendFile())
	// if err != nil {
	// 	return err
	// }

	// if err := driver.RemoveDependencies(ctx, opts.Reference); err != nil {
	// 	return err
	// }

	return nil
}
