package eventservice

import (
	"errors"
	"reflect"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
)

type (
	Service interface {
		On(name string, fn interface{}) error
		Go(name string, params ...interface{}) error
		Has(name string) bool
		List() []string
		Remove(names ...string)
	}

	service struct {
		logger logger.Logger

		sync.RWMutex

		events map[string][]interface{}
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
		events: make(map[string][]interface{}),
	}, nil
}

// On set new listener
func (s *service) On(name string, fn interface{}) error {
	s.Lock()
	defer s.Unlock()

	if fn == nil {
		err := errors.New("fn is nil")
		s.logger.Error("fn is nil", map[string]interface{}{"error": err})
		return err
	}
	if _, ok := fn.(handle); ok {
		s.events[name] = append(s.events[name], fn)
		return nil
	}

	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		return errors.New("fn is not a function")
	}
	if t.NumOut() != 1 {
		return errors.New("fn must have one return value")
	}
	if t.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return errors.New("fn must return an error message")
	}
	if list, ok := s.events[name]; ok && len(list) > 0 {
		tt := reflect.TypeOf(list[0])
		if tt.NumIn() != t.NumIn() {
			return errors.New("fn signature is not equal")
		}
		for i := 0; i < tt.NumIn(); i++ {
			if tt.In(i) != t.In(i) {
				return errors.New("fn signature is not equal")
			}
		}
	}

	s.events[name] = append(s.events[name], fn)
	s.logger.Info("listener registered", map[string]interface{}{"event": name})
	return nil
}

// Go firing an event
func (s *service) Go(name string, params ...interface{}) error {
	s.logger.Trace("Firing event", map[string]interface{}{"event": name})

	s.RLock()
	defer s.RUnlock()

	fns := s.events[name]
	for i := len(fns) - 1; i >= 0; i-- {
		stopped, err := s.call(fns[i], params...)
		if err != nil {
			s.logger.Error("error in event handling", map[string]interface{}{"event": name, "error": err})
			return err
		}
		if stopped {
			break
		}
	}

	return nil
}

// Has returns true if a event exists
func (s *service) Has(name string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.events[name]
	return ok
}

// List returns list events
func (s *service) List() []string {
	s.RLock()
	defer s.RUnlock()
	list := make([]string, 0, len(s.events))
	for name := range s.events {
		list = append(list, name)
	}
	return list
}

// Remove delete events from the event list
func (s *service) Remove(names ...string) {
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
	s.events = make(map[string][]interface{})
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
