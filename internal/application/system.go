package application

import (
	"context"

	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	rbruntime "github.com/rocketblend/rocketblend/pkg/runtime"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

type (
	System struct {
		logger types.Logger

		tracker   types.Tracker
		portfolio types.Portfolio

		configurator   types.Configurator
		rbConfigurator rbtypes.Configurator

		ctx      context.Context
		ui       UI
		version  string
		platform rbruntime.Platform
	}
)

func newSystem(ctx context.Context, deps *dependencies, ui UI, version string, platform rbruntime.Platform) *System {
	return &System{
		logger:         deps.logger,
		tracker:        deps.tracker,
		portfolio:      deps.portfolio,
		configurator:   deps.configurator,
		rbConfigurator: deps.rbConfigurator,
		ctx:            ctx,
		ui:             ui,
		version:        version,
		platform:       platform,
	}
}

// Quit quits the application
func (s *System) Quit() {
	s.logger.Debug("quitting application")
	s.ui.Quit()
}
