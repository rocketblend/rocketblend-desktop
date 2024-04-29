package events

import (
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/indextype"
)

const (
	StoreInsertChannel = "store.insert"
	StoreRemoveChannel = "store.remove"
)

type (
	StoreEvent struct {
		Event

		ID        uuid.UUID           `json:"id"`
		IndexType indextype.IndexType `json:"indexType"`
	}
)
