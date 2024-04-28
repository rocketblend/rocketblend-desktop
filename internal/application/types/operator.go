package types

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/listoption"
)

type (
	Operation struct {
		ID        uuid.UUID   `json:"id"`
		Completed bool        `json:"completed"`
		ErrorMsg  string      `json:"error,omitempty"`
		Result    interface{} `json:"result,omitempty"`
	}

	Operator interface {
		Create(ctx context.Context, opFunc func(ctx context.Context, opid uuid.UUID) (interface{}, error)) (uuid.UUID, error)
		Get(ctx context.Context, opid uuid.UUID) (*Operation, error)
		List(ctx context.Context, opts ...listoption.ListOption) ([]*Operation, error)
		Cancel(opid uuid.UUID) error
	}
)
