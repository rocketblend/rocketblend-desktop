package projectservice

import (
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/eventservice"
)

const (
	RunEventChannel      = "projectservice.run"
	ExoploreEventChannel = "projectservice.explore"
)

type (
	Event struct {
		eventservice.Event

		ID uuid.UUID `json:"id"`
	}
)

func NewEvent(id uuid.UUID) eventservice.Eventer {
	return &Event{
		ID: id,
	}
}
