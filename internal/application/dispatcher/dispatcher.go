package dispatcher

import (
	"context"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
)

type (
	eventContext struct {
		name   string
		cancel context.CancelFunc
	}

	Dispatcher struct {
		logger logger.Logger

		sync.RWMutex

		events   map[string][]types.EventListener
		register sync.Map
	}

	Options struct {
		Logger logger.Logger
	}

	Option func(*Options)
)

// WithLogger sets the logger on the Options struct.
func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

// New creates a new event service
func New(opts ...Option) (*Dispatcher, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	return &Dispatcher{
		logger: options.Logger,
		events: make(map[string][]types.EventListener),
	}, nil
}

// EventExists checks if an event exists
func (d *Dispatcher) EventExists(name string) bool {
	d.RLock()
	defer d.RUnlock()
	_, ok := d.events[name]
	return ok
}

// ListEvents returns a list of all events
func (d *Dispatcher) ListEvents() []string {
	d.RLock()
	defer d.RUnlock()
	list := make([]string, 0, len(d.events))
	for name := range d.events {
		list = append(list, name)
	}
	return list
}

// FilterEvents returns a list of events filtered by a filter function
func (d *Dispatcher) FilterEvents(filterFunc func(string, []types.EventListener) bool) []string {
	d.RLock()
	defer d.RUnlock()
	var filteredEvents []string
	for name, listeners := range d.events {
		if filterFunc(name, listeners) {
			filteredEvents = append(filteredEvents, name)
		}
	}

	return filteredEvents
}

// CountListeners returns the number of listeners for an event
func (d *Dispatcher) CountListeners(eventName string) int {
	d.RLock()
	defer d.RUnlock()
	if eventName != "" {
		return len(d.events[eventName])
	}
	count := 0
	for _, listeners := range d.events {
		count += len(listeners)
	}

	return count
}

// Close closes the event service
func (d *Dispatcher) Close() error {
	// TODO: Wait for all listeners to finish exiting.
	return nil
}

func (d *Dispatcher) generateListenerID() string {
	return uuid.New().String()
}
