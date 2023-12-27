package application

import "time"

type (
	LaunchEvent struct {
		Args []string `json:"args"`
	}

	LogEvent struct {
		Level   string                 `json:"level"`
		Message string                 `json:"message"`
		Time    time.Time              `json:"time"`
		Fields  map[string]interface{} `json:"fields"`
	}
)
