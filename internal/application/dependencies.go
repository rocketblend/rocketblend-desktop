package application

import (
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

type (
	dependencies struct {
		logger types.Logger

		dispatcher types.Dispatcher
		tracker    types.Tracker
		operator   types.Operator

		portfolio types.Portfolio
		catalog   types.Catalog

		configurator   types.Configurator
		rbConfigurator rbtypes.Configurator
	}
)

func setupDependencies(container types.Container) (*dependencies, error) {
	logger, err := container.GetLogger()
	if err != nil {
		return nil, err
	}

	dispatcher, err := container.GetDispatcher()
	if err != nil {
		return nil, err
	}

	tracker, err := container.GetTracker()
	if err != nil {
		return nil, err
	}

	operator, err := container.GetOperator()
	if err != nil {
		return nil, err
	}

	configurator, err := container.GetConfigurator()
	if err != nil {
		return nil, err
	}

	portfolio, err := container.GetPortfolio()
	if err != nil {
		return nil, err
	}

	catalog, err := container.GetCatalog()
	if err != nil {
		return nil, err
	}

	rbConfigurator, err := container.GetRBConfigurator()
	if err != nil {
		return nil, err
	}

	return &dependencies{
		logger:         logger,
		dispatcher:     dispatcher,
		tracker:        tracker,
		operator:       operator,
		configurator:   configurator,
		rbConfigurator: rbConfigurator,
		portfolio:      portfolio,
		catalog:        catalog,
	}, nil
}
