package container

import (
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

func (c *Container) GetRBConfigurator() (rbtypes.Configurator, error) {
	configurator, err := c.rbContainer.GetConfigurator()
	if err != nil {
		return nil, err
	}

	return configurator, nil
}
