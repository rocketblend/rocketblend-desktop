package config

type (
	Watcher struct {
		Paths          []string `mapstructure:"paths"`
		FileExtensions []string `mapstructure:"fileExtensions"`
	}

	Project struct {
		Watcher Watcher `mapstructure:"watcher"`
	}

	Package struct {
	}

	Config struct {
		Project Project `mapstructure:"project"`
		Package Package `mapstructure:"package"`
	}
)
