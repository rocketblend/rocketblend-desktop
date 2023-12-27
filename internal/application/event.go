package application

type (
	LaunchEvent struct {
		Args []string `json:"args"`
	}

	LogEvent struct {
		Level   string
		Message string
		Fields  []map[string]interface{}
	}
)
