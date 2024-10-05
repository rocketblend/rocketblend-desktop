package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/rocketblend/rocketblend-desktop/internal/application"
)

// 'wails dev' should properly launch vite to serve the site
// for live development without needing to seperately launch
// 'npm run dev' or your flavor such as pnpm in the frontend
// directory seperately.

// The comment below chooses what gets packaged with
// the application.

//go:embed all:frontend/build
var assets embed.FS

var (
	Version               = "dev"
	BuildType      string = "debug"
	BuildTimestamp string = "NOW"
	CommitSha      string = "HEAD"
	BuildLink      string = "http://localhost"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Println("Something went wrong:", err)
	}
}

func run(args []string) error {
	app, err := application.New(application.ApplicationOpts{
		Assets:  assets,
		Version: Version,
		Args:    args,
	})
	if err != nil {
		return err
	}

	if len(os.Args) > 1 {
		if err := app.Open(context.Background(), os.Args[1]); err == nil {
			// If we successfully launched a project, we're done.
			return nil
		}

		// If we failed to launch a project directly, open with application.
	}

	if err := app.Execute(); err != nil {
		return err
	}

	return nil
}
