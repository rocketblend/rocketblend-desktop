package main

import (
	"embed"

	"github.com/rocketblend/rocketblend-desktop/internal/frontend"
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
	frontend := frontend.New(assets)

	err := frontend.Execute()
	if err != nil {
		println("Error:", err.Error())
	}
}
