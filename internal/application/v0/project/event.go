package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/eventservice"
)

const (
	CreateEventChannel = "project.create"
	UpdateEventChannel = "project.update"

	RunEventChannel     = "project.run"
	ExploreEventChannel = "project.explore"
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

func (r *repository) emitEvent(ctx context.Context, id uuid.UUID, channel string) {
	event := NewEvent(id)
	if err := r.dispatcher.EmitEvent(ctx, channel, event); err != nil {
		r.logger.Error("error emitting event", map[string]interface{}{
			"error":   err,
			"event":   event,
			"channel": channel,
		})
	}
}
