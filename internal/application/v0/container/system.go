package container

import (
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

func (c *Container) GetDispatcher() (types.Dispatcher, error) {
	return nil, types.ErrNotImplement
}

func (c *Container) GetTracker() (types.Tracker, error) {
	return nil, types.ErrNotImplement
}

func (c *Container) GetOperator() (types.Operator, error) {
	return nil, types.ErrNotImplement
}
