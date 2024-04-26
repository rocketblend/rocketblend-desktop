package project

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
	"github.com/rocketblend/rocketblend/pkg/reference"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

func (r *repository) RemovePackage(ctx context.Context, opts *types.RemoveProjectPackageOpts) error {
	if err := r.removePackage(ctx, opts.ID, opts.Reference); err != nil {
		return err
	}

	return nil
}

func (r *repository) removePackage(ctx context.Context, id uuid.UUID, reference reference.Reference) error {
	project, err := r.get(ctx, id)
	if err != nil {
		return err
	}

	if !project.HasDependency(reference) {
		return errors.New("package does not exist on project")
	}

	profile := project.Profile()
	profile.RemoveDependencies(&rbtypes.Dependency{
		Reference: reference,
	})

	if err := r.rbDriver.SaveProfiles(ctx, &rbtypes.SaveProfilesOpts{
		Profiles: map[string]*rbtypes.Profile{
			project.Path: profile,
		},
	}); err != nil {
		return err
	}

	return nil
}
