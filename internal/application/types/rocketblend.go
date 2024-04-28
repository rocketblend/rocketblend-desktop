package types

import "github.com/rocketblend/rocketblend/pkg/types"

type (
	Logger interface {
		types.Logger
	}

	Validator interface {
		types.Validator
	}

	RBConfigurator interface {
		types.Configurator
	}

	RBRepository interface {
		types.Repository
	}

	Blender interface {
		types.Blender
	}

	RBDriver interface {
		types.Driver
	}
)
