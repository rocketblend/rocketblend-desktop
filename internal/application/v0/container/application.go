package container

import (
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

func (c *Container) GetConfigurator() (types.Configurator, error) {
	return nil, types.ErrNotImplement
}

func (c *Container) GetStore() (types.Store, error) {
	return nil, types.ErrNotImplement
}

func (c *Container) GetPortfolio() (types.Portfolio, error) {
	return nil, types.ErrNotImplement
}

func (c *Container) GetCatalog() (types.Catalog, error) {
	return nil, types.ErrNotImplement
}
