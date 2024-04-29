package types

import rbtypes "github.com/rocketblend/rocketblend/pkg/types"

type (
	Container interface {
		GetLogger() (Logger, error)
		GetValidator() (Validator, error)

		GetDispatcher() (Dispatcher, error)
		GetTracker() (Tracker, error)
		GetOperator() (Operator, error)

		GetConfigurator() (Configurator, error)
		GetRBConfigurator() (rbtypes.Configurator, error)

		GetStore() (Store, error)
		GetPortfolio() (Portfolio, error)
		GetCatalog() (Catalog, error)

		// Preload() error
		// Close() error
	}
)
