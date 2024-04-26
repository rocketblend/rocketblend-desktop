package project

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
	"github.com/rocketblend/rocketblend/pkg/reference"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

func (r *repository) AddPackage(ctx context.Context, opts *types.AddProjectPackageOpts) error {
	if err := r.addPackage(ctx, opts.ID, opts.Reference); err != nil {
		return err
	}

	return nil
}

func (r *repository) addPackage(ctx context.Context, id uuid.UUID, reference reference.Reference) error {
	project, err := r.get(ctx, id)
	if err != nil {
		return err
	}

	if project.HasDependency(reference) {
		return errors.New("package already exists on project")
	}

	profile := project.Profile()
	profile.AddDependencies(&rbtypes.Dependency{
		Reference: reference,
	})

	if err := r.rbDriver.TidyProfiles(ctx, &rbtypes.TidyProfilesOpts{
		Profiles: []*rbtypes.Profile{profile},
	}); err != nil {
		return err
	}

	if err := r.rbDriver.SaveProfiles(ctx, &rbtypes.SaveProfilesOpts{
		Profiles: map[string]*rbtypes.Profile{
			project.Path: profile,
		},
	}); err != nil {
		return err
	}

	return nil
}
