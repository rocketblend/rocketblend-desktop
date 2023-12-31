package eventservice

import (
	"context"
	"errors"
	"reflect"
)

// TriggerEvent fires an event
func (s *service) TriggerEvent(ctx context.Context, name string, params ...interface{}) (err error) {
	s.logger.Trace("firing event", map[string]interface{}{"event": name})

	listeners := s.copyListeners(name)

	for _, listener := range listeners {
		if err = s.processListener(ctx, &listener, name, params...); err != nil {
			break
		}
	}

	s.updateListeners(name, listeners)

	return
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

func (s *service) copyListeners(name string) []eventListener {
	s.RLock()
	defer s.RUnlock()

	// Create a copy to avoid modifying the original slice
	copied := make([]eventListener, len(s.events[name]))
	copy(copied, s.events[name])

	return copied
}

func (s *service) processListener(ctx context.Context, listener *eventListener, name string, params ...interface{}) error {
	select {
	case <-ctx.Done():
		s.logger.Debug("event processing canceled", map[string]interface{}{"event": name})
		return ctx.Err()
	default:
		stopped, err := s.call(listener.fn, params...)
		if err != nil {
			s.logger.Error("error in event handling", map[string]interface{}{"event": name, "error": err.Error()})
			return err
		}

		s.updateListenerCount(listener)

		if stopped {
			s.logger.Trace("event propagation stopped", map[string]interface{}{"event": name})
			return errors.New("event propagation stopped")
		}

		return nil
	}
}

func (s *service) call(fn interface{}, params ...interface{}) (stopped bool, err error) {
	switch f := fn.(type) {
	case handle:
		return s.callHandle(f, params...)
	default:
		return s.callWithReflection(fn, params...)
	}
}

func (s *service) callHandle(f handle, params ...interface{}) (bool, error) {
	if len(params) != 1 {
		return false, errors.New("handle function requires exactly one parameter")
	}

	event, ok := params[0].(Eventer)
	if !ok {
		return false, errors.New("parameter is not of type Eventer")
	}

	err := f(event)
	return event.IsPropagationStopped(), err
}

func (s *service) callWithReflection(fn interface{}, params ...interface{}) (bool, error) {
	f := reflect.ValueOf(fn)
	t := f.Type()

	if err := s.validateFunctionParams(t, params); err != nil {
		return false, err
	}

	var in []reflect.Value
	if t.IsVariadic() {
		in = s.prepareVariadicParams(t, params)
	} else {
		in = s.prepareNonVariadicParams(t, params)
	}

	return s.invokeFunction(f, in)
}

func (s *service) validateFunctionParams(t reflect.Type, params []interface{}) error {
	expected := t.NumIn()
	if t.IsVariadic() {
		if len(params) < expected-1 {
			return errors.New("insufficient parameters for variadic function")
		}
	} else if len(params) != expected {
		return errors.New("parameter count mismatch for function")
	}

	return nil
}

func (s *service) prepareVariadicParams(t reflect.Type, params []interface{}) []reflect.Value {
	numIn := t.NumIn()
	in := make([]reflect.Value, 0, numIn)

	for _, param := range params[:numIn-1] {
		in = append(in, reflect.ValueOf(param))
	}
	variadicParams := params[numIn-1:]
	slice := reflect.MakeSlice(t.In(numIn-1), 0, len(variadicParams))
	for _, param := range variadicParams {
		slice = reflect.Append(slice, reflect.ValueOf(param))
	}
	in = append(in, slice)

	return in
}

func (s *service) prepareNonVariadicParams(t reflect.Type, params []interface{}) []reflect.Value {
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}

	return in
}

func (s *service) invokeFunction(f reflect.Value, in []reflect.Value) (bool, error) {
	results := f.Call(in)

	var err error
	if len(results) > 0 {
		err, _ = results[0].Interface().(error)
	}

	if err != nil {
		s.logger.Error("error in function call", map[string]interface{}{"error": err})
	} else {
		s.logger.Trace("function call successful")
	}

	return false, err
}

func (s *service) updateListenerCount(listener *eventListener) {
	if listener.count != nil {
		*listener.count--
		if *listener.count <= 0 {
			s.unregisterListener(listener)
		}
	}
}

func (s *service) unregisterListener(eventListener *eventListener) {
	if val, ok := s.register.Load(eventListener.id); ok {
		if lstnr, ok := val.(listener); ok {
			lstnr.cancel()
		}
	}
}

func (s *service) updateListeners(name string, listeners []eventListener) {
	s.Lock()
	defer s.Unlock()

	var remainingListeners []eventListener
	for _, listener := range listeners {
		if listener.count == nil || *listener.count > 0 {
			remainingListeners = append(remainingListeners, listener)
		}
	}

	s.events[name] = remainingListeners
}
