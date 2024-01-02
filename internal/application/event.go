package application

import (
	"context"
	"time"

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
