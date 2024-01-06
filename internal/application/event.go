package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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

	StoreEvent struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}
)

func (d *Driver) setupDriverEventHandlers(ctx context.Context) error {
	return nil
}

func (d *Driver) setupRuntimeEventHandlers(ctx context.Context) error {
	runtime.EventsOn(ctx, "operation.cancel", d.handleOperationCancel)
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
