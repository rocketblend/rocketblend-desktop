package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application"
	"github.com/rocketblend/rocketblend/pkg/driver"
	"github.com/rocketblend/rocketblend/pkg/driver/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketfile"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/factory"
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
	rbFactory, err := factory.New()
	if err != nil {
		return err
	}

	logger := logger.New(
		logger.WithLogLevel("debug"),
	)

	if err := rbFactory.SetLogger(logger); err != nil {
		return err
	}

	var messages []string
	if len(os.Args) > 1 {
		ctx := context.Background()
		if err := open(ctx, os.Args[1], logger, rbFactory); err == nil {
			// If we successfully launched a project, we're done.
			return nil
		}

		// If we failed to launch a project directly, open with application.
		messages = append(messages, err.Error())
	}

	app, err := application.New(logger, rbFactory, assets, messages...)
	if err != nil {
		return err
	}

	if err := app.Execute(); err != nil {
		return err
	}

	return nil
}

func open(ctx context.Context, blendFilePath string, logger logger.Logger, factory factory.Factory) error {
	rocketPackService, err := factory.GetRocketPackService()
	if err != nil {
		return err
	}

	installationService, err := factory.GetInstallationService()
	if err != nil {
		return err
	}

	blendFileService, err := factory.GetBlendFileService()
	if err != nil {
		return err
	}

	blendFile, err := blendconfig.Load(blendFilePath, filepath.Join(filepath.Dir(blendFilePath), rocketfile.FileName))
	if err != nil {
		return err
	}

	driver, err := driver.New(
		driver.WithLogger(logger),
		driver.WithRocketPackService(rocketPackService),
		driver.WithInstallationService(installationService),
		driver.WithBlendFileService(blendFileService),
		driver.WithBlendConfig(blendFile),
	)

	if err != nil {
		return err
	}

	return driver.Run(ctx)
}
