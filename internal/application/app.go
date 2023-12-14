package application

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/packageservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/factory"
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
		id      uuid.UUID
		driver  *Driver
		handler http.Handler
		assets  fs.FS
	}
)

func New(logger logger.Logger, factory factory.Factory, assets fs.FS) (Application, error) {
	id, err := uuid.Parse(id)
	if err != nil {
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
		id:      id,
		assets:  assets,
		driver:  driver,
		handler: handler,
	}, nil
}

func (a *application) Execute() error {
	// Create application with options
	return wails.Run(&options.App{
		Title:  "RocketBlend Desktop",
		Width:  1150,
		Height: 780,
		AssetServer: &assetserver.Options{
			Assets:  a.assets,
			Handler: a.handler,
		},
		Mac: &mac.Options{
			OnFileOpen: func(filePath string) { fmt.Println(filePath) },
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               a.id.String(),
			OnSecondInstanceLaunch: a.driver.onSecondInstanceLaunch,
		},
		MinHeight:        580,
		MinWidth:         800,
		BackgroundColour: &options.RGBA{R: 00, G: 00, B: 00, A: 1},
		OnStartup:        a.driver.startup,
		OnShutdown:       a.driver.shutdown,
		OnDomReady:       a.driver.onDomReady,
		Frameless:        true,
		Bind: []interface{}{
			a.driver,
		},
	})
}
