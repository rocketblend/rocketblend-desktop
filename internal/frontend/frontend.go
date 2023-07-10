package frontend

import (
	"io/fs"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

type Frontend struct {
	app    *App
	assets fs.FS
}

func New(assets fs.FS) *Frontend {
	return &Frontend{
		assets: assets,
		app:    NewApp(),
	}
}

func (f *Frontend) Execute() error {
	// Create application with options
	return wails.Run(&options.App{
		Title:  "RocketBlend Desktop",
		Width:  1024,
		Height: 768,
		Menu:   f.app.menu(),
		AssetServer: &assetserver.Options{
			Assets: f.assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        f.app.startup,
		OnShutdown:       f.app.shutdown,
		Bind: []interface{}{
			f.app,
		},
	})
}
