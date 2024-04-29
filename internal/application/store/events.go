package store

import (
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/events"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
)

func newEvent(id uuid.UUID, indexType indextype.IndexType) types.Eventer {
	return &events.StoreEvent{
		ID:        id,
		IndexType: indexType,
	}
}
