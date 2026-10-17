package application

import (
	"context"
	"errors"

	"github.com/rocketblend/rocketblend-desktop/internal/application/events"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	"github.com/rocketblend/rocketblend-desktop/internal/eventwriter"
)

const (
	ApplicationLogChannel      = "application.log"
	ApplicationArgumentChannel = "application.argument"
)

type (
	LaunchEvent struct {
		Args []string `json:"args"`
	}
)

// TODO: We could also create wrapper structs over rocketblend functions to add events to them.
func (l *lifecycle) setupDriverEventHandlers() error {
	if err := l.ctx.Err(); err != nil {
		return err
	}

	if err := l.subscribeToEvent(events.StoreInsertChannel, l.handleStoreInsertEvent); err != nil {
		return err
	}

	if err := l.subscribeToEvent(events.StoreRemoveChannel, l.handleStoreRemoveEvent); err != nil {
		return err
	}

	if err := l.subscribeToEvent(events.ProjectRunChannel, l.handleProjectRunEvent); err != nil {
		return err
	}

	return nil
}

func (l *lifecycle) subscribeToEvent(channel string, handler func(types.Eventer) error) error {
	if err := l.ctx.Err(); err != nil {
		return err
	}

	_, err := l.dispatcher.Subscribe(l.ctx, channel, handler, 0)
	if err != nil {
		l.logger.Error("failed to subscribe to event", map[string]interface{}{
			"channel": channel,
			"error":   err.Error(),
		})

		return err
	}

	return nil
}

func (l *lifecycle) handleProjectRunEvent(e types.Eventer) error {
	ev, ok := e.(*events.ProjectEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	// TODO: Move the metric creation into the portfolio functions???
	if err := l.tracker.CreateMetric(context.Background(), &types.CreateMetricOpts{
		Domain: ev.ID.String(),
		Name:   ProjectRunMetric,
		Value:  1,
	}); err != nil {
		return err
	}

	return nil
}

func (l *lifecycle) handleStoreInsertEvent(e types.Eventer) error {
	ev, ok := e.(*events.StoreEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	l.ui.Emit(events.StoreInsertChannel, ev)
	return nil
}

func (l *lifecycle) handleStoreRemoveEvent(e types.Eventer) error {
	ev, ok := e.(*events.StoreEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	l.ui.Emit(events.StoreRemoveChannel, ev)
	return nil
}

func (l *lifecycle) listenToLogEvents() {
	for {
		select {
		case <-l.ctx.Done():
			return
		case data, ok := <-l.events.Channel():
			if !ok {
				return
			}

			if logEvent, isLogEvent := data.(eventwriter.Event); isLogEvent {
				l.ui.Emit(ApplicationLogChannel, logEvent)
			}
		}
	}
}

func (l *lifecycle) eventEmitLaunchArgs(event LaunchEvent) {
	l.logger.Debug("emitting application.argument event", map[string]interface{}{
		"event": event,
	})

	l.ui.Emit(ApplicationArgumentChannel, event)
}
