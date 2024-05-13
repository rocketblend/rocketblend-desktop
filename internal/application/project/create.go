package project

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/events"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	"github.com/rocketblend/rocketblend/pkg/reference"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

func (r *Repository) CreateProject(ctx context.Context, opts *types.CreateProjectOpts) (*types.CreateProjectResult, error) {
	profile, err := r.newProfile(ctx, opts.Build)
	if err != nil {
		return nil, err
	}

	// We create a temporary ignore file to avoid adding the project to the index before it is fully created.
	if err := createIgnoreFile(opts.Path); err != nil {
		return nil, err
	}

	defer func() {
		err := removeIgnoreFile(opts.Path)
		if err != nil {
			r.logger.Error("failed to remove temporarily project ignore file", map[string]interface{}{"error": err})
		}
	}()

	if err := r.createBlendFile(ctx, opts.DisplayName, filepath.Join(opts.Path, opts.BlendFileName), profile); err != nil {
		return nil, err
	}

	id := uuid.New()
	if err := r.saveDetail(opts.Path, &types.Detail{
		ID:        id,
		Name:      opts.DisplayName,
		MediaPath: DefaultMediaPath,
	}); err != nil {
		return nil, err
	}

	r.emitEvent(ctx, id, events.ProjectCreateChannel)

	return &types.CreateProjectResult{
		ID: id,
	}, nil
}

func (r *Repository) newProfile(ctx context.Context, build reference.Reference) (*rbtypes.Profile, error) {
	profiles := []*rbtypes.Profile{
		{
			Dependencies: []*rbtypes.Dependency{
				{
					Reference: build,
					Type:      rbtypes.PackageBuild,
				},
			},
		},
	}

	if err := r.rbDriver.TidyProfiles(ctx, &rbtypes.TidyProfilesOpts{
		Profiles: profiles,
	}); err != nil {
		return nil, err
	}

	return profiles[0], nil
}

func (r *Repository) createBlendFile(ctx context.Context, displayName string, filePath string, profile *rbtypes.Profile) error {
	if !strings.HasSuffix(filePath, ".blend") {
		return errors.New("filename must have .blend extension")
	}

	resolved, err := r.rbDriver.ResolveProfiles(ctx, &rbtypes.ResolveProfilesOpts{
		Profiles: []*rbtypes.Profile{profile},
	})
	if err != nil {
		return err
	}

	r.logger.Debug("creating blend file", map[string]interface{}{
		"displayName":  displayName,
		"filePath":     filePath,
		"dependencies": resolved.Installations[0],
	})

	if err := r.blender.Create(ctx, &rbtypes.CreateOpts{
		BlenderOpts: rbtypes.BlenderOpts{
			BlendFile: &rbtypes.BlendFile{
				Path:         filePath,
				Dependencies: resolved.Installations[0],
			},
			Background: true,
		},
	}); err != nil {
		return err
	}

	return nil
}

func createIgnoreFile(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(path, types.IgnoreFileName))
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func removeIgnoreFile(path string) error {
	return os.Remove(filepath.Join(path, types.IgnoreFileName))
}
