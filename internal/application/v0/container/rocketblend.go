package container

import (
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

func (c *Container) GetRBContainer() (rbtypes.Container, error) {
	return nil, types.ErrNotImplement
}
