package config

type (
	Config struct {
		LogLevel          string   `mapstructure:"logLevel"`
		ProjectWatchPaths []string `mapstructure:"projectWatchPaths"`
	}
)
