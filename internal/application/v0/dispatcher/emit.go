package dispatcher

import (
	"context"
	"errors"
	"reflect"

	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

type (
	handle = func(types.Eventer) error
)

// EmitEvent fires an event
func (d *Dispatcher) EmitEvent(ctx context.Context, name string, params ...interface{}) (err error) {
	d.logger.Trace("firing event", map[string]interface{}{"event": name})

	listeners := d.copyListeners(name)

	for _, listener := range listeners {
		if err = d.processListener(ctx, &listener, name, params...); err != nil {
			break
		}
	}

	d.updateListeners(name, listeners)

	return
}

// Broadcast fires an event to all listeners
func (d *Dispatcher) Broadcast(ctx context.Context, params ...interface{}) error {
	d.RLock()
	eventNames := make([]string, 0, len(d.events))
	for name := range d.events {
		eventNames = append(eventNames, name)
	}
	d.RUnlock()

	for _, name := range eventNames {
		if err := d.EmitEvent(ctx, name, params...); err != nil {
			return err
		}
	}

	return nil
}

func (d *Dispatcher) copyListeners(name string) []types.EventListener {
	d.RLock()
	defer d.RUnlock()

	// Create a copy to avoid modifying the original slice
	copied := make([]types.EventListener, len(d.events[name]))
	copy(copied, d.events[name])

	return copied
}

func (d *Dispatcher) processListener(ctx context.Context, listener *types.EventListener, name string, params ...interface{}) error {
	select {
	case <-ctx.Done():
		d.logger.Debug("event processing canceled", map[string]interface{}{"event": name})
		return ctx.Err()
	default:
		stopped, err := d.call(listener.FN, params...)
		if err != nil {
			d.logger.Error("error in event handling", map[string]interface{}{"event": name, "error": err.Error()})
			return err
		}

		d.updateListenerCount(listener)

		if stopped {
			d.logger.Trace("event propagation stopped", map[string]interface{}{"event": name})
			return errors.New("event propagation stopped")
		}

		return nil
	}
}

func (d *Dispatcher) call(fn interface{}, params ...interface{}) (stopped bool, err error) {
	switch f := fn.(type) {
	case handle:
		return d.callHandle(f, params...)
	default:
		return d.callWithReflection(fn, params...)
	}
}

func (d *Dispatcher) callHandle(f handle, params ...interface{}) (bool, error) {
	if len(params) != 1 {
		return false, errors.New("handle function requires exactly one parameter")
	}

	event, ok := params[0].(types.Eventer)
	if !ok {
		return false, errors.New("parameter is not of type Eventer")
	}

	err := f(event)
	return event.IsPropagationStopped(), err
}

func (d *Dispatcher) callWithReflection(fn interface{}, params ...interface{}) (bool, error) {
	f := reflect.ValueOf(fn)
	t := f.Type()

	if err := d.validateFunctionParams(t, params); err != nil {
		return false, err
	}

	var in []reflect.Value
	if t.IsVariadic() {
		in = d.prepareVariadicParams(t, params)
	} else {
		in = d.prepareNonVariadicParams(t, params)
	}

	return d.invokeFunction(f, in)
}

func (d *Dispatcher) validateFunctionParams(t reflect.Type, params []interface{}) error {
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

func (d *Dispatcher) prepareVariadicParams(t reflect.Type, params []interface{}) []reflect.Value {
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

func (d *Dispatcher) prepareNonVariadicParams(t reflect.Type, params []interface{}) []reflect.Value {
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}

	return in
}

func (d *Dispatcher) invokeFunction(f reflect.Value, in []reflect.Value) (bool, error) {
	results := f.Call(in)

	var err error
	if len(results) > 0 {
		err, _ = results[0].Interface().(error)
	}

	if err != nil {
		d.logger.Error("error in function call", map[string]interface{}{"error": err})
	} else {
		d.logger.Trace("function call successful")
	}

	return false, err
}

func (d *Dispatcher) updateListenerCount(listener *types.EventListener) {
	if listener.Count != nil {
		*listener.Count--
		if *listener.Count <= 0 {
			d.unregisterListener(listener)
		}
	}
}

func (d *Dispatcher) unregisterListener(eventListener *types.EventListener) {
	if val, ok := d.register.Load(eventListener.ID); ok {
		if lstnr, ok := val.(eventContext); ok {
			lstnr.cancel()
		}
	}
}

func (d *Dispatcher) updateListeners(name string, listeners []types.EventListener) {
	d.Lock()
	defer d.Unlock()

	var remainingListeners []types.EventListener
	for _, listener := range listeners {
		if listener.Count == nil || *listener.Count > 0 {
			remainingListeners = append(remainingListeners, listener)
		}
	}

	d.events[name] = remainingListeners
}
