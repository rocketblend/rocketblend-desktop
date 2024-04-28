package container

import (
	"fmt"

	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/dispatcher"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

func (c *Container) GetLogger() (types.Logger, error) {
	return c.logger, nil
}

func (c *Container) GetValidator() (types.Validator, error) {
	return c.validator, nil
}

func (c *Container) GetDispatcher() (types.Dispatcher, error) {
	var err error
	c.dispatcherHolder.once.Do(func() {
		c.dispatcherHolder.instance, err = dispatcher.New(
			dispatcher.WithLogger(c.logger),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get/create dispatcher: %w", err)
	}

	return c.dispatcherHolder.instance, nil
}
