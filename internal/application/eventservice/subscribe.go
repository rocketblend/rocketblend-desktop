package eventservice

import (
	"context"
	"errors"
	"reflect"
)

// Subscribe registers a listener for an event with optional max call limit
func (s *service) Subscribe(ctx context.Context, name string, fn interface{}, maxTrigger int) (context.CancelFunc, error) {
	listenerID, err := s.subscribe(name, fn, maxTrigger)
	if err != nil {
		return nil, err
	}

	return s.setupListenerCancellation(ctx, listenerID, name), nil
}

func (s *service) setupListenerCancellation(ctx context.Context, listenerID, name string) context.CancelFunc {
	lctx, lcancel := context.WithCancel(ctx)
	s.register.Store(listenerID, listener{name: name, cancel: lcancel})

	go s.listenerCleanupRoutine(lctx, listenerID)

	return lcancel
}

func (s *service) listenerCleanupRoutine(ctx context.Context, listenerID string) {
	<-ctx.Done()
	s.unsubscribe(listenerID)
}

func (s *service) subscribe(name string, fn interface{}, maxTrigger int) (string, error) {
	if err := s.validateFunction(fn); err != nil {
		return "", err
	}

	s.Lock()
	defer s.Unlock()

	id := s.generateListenerID()
	count := s.initializeCount(maxTrigger)

	if err := s.checkFunctionCompatibility(name, fn); err != nil {
		return "", err
	}

	s.events[name] = append(s.events[name], eventListener{id: id, fn: fn, count: count})
	s.logger.Debug("listener registered", map[string]interface{}{"event": name, "id": id})

	return id, nil
}

func (s *service) validateFunction(fn interface{}) error {
	if fn == nil {
		return errors.New("fn is nil")
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

	return nil
}

func (s *service) initializeCount(maxTrigger int) *int {
	if maxTrigger > 0 {
		count := new(int)
		*count = maxTrigger
		return count
	}

	return nil
}

func (s *service) checkFunctionCompatibility(name string, fn interface{}) error {
	if list, ok := s.events[name]; ok && len(list) > 0 {
		t := reflect.TypeOf(fn)
		tt := reflect.TypeOf(list[0].fn)
		if tt.NumIn() != t.NumIn() {
			return errors.New("fn signature mismatch: number of parameters does not match")
		}
		for i := 0; i < tt.NumIn(); i++ {
			if tt.In(i) != t.In(i) {
				return errors.New("fn signature mismatch: parameter types do not match")
			}
		}
	}

	return nil
}

func (s *service) unsubscribe(id string) {
	value, exists := s.register.Load(id)
	if !exists {
		return
	}

	registered, ok := value.(listener)
	if !ok {
		s.logger.Error("Type assertion failed for listener", map[string]interface{}{"id": id})
		return
	}

	// Remove the listener from the register map
	s.register.Delete(id)

	s.Lock()
	defer s.Unlock()

	listeners, ok := s.events[registered.name]
	if !ok {
		// This should not happen in a consistent state, but log if it does
		s.logger.Error("No event found for registered listener", map[string]interface{}{"event": registered.name})
		return
	}

	for i, listener := range listeners {
		if listener.id == id {
			s.events[registered.name][i] = s.events[registered.name][len(s.events[registered.name])-1]
			s.events[registered.name] = s.events[registered.name][:len(s.events[registered.name])-1]
			break
		}
	}

	// If there are no more listeners for this event, delete the event entry
	if len(s.events[registered.name]) == 0 {
		delete(s.events, registered.name)
	}

	s.logger.Debug("listener unsubscribed", map[string]interface{}{"event": registered.name, "id": id})
}
