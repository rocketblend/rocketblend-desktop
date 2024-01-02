package eventservice

type (
	Eventer interface {
		StopPropagation()
		IsPropagationStopped() bool
	}

	Event struct {
		stopped bool
	}

	handle = func(Eventer) error
)

// StopPropagation Stops the propagation of the event to further event listeners
func (e *Event) StopPropagation() {
	e.stopped = true
}

// IsPropagationStopped returns whether further event listeners should be triggered
func (e *Event) IsPropagationStopped() bool {
	return e.stopped
}
