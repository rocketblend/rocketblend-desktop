package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application"
	"github.com/rocketblend/rocketblend-desktop/internal/application/build"
	"github.com/rocketblend/rocketblend/pkg/container"
	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/types"
)

// 'wails dev' should properly launch vite to serve the site
// for live development without needing to seperately launch
// 'npm run dev' or your flavor such as pnpm in the frontend
// directory seperately.

// The comment below chooses what gets packaged with
// the application.

//go:embed all:frontend/build
var assets embed.FS

func main() {
	if err := run(os.Args); err != nil {
		fmt.Println("Error:", err)
	}
}

func run(args []string) error {
	logger := logger.NoOp()

	if len(os.Args) > 1 {
		var err error
		if err = open(context.Background(), os.Args[1], logger); err == nil {
			// If we successfully launched a project, we're done.
			return nil
		}

		// If we failed to launch a project directly, open with application.
	}

	app, err := application.New(assets)
	if err != nil {
		return err
	}

	if err := app.Execute(); err != nil {
		return err
	}

	return nil
}

func open(ctx context.Context, blendFilePath string, logger logger.Logger) error {
	development := false
	if build.Version == "dev" {
		development = true
	}

	container, err := container.New(
		container.WithLogger(logger),
		container.WithApplicationName(types.ApplicationName),
		container.WithDevelopmentMode(development),
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	driver, err := container.GetDriver()
	if err != nil {
		return err
	}

	profiles, err := driver.LoadProfiles(ctx, &types.LoadProfilesOpts{
		Paths: []string{filepath.Dir(blendFilePath)},
	})
	if err != nil {
		return err
	}

	resolve, err := driver.ResolveProfiles(ctx, &types.ResolveProfilesOpts{
		Profiles: profiles.Profiles,
	})
	if err != nil {
		return err
	}

	blender, err := container.GetBlender()
	if err != nil {
		return err
	}

	if err := blender.Run(ctx, &types.RunOpts{
		BlenderOpts: types.BlenderOpts{
			BlendFile: &types.BlendFile{
				Name:         helpers.ExtractName(blendFilePath),
				Path:         blendFilePath,
				Dependencies: resolve.Installations[0],
			},
		},
	}); err != nil {
		return err
	}

	return nil
}
