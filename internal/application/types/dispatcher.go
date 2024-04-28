package types

import "context"

type (
	Eventer interface {
		StopPropagation()
		IsPropagationStopped() bool
	}

	EventListener struct {
		ID    string
		Count *int
		FN    interface{}
	}

	Dispatcher interface {
		Subscribe(ctx context.Context, name string, fn interface{}, maxTrigger int) (context.CancelFunc, error)
		EmitEvent(ctx context.Context, name string, params ...interface{}) error
		Broadcast(ctx context.Context, params ...interface{}) error
		EventExists(name string) bool
		ListEvents() []string
		FilterEvents(filterFunc func(string, []EventListener) bool) []string
		CountListeners(eventName string) int
		Close() error
	}
)
