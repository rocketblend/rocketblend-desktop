package store

import (
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/events"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/store/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

func newEvent(id uuid.UUID, indexType indextype.IndexType) types.Eventer {
	return &events.StoreEvent{
		ID:        id,
		IndexType: indexType,
	}
}
