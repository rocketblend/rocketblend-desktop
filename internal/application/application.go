package application

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/container"
	"github.com/rocketblend/rocketblend-desktop/internal/application/enums"
	"github.com/rocketblend/rocketblend-desktop/internal/application/fileserver"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	"github.com/rocketblend/rocketblend-desktop/internal/buffer"
	"github.com/rocketblend/rocketblend-desktop/internal/eventwriter"
	"github.com/rocketblend/rocketblend/pkg/runtime"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
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
		id          uuid.UUID
		platform    runtime.Platform
		rocketblend rbtypes.Driver
		blender     rbtypes.Blender
		handler     http.Handler
		assets      fs.FS
		logger      types.Logger

		cancel context.CancelFunc

		ui        *wailsUI
		lifecycle *lifecycle

		project   *Project
		pkg       *Package
		operation *Operation
		system    *System
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

	deps, err := setupDependencies(container)
	if err != nil {
		return nil, fmt.Errorf("failed to setup dependencies: %w", err)
	}

	rocketblend, err := container.GetRBDriver()
	if err != nil {
		return nil, err
	}

	blender, err := container.GetBlender()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	ui := newWailsUI()

	return &Application{
		id:          id,
		platform:    platform,
		assets:      opts.Assets,
		rocketblend: rocketblend,
		blender:     blender,
		handler:     handler,
		logger:      deps.logger,
		cancel:      cancel,
		ui:          ui,
		lifecycle:   newLifecycle(ctx, deps, ui, events, opts.Args),
		project:     newProject(ctx, deps),
		pkg:         newPackage(ctx, deps),
		operation:   newOperation(ctx, deps),
		system:      newSystem(ctx, deps, ui, opts.Version, platform),
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
		Windows: &windows.Options{},
		Linux: &linux.Options{
			ProgramName: title,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.DefaultAppearance,
			WebviewIsTransparent: true,
			About: &mac.AboutInfo{
				Title:   title,
				Message: fmt.Sprintf("© %d RocketBlend. All rights reserved.", time.Now().Year()),
			},
			OnFileOpen: func(filePath string) {
				if err := a.Open(context.Background(), filePath); err != nil {
					a.logger.Error("failed to open blend file", map[string]interface{}{"error": err})
				}
			},
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: a.id.String(),
			OnSecondInstanceLaunch: func(data options.SecondInstanceData) {
				a.lifecycle.onSecondInstanceLaunch(data.Args, data.WorkingDirectory)
			},
		},
		EnumBind: []interface{}{
			enums.PackageStates,
			enums.PackageTypes,
		},
		MinHeight:        580,
		MinWidth:         800,
		BackgroundColour: &options.RGBA{R: 00, G: 00, B: 00, A: 1},
		OnStartup:        a.startup,
		OnShutdown:       a.shutdown,
		OnDomReady:       a.domReady,
		Frameless:        frameless,
		Bind: []interface{}{
			a.project,
			a.pkg,
			a.operation,
			a.system,
		},
	})
}

func (a *Application) startup(ctx context.Context) {
	a.ui.startup(ctx)
	a.lifecycle.startup()
}

func (a *Application) shutdown(ctx context.Context) {
	a.cancel()
	a.lifecycle.shutdown()
}

func (a *Application) domReady(ctx context.Context) {
	a.lifecycle.onDomReady()
}

func (a *Application) Open(ctx context.Context, filePath string) error {
	if err := openWithRocketBlend(ctx, a.rocketblend, a.blender, filePath); err != nil {
		return err
	}

	return nil
}
