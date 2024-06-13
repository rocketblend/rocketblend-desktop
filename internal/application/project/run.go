package project

import (
	"context"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/events"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

func (r *Repository) RunProject(ctx context.Context, opts *types.RunProjectOpts) error {
	if err := r.run(ctx, opts.ID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) run(ctx context.Context, id uuid.UUID) error {
	project, err := r.get(ctx, id)
	if err != nil {
		return err
	}

	result, err := r.rbDriver.ResolveProfiles(ctx, &rbtypes.ResolveProfilesOpts{
		Profiles: []*rbtypes.Profile{
			project.Profile(),
		},
	})
	if err != nil {
		return err
	}

	go func() {
		if err := r.blender.Run(ctx, &rbtypes.RunOpts{
			BlenderOpts: rbtypes.BlenderOpts{
				BlendFile: &rbtypes.BlendFile{
					Path:         filepath.Join(project.Path, project.FileName),
					Dependencies: result.Installations[0],
					Strict:       project.Strict,
				},
			},
		}); err != nil {
			r.logger.Error("failed to run project", map[string]interface{}{"error": err})
		}
	}()

	r.emitEvent(ctx, id, events.ProjectRunChannel)

	return nil
}
