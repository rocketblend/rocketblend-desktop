package application

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	pack "github.com/rocketblend/rocketblend-desktop/internal/application/package"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/container"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/fileserver"
	"github.com/rocketblend/rocketblend-desktop/internal/buffer"
	"github.com/rocketblend/rocketblend-desktop/internal/eventwriter"
	"github.com/rocketblend/rocketblend/pkg/runtime"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

const id = "3088d60d-70ae-47e5-bf1f-cad698da3620"

type (
	Application interface {
		Execute() error
	}

	application struct {
		id       uuid.UUID
		platform runtime.Platform
		driver   *Driver
		handler  http.Handler
		assets   fs.FS
	}
)

func New(assets fs.FS) (Application, error) {
	events := buffer.New(buffer.WithMaxBufferSize(50))
	logger := logger.New(
		logger.WithLogLevel("debug"),
		logger.WithWriters(
			logger.PrettyWriter(),
			eventwriter.New(events),
		),
	)

	id, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	platform := runtime.DetectPlatform()
	if platform == runtime.Undefined {
		return nil, fmt.Errorf("unsupported platform")
	}

	container, err := container.New(
		container.WithLogger(logger),
		container.WithDevelopmentMode(true), // TODO: use build flag.
	)
	if err != nil {
		return nil, err
	}

	handler, err := fileserver.New(
		fileserver.WithContainer(container),
	)
	if err != nil {
		return nil, err
	}

	// TOOD: use options.
	driver, err := NewDriver(factory, events, platform)
	if err != nil {
		return nil, err
	}

	return &application{
		id:       id,
		platform: platform,
		assets:   assets,
		driver:   driver,
		handler:  handler,
	}, nil
}

func (a *application) Execute() error {
	frameless := false
	if a.platform == runtime.Windows {
		frameless = true
	}

	// Create application with options
	return wails.Run(&options.App{
		Title:  "RocketBlend Desktop",
		Width:  1400,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets:  a.assets,
			Handler: a.handler,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.DefaultAppearance,
			WebviewIsTransparent: true,
			About: &mac.AboutInfo{
				Title:   "RocketBlend Desktop",
				Message: "© 2024 RocketBlend",
			},
			OnFileOpen: func(filePath string) { fmt.Println(filePath) },
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               a.id.String(),
			OnSecondInstanceLaunch: a.driver.onSecondInstanceLaunch,
		},
		EnumBind: []interface{}{
			pack.AllPackageTypes,
			pack.AllPackageStates,
		},
		MinHeight:        580,
		MinWidth:         800,
		BackgroundColour: &options.RGBA{R: 00, G: 00, B: 00, A: 1},
		OnStartup:        a.driver.startup,
		OnShutdown:       a.driver.shutdown,
		OnDomReady:       a.driver.onDomReady,
		Frameless:        frameless,
		Bind: []interface{}{
			a.driver,
		},
	})
}
