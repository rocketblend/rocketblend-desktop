package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/events"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	"github.com/rocketblend/rocketblend-desktop/internal/eventwriter"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	OperationCancelChannel = "operation.cancel"

	ApplicationLogChannel      = "application.log"
	ApplicationArgumentChannel = "application.argument"
)

type (
	LaunchEvent struct {
		Args []string `json:"args"`
	}
)

// TODO: Might be better to just create an interface for the runtime.EventsEmit functions and pass that in as a dependency.
// TODO: We could also create wrapper structs over rocketblend functions to add events to them.
func (d *Driver) setupDriverEventHandlers() error {
	if err := d.ctx.Err(); err != nil {
		return err
	}

	if err := d.subscribeToEvent(events.StoreInsertChannel, d.handleStoreInsertEvent); err != nil {
		return err
	}

	if err := d.subscribeToEvent(events.StoreRemoveChannel, d.handleStoreRemoveEvent); err != nil {
		return err
	}

	if err := d.subscribeToEvent(events.ProjectRunChannel, d.handleProjectRunEvent); err != nil {
		return err
	}

	return nil
}

func (d *Driver) subscribeToEvent(channel string, handler func(types.Eventer) error) error {
	if err := d.ctx.Err(); err != nil {
		return err
	}

	_, err := d.dispatcher.Subscribe(d.ctx, channel, handler, 0)
	if err != nil {
		d.logger.Error("failed to subscribe to event", map[string]interface{}{
			"channel": channel,
			"error":   err.Error(),
		})

		return err
	}

	return nil
}

func (d *Driver) handleProjectRunEvent(e types.Eventer) error {
	ev, ok := e.(*events.ProjectEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	// TODO: Move the metric creation into the portfolio functions???
	if err := d.tracker.CreateMetric(context.Background(), &types.CreateMetricOpts{
		Domain: ev.ID.String(),
		Name:   ProjectRunMetric,
		Value:  1,
	}); err != nil {
		return err
	}

	return nil
}

func (d *Driver) handleStoreInsertEvent(e types.Eventer) error {
	ev, ok := e.(*events.StoreEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	runtime.EventsEmit(d.ctx, events.StoreInsertChannel, ev)
	return nil
}

func (d *Driver) handleStoreRemoveEvent(e types.Eventer) error {
	ev, ok := e.(*events.StoreEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	runtime.EventsEmit(d.ctx, events.StoreRemoveChannel, ev)
	return nil
}

func (d *Driver) setupRuntimeEventHandlers() error {
	if err := d.ctx.Err(); err != nil {
		return err
	}

	runtime.EventsOn(d.ctx, OperationCancelChannel, d.handleOperationCancel)
	return nil
}

// TODO: does this need to be an event? Can we just call this via api directly?
func (d *Driver) handleOperationCancel(optionalData ...interface{}) {
	if len(optionalData) == 0 {
		d.logger.Error("no operation ID provided for cancellation")
		return
	}

	cidStr, ok := optionalData[0].(string)
	if !ok {
		d.logger.Error("invalid data type for operation ID", map[string]interface{}{
			"type": fmt.Sprintf("%T", optionalData[0]),
		})
		return
	}

	cid, err := uuid.Parse(cidStr)
	if err != nil {
		d.logger.Error("invalid operation ID format for cancellation", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	d.cancelTokens.Store(cid.String(), true)
	d.logger.Debug("cancellation requested", map[string]interface{}{
		"cid": cid,
	})
}

func (d *Driver) listenToLogEvents() {
	for {
		select {
		case <-d.ctx.Done():
			return
		default:
			data, ok := d.events.GetNextData()
			if ok {
				if logEvent, isLogEvent := data.(eventwriter.Event); isLogEvent {
					runtime.EventsEmit(d.ctx, ApplicationLogChannel, logEvent)
				}
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}

func (d *Driver) eventEmitLaunchArgs(ctx context.Context, event LaunchEvent) {
	d.logger.Debug("emitting application.argument event", map[string]interface{}{
		"event": event,
	})

	runtime.EventsEmit(ctx, ApplicationArgumentChannel, event)
}
