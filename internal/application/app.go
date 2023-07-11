package application

import (
	"io/fs"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

type (
	Application interface {
		Execute() error
	}

	application struct {
		driver *Driver
		assets fs.FS
	}
)

func New(assets fs.FS) Application {
	return &application{
		assets: assets,
		driver: NewDriver(),
	}
}

func (a *application) Execute() error {
	// Create application with options
	return wails.Run(&options.App{
		Title:  "RocketBlend Desktop",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: a.assets,
		},
		MinHeight:        580,
		MinWidth:         764,
		BackgroundColour: &options.RGBA{R: 00, G: 00, B: 00, A: 1},
		OnStartup:        a.driver.startup,
		OnShutdown:       a.driver.shutdown,
		Frameless:        true,
		Bind: []interface{}{
			a.driver,
		},
	})
}
