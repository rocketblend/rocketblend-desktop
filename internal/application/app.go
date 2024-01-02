package application

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/buffermanager"
	"github.com/rocketblend/rocketblend-desktop/internal/application/factory"
	pack "github.com/rocketblend/rocketblend-desktop/internal/application/package"
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

func New(assets fs.FS) (Application, error) {
	events := buffermanager.New(buffermanager.WithMaxBufferSize(50))
	logger := logger.New(
		logger.WithLogLevel("debug"),
		//logger.WithPretty(),
		logger.WithWriters(NewEventBufferWriter(events)),
	)

	id, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	factory, err := factory.New(
		factory.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	driver, err := NewDriver(factory, events)
	if err != nil {
		return nil, err
	}

	cacheTimeout := 3600 // 1 hour
	handler, err := NewFileLoader(logger, factory, cacheTimeout)
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
		Frameless:        true,
		Bind: []interface{}{
			a.driver,
		},
	})
}
