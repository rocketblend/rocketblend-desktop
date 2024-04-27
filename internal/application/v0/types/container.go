package types

import (
	"github.com/rocketblend/rocketblend-desktop/internal/application/eventservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/metricservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/operationservice"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Container interface {
		GetApplicationConfigurator() (Configurator, error)

		GetDispatcher() (eventservice.Service, error)
		GetTracker() (metricservice.Service, error)

		GetStore() (Store, error)
		GetPortfolio() (Portfolio, error)
		GetCatalog() (Catalog, error)
		GetOperator() (operationservice.Service, error)

		Preload() error
		Close() error

		rbtypes.Container
	}
)
