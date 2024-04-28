package events

type (
	Event struct {
		stopped bool
	}
)

// StopPropagation Stops the propagation of the event to further event listeners
func (e *Event) StopPropagation() {
	e.stopped = true
}

// IsPropagationStopped returns whether further event listeners should be triggered
func (e *Event) IsPropagationStopped() bool {
	return e.stopped
}
