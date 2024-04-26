package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

func (r *repository) AddPackage(ctx context.Context, opts *types.AddProjectPackageOpts) error {
	if err := r.addPackage(ctx, opts.ID); err != nil {
		return err
	}

	return nil
}

func (r *repository) addPackage(ctx context.Context, id uuid.UUID) error {
	// project, err := r.get(ctx, opts.ID)
	// if err != nil {
	// 	return err
	// }

	// // TODO: Check should be done on the driver function.
	// if project.HasDependency(opts.Reference) {
	// 	return errors.New("package already exists on project")
	// }

	// driver, err := r.createDriver(project.BlendFile())
	// if err != nil {
	// 	return err
	// }

	// if err := driver.AddDependencies(ctx, false, opts.Reference); err != nil {
	// 	return err
	// }

	return nil
}
