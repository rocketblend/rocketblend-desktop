package application

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/container"
	"github.com/rocketblend/rocketblend-desktop/internal/application/fileserver"
	"github.com/rocketblend/rocketblend-desktop/internal/buffer"
	"github.com/rocketblend/rocketblend-desktop/internal/eventwriter"
	"github.com/rocketblend/rocketblend/pkg/runtime"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

const (
	id    = "3088d60d-70ae-47e5-bf1f-cad698da3620"
	title = "RocketBlend Desktop"
)

type (
	ApplicationOpts struct {
		Assets  fs.FS
		Version string
		Args    []string
	}

	Application struct {
		id       uuid.UUID
		platform runtime.Platform
		driver   *Driver
		handler  http.Handler
		assets   fs.FS
	}
)

func New(opts ApplicationOpts) (*Application, error) {
	if opts.Assets == nil {
		return nil, fmt.Errorf("assets are required")
	}

	logLevel := "info"
	development := false
	if opts.Version == "dev" {
		development = true
		logLevel = "debug"
	}

	events := buffer.New(buffer.WithMaxBufferSize(50))
	logger := logger.New(
		logger.WithLogLevel(logLevel),
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
		container.WithDevelopmentMode(development),
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

	driver, err := NewDriver(
		WithContainer(container),
		WithWriter(events),
		WithPlatform(platform),
		WithVersion(opts.Version),
		WithArgs(opts.Args...),
	)
	if err != nil {
		return nil, err
	}

	return &Application{
		id:       id,
		platform: platform,
		assets:   opts.Assets,
		driver:   driver,
		handler:  handler,
	}, nil
}

func (a *Application) Execute() error {
	frameless := false
	if a.platform == runtime.Windows {
		frameless = true
	}

	// Create application with options
	return wails.Run(&options.App{
		Title:  title,
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
				Title:   title,
				Message: "© 2024 RocketBlend",
			},
			OnFileOpen: func(filePath string) { fmt.Println(filePath) },
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               a.id.String(),
			OnSecondInstanceLaunch: a.driver.onSecondInstanceLaunch,
		},
		EnumBind: []interface{}{
			//pack.AllPackageTypes,
			//pack.AllPackageStates,
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
