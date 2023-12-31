package eventservice

import (
	"context"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
)

type (
	Service interface {
		Subscribe(ctx context.Context, name string, fn interface{}, maxTrigger int) (context.CancelFunc, error)
		TriggerEvent(ctx context.Context, name string, params ...interface{}) error
		Broadcast(ctx context.Context, params ...interface{}) error
		EventExists(name string) bool
		ListEvents() []string
		FilterEvents(filterFunc func(string, []eventListener) bool) []string
		CountListeners(eventName string) int
	}

	eventListener struct {
		id    string
		count *int
		fn    interface{}
	}

	listener struct {
		name   string
		cancel context.CancelFunc
	}

	service struct {
		logger logger.Logger

		sync.RWMutex

		events   map[string][]eventListener
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
func New(opts ...Option) (Service, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	return &service{
		logger: options.Logger,
		events: make(map[string][]eventListener),
	}, nil
}

// EventExists checks if an event exists
func (s *service) EventExists(name string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.events[name]
	return ok
}

// ListEvents returns a list of all events
func (s *service) ListEvents() []string {
	s.RLock()
	defer s.RUnlock()
	list := make([]string, 0, len(s.events))
	for name := range s.events {
		list = append(list, name)
	}
	return list
}

// FilterEvents returns a list of events filtered by a filter function
func (s *service) FilterEvents(filterFunc func(string, []eventListener) bool) []string {
	s.RLock()
	defer s.RUnlock()
	var filteredEvents []string
	for name, listeners := range s.events {
		if filterFunc(name, listeners) {
			filteredEvents = append(filteredEvents, name)
		}
	}

	return filteredEvents
}

// CountListeners returns the number of listeners for an event
func (s *service) CountListeners(eventName string) int {
	s.RLock()
	defer s.RUnlock()
	if eventName != "" {
		return len(s.events[eventName])
	}
	count := 0
	for _, listeners := range s.events {
		count += len(listeners)
	}

	return count
}

func (s *service) generateListenerID() string {
	return uuid.New().String()
}
