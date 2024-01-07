package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/eventservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/metricservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	LaunchEvent struct {
		Args []string `json:"args"`
	}

	LogEvent struct {
		Level   string                 `json:"level"`
		Message string                 `json:"message"`
		Time    time.Time              `json:"time"`
		Fields  map[string]interface{} `json:"fields"`
	}
)

func (d *Driver) ListMetrics(options metricservice.FilterOptions) ([]*metricservice.Metric, error) {
	metric, err := d.factory.GetMetricService()
	if err != nil {
		return nil, err
	}

	return metric.List(context.Background(), options)
}

func (d *Driver) AggregateMetrics(options metricservice.FilterOptions) (*metricservice.Aggregate, error) {
	metric, err := d.factory.GetMetricService()
	if err != nil {
		return nil, err
	}

	return metric.Aggregate(context.Background(), options)
}

func (d *Driver) setupDriverEventHandlers() error {
	if err := d.ctx.Err(); err != nil {
		return err
	}

	eventService, err := d.factory.GetEventService()
	if err != nil {
		d.logger.Error("failed to get event service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := d.subscribeToEvent(eventService, searchstore.InsertEventChannel, d.handleStoreInsertEvent); err != nil {
		return err
	}

	if err := d.subscribeToEvent(eventService, projectservice.RunEventChannel, d.handleProjectRunEvent); err != nil {
		return err
	}

	return nil
}

func (d *Driver) subscribeToEvent(eventService eventservice.Service, channel string, handler func(eventservice.Eventer) error) error {
	if err := d.ctx.Err(); err != nil {
		return err
	}

	_, err := eventService.Subscribe(d.ctx, channel, handler, 0)
	if err != nil {
		d.logger.Error("failed to subscribe to event", map[string]interface{}{
			"channel": channel,
			"error":   err.Error(),
		})
		return err
	}

	return nil
}

func (d *Driver) handleProjectRunEvent(e eventservice.Eventer) error {
	ev, ok := e.(*projectservice.Event)
	if !ok {
		return errors.New("invalid event type")
	}

	metric, err := d.factory.GetMetricService()
	if err != nil {
		return err
	}

	if err := metric.Add(context.Background(), metricservice.AddOptions{
		Domain: ev.ID.String(),
		Name:   ProjectRunMetric,
		Value:  1,
	}); err != nil {
		return err
	}

	return nil
}

func (d *Driver) handleStoreInsertEvent(e eventservice.Eventer) error {
	ev, ok := e.(*searchstore.Event)
	if !ok {
		return errors.New("invalid event type")
	}

	runtime.EventsEmit(d.ctx, searchstore.InsertEventChannel, ev)
	return nil
}

func (d *Driver) setupRuntimeEventHandlers() error {
	if err := d.ctx.Err(); err != nil {
		return err
	}

	runtime.EventsOn(d.ctx, "operation.cancel", d.handleOperationCancel)
	return nil
}

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
				if logEvent, isLogEvent := data.(LogEvent); isLogEvent {
					runtime.EventsEmit(d.ctx, "debug.log", logEvent)
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

	runtime.EventsEmit(ctx, "application.argument", event)
}
