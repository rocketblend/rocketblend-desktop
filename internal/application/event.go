package application

type (
	LaunchEvent struct {
		Args     []string `json:"args"`
		Messages []string `json:"messages"`
	}
)
