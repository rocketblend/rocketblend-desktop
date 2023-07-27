package main

import (
	"embed"
	"fmt"

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

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func run() error {
	app, err := application.New(assets)
	if err != nil {
		return err
	}

	if err := app.Execute(); err != nil {
		return err
	}

	return nil
}
