package eventservice

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
)

type (
	Service interface {
		SubscribeOnce(ctx context.Context, name string, fn interface{}) error
		Subscribe(ctx context.Context, name string, fn interface{}) error
		Unsubscribe(id string)
		TriggerEvent(ctx context.Context, name string, params ...interface{}) error
		Broadcast(ctx context.Context, params ...interface{}) error
		EventExists(name string) bool
		ListEvents() []string
		FilterEvents(filterFunc func(string, []eventListener) bool) []string
		CountListeners(eventName string) int
		RemoveEvents(names ...string)
	}

	eventListener struct {
		id   string
		fn   interface{}
		once bool
	}

	service struct {
		logger logger.Logger

		sync.RWMutex

		events    map[string][]eventListener
		listeners sync.Map
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

// SubscribeOnce registers a listener for an event that is called only once
func (s *service) SubscribeOnce(ctx context.Context, name string, fn interface{}) error {
	unregister, err := s.subscribe(ctx, name, fn, true)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		unregister()
	}()
	return nil
}

// Subscribe registers a listener for an event
func (s *service) Subscribe(ctx context.Context, name string, fn interface{}) error {
	unregister, err := s.subscribe(ctx, name, fn, false)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		unregister()
	}()
	return nil
}

// Unsubscribe removes a listener from the event list
func (s *service) Unsubscribe(id string) {
	value, exists := s.listeners.Load(id)
	if !exists {
		return
	}

	name, _ := value.(string)
	s.listeners.Delete(id)

	s.Lock()
	defer s.Unlock()

	listeners := s.events[name]
	for i, listener := range listeners {
		if listener.id == id {
			s.events[name] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}

	if len(s.events[name]) == 0 {
		delete(s.events, name)
	}
}

// TriggerEvent fires an event
func (s *service) TriggerEvent(ctx context.Context, name string, params ...interface{}) error {
	s.logger.Trace("firing event", map[string]interface{}{"event": name})

	s.RLock()
	fns := append([]eventListener(nil), s.events[name]...)
	s.RUnlock()

	var err error

	// Temporary slice to store remaining listeners
	remainingListeners := make([]eventListener, 0, len(fns))

eventLoop:
	for _, fn := range fns {
		select {
		case <-ctx.Done():
			s.logger.Debug("event processing canceled", map[string]interface{}{"event": name})
			return ctx.Err()
		default:
			stopped, e := s.call(fn.fn, params...)
			if e != nil {
				s.logger.Error("error in event handling", map[string]interface{}{"event": name, "error": e.Error()})
				err = e
				break eventLoop
			}
			if !fn.once {
				remainingListeners = append(remainingListeners, fn)
			}
			if stopped {
				s.logger.Trace("event propagation stopped", map[string]interface{}{"event": name})
				break eventLoop
			}
		}
	}

	// Update the event listeners, removing 'once' listeners that were executed
	s.Lock()
	s.events[name] = remainingListeners
	s.Unlock()

	return err
}

// Broadcast fires an event to all listeners
func (s *service) Broadcast(ctx context.Context, params ...interface{}) error {
	s.RLock()
	eventNames := make([]string, 0, len(s.events))
	for name := range s.events {
		eventNames = append(eventNames, name)
	}
	s.RUnlock()

	for _, name := range eventNames {
		if err := s.TriggerEvent(ctx, name, params...); err != nil {
			return err
		}
	}

	return nil
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

// RemoveEvents removes one or more events
func (s *service) RemoveEvents(names ...string) {
	s.Lock()
	defer s.Unlock()

	if len(names) > 0 {
		for _, name := range names {
			s.logger.Debug("removing event", map[string]interface{}{"event": name})
			delete(s.events, name)
		}
		return
	}

	s.logger.Debug("clearing all events")
	s.events = make(map[string][]eventListener)
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

func (s *service) subscribe(ctx context.Context, name string, fn interface{}, once bool) (func(), error) {
	if fn == nil {
		return nil, errors.New("fn is nil")
	}

	s.Lock()
	defer s.Unlock()

	if _, ok := fn.(handle); ok {
		id := s.generateListenerID()
		s.events[name] = append(s.events[name], eventListener{id: id, fn: fn, once: once})
		s.listeners.Store(id, name)

		return func() { s.Unsubscribe(id) }, nil
	}

	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		return nil, errors.New("fn is not a function")
	}
	if t.NumOut() != 1 {
		return nil, errors.New("fn must have one return value")
	}
	if t.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil, errors.New("fn must return an error message")
	}
	if list, ok := s.events[name]; ok && len(list) > 0 {
		tt := reflect.TypeOf(list[0])
		if tt.NumIn() != t.NumIn() {
			return nil, errors.New("fn signature is not equal")
		}
		for i := 0; i < tt.NumIn(); i++ {
			if tt.In(i) != t.In(i) {
				return nil, errors.New("fn signature is not equal")
			}
		}
	}

	id := s.generateListenerID()
	s.events[name] = append(s.events[name], eventListener{id: id, fn: fn, once: once})
	s.listeners.Store(id, name)

	s.logger.Info("listener registered", map[string]interface{}{"event": name, "id": id})
	return func() { s.Unsubscribe(id) }, nil
}

func (s *service) call(fn interface{}, params ...interface{}) (stopped bool, err error) {
	if f, ok := fn.(handle); ok {
		if len(params) != 1 {
			return stopped, errors.New("parameters mismatched")
		}
		event, ok := (params[0]).(Eventer)
		if !ok {
			return stopped, errors.New("parameters mismatched")
		}
		err = f(event)
		return event.IsPropagationStopped(), err
	}

	var (
		f     = reflect.ValueOf(fn)
		t     = f.Type()
		numIn = t.NumIn()
		in    = make([]reflect.Value, 0, numIn)
	)

	if t.IsVariadic() {
		n := numIn - 1
		if len(params) < n {
			return stopped, errors.New("parameters mismatched")
		}
		for _, param := range params[:n] {
			in = append(in, reflect.ValueOf(param))
		}
		s := reflect.MakeSlice(t.In(n), 0, len(params[n:]))
		for _, param := range params[n:] {
			s = reflect.Append(s, reflect.ValueOf(param))
		}
		in = append(in, s)

		err, _ = f.CallSlice(in)[0].Interface().(error)
		return stopped, err
	}

	if len(params) != numIn {
		return stopped, errors.New("parameters mismatched")
	}
	for _, param := range params {
		in = append(in, reflect.ValueOf(param))
	}

	err, _ = f.Call(in)[0].Interface().(error)
	if err != nil {
		s.logger.Error("error in function call", map[string]interface{}{"error": err})
		return stopped, err
	}

	s.logger.Trace("function call successful")
	return stopped, nil
}

func (s *service) generateListenerID() string {
	return uuid.New().String()
}
