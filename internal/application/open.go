package application

import (
	"context"
	"path/filepath"

	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

func openWithRocketBlend(ctx context.Context, driver rbtypes.Driver, blender rbtypes.Blender, blendFilePath string) error {
	profiles, err := driver.LoadProfiles(ctx, &rbtypes.LoadProfilesOpts{
		Paths: []string{filepath.Dir(blendFilePath)},
	})
	if err != nil {
		return err
	}

	resolve, err := driver.ResolveProfiles(ctx, &rbtypes.ResolveProfilesOpts{
		Profiles: profiles.Profiles,
	})
	if err != nil {
		return err
	}

	if err := blender.Run(ctx, &rbtypes.RunOpts{
		BlenderOpts: rbtypes.BlenderOpts{
			BlendFile: &rbtypes.BlendFile{
				Path:         blendFilePath,
				Dependencies: resolve.Installations[0],
			},
		},
	}); err != nil {
		return err
	}

	return nil
}
