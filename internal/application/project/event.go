package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/events"
)

func (r *Repository) emitEvent(ctx context.Context, id uuid.UUID, channel string) {
	event := events.ProjectEvent{ID: id}
	if err := r.dispatcher.EmitEvent(ctx, channel, event); err != nil {
		r.logger.Error("error emitting event", map[string]interface{}{
			"error":   err,
			"event":   event,
			"channel": channel,
		})
	}
}
