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

func (c *Container) GetRBDriver() (rbtypes.Driver, error) {
	driver, err := c.rbContainer.GetDriver()
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func (c *Container) GetBlender() (rbtypes.Blender, error) {
	blender, err := c.rbContainer.GetBlender()
	if err != nil {
		return nil, err
	}

	return blender, nil
}
