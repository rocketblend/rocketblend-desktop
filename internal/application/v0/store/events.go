package store

import (
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/eventservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
)

type (
	Event struct {
		eventservice.Event

		ID        uuid.UUID           `json:"id"`
		IndexType indextype.IndexType `json:"indexType"`
	}
)

func newEvent(id uuid.UUID, indexType indextype.IndexType) eventservice.Eventer {
	return &Event{
		ID:        id,
		IndexType: indexType,
	}
}
