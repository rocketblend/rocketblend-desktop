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

func newEvent(id uuid.UUID) eventservice.Eventer {
	return &Event{
		ID: id,
	}
}

func (r *Repository) emitEvent(ctx context.Context, id uuid.UUID, channel string) {
	event := newEvent(id)
	if err := r.dispatcher.EmitEvent(ctx, channel, event); err != nil {
		r.logger.Error("error emitting event", map[string]interface{}{
			"error":   err,
			"event":   event,
			"channel": channel,
		})
	}
}
