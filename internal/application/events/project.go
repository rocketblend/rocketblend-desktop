package events

import (
	"github.com/google/uuid"
)

const (
	ProjectCreateChannel = "project.create"
	ProjectUpdateChannel = "project.update"

	ProjectRunChannel = "project.run"
)

type (
	ProjectEvent struct {
		Event
		ID uuid.UUID `json:"id"`
	}
)
