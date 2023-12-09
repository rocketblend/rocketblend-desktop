package application

import (
	"io/fs"
	"net/http"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/packageservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/factory"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

type (
	Application interface {
		Execute() error
	}

	application struct {
		driver  *Driver
		handler http.Handler
		assets  fs.FS
	}
)

func New(assets fs.FS) (Application, error) {
	logger := logger.New(
		logger.WithLogLevel("debug"),
	)

	factory, err := factory.New()
	if err != nil {
		return nil, err
	}

	if err := factory.SetLogger(logger); err != nil {
		return nil, err
	}

	storeService, err := searchstore.New(
		searchstore.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	configService, err := factory.GetConfigService()
	if err != nil {
		return nil, err
	}

	projectService, err := projectservice.New(
		projectservice.WithLogger(logger),
		projectservice.WithFactory(factory),
		projectservice.WithStore(storeService),
	)
	if err != nil {
		return nil, err
	}

	packageService, err := packageservice.New(
		packageservice.WithLogger(logger),
		packageservice.WithStore(storeService),
		packageservice.WithConfig(configService),
	)
	if err != nil {
		return nil, err
	}

	driver, err := NewDriver(logger, configService, projectService, packageService)
	if err != nil {
		return nil, err
	}

	cacheTimeout := 3600 // 1 hour
	handler, err := NewFileLoader(logger, storeService, cacheTimeout)
	if err != nil {
		return nil, err
	}

	return &application{
		assets:  assets,
		driver:  driver,
		handler: handler,
	}, nil
}

func (a *application) Execute() error {
	// Create application with options
	return wails.Run(&options.App{
		Title:  "RocketBlend Desktop",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  a.assets,
			Handler: a.handler,
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
