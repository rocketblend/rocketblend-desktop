package application

type (
	LaunchEvent struct {
		Args []string `json:"args"`
	}

	LogEvent struct {
		Level   string                 `json:"level"`
		Message string                 `json:"message"`
		Fields  map[string]interface{} `json:"fields"`
	}
)
