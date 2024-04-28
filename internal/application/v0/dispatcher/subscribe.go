package dispatcher

import (
	"context"
	"errors"
	"reflect"

	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

// Subscribe registers a listener for an event with optional max call limit
func (d *dispatcher) Subscribe(ctx context.Context, name string, fn interface{}, maxTrigger int) (context.CancelFunc, error) {
	listenerID, err := d.subscribe(name, fn, maxTrigger)
	if err != nil {
		return nil, err
	}

	return d.setupListenerCancellation(ctx, listenerID, name), nil
}

func (d *dispatcher) setupListenerCancellation(ctx context.Context, listenerID, name string) context.CancelFunc {
	lctx, lcancel := context.WithCancel(ctx)
	d.register.Store(listenerID, eventContext{name: name, cancel: lcancel})

	go d.listenerCleanupRoutine(lctx, listenerID)

	return lcancel
}

func (d *dispatcher) listenerCleanupRoutine(ctx context.Context, listenerID string) {
	<-ctx.Done()
	d.logger.Debug("listener context canceled", map[string]interface{}{"id": listenerID})
	d.unsubscribe(listenerID)
}

func (d *dispatcher) subscribe(name string, fn interface{}, maxTrigger int) (string, error) {
	if err := d.validateFunction(fn); err != nil {
		return "", err
	}

	d.Lock()
	defer d.Unlock()

	id := d.generateListenerID()
	count := d.initializeCount(maxTrigger)

	if err := d.checkFunctionCompatibility(name, fn); err != nil {
		return "", err
	}

	d.events[name] = append(d.events[name], types.EventListener{ID: id, FN: fn, Count: count})
	d.logger.Debug("listener registered", map[string]interface{}{"event": name, "id": id})

	return id, nil
}

func (d *dispatcher) validateFunction(fn interface{}) error {
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

func (d *dispatcher) initializeCount(maxTrigger int) *int {
	if maxTrigger > 0 {
		count := new(int)
		*count = maxTrigger
		return count
	}

	return nil
}

func (d *dispatcher) checkFunctionCompatibility(name string, fn interface{}) error {
	if list, ok := d.events[name]; ok && len(list) > 0 {
		t := reflect.TypeOf(fn)
		tt := reflect.TypeOf(list[0].FN)
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

func (d *dispatcher) unsubscribe(id string) {
	value, exists := d.register.Load(id)
	if !exists {
		return
	}

	registered, ok := value.(eventContext)
	if !ok {
		d.logger.Error("Type assertion failed for listener", map[string]interface{}{"id": id})
		return
	}

	// Remove the listener from the register map
	d.register.Delete(id)

	d.Lock()
	defer d.Unlock()

	listeners, ok := d.events[registered.name]
	if !ok {
		// This should not happen in a consistent state, but log if it does
		d.logger.Error("No event found for registered listener", map[string]interface{}{"event": registered.name})
		return
	}

	for i, listener := range listeners {
		if listener.ID == id {
			d.events[registered.name][i] = d.events[registered.name][len(d.events[registered.name])-1]
			d.events[registered.name] = d.events[registered.name][:len(d.events[registered.name])-1]
			break
		}
	}

	// If there are no more listeners for this event, delete the event entry
	if len(d.events[registered.name]) == 0 {
		delete(d.events, registered.name)
	}

	d.logger.Debug("listener unsubscribed", map[string]interface{}{"event": registered.name, "id": id})
}
