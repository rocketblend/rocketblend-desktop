package config

type (
	Watcher struct {
		FileExtensions []string `mapstructure:"fileExtensions"`
	}

	Project struct {
		Paths   []string `mapstructure:"paths"`
		Watcher Watcher  `mapstructure:"watcher"`
	}

	Package struct {
		Watcher Watcher `mapstructure:"watcher"`
	}

	Feature struct {
		Addon     bool `mapstructure:"addon"`
		Terminal  bool `mapstructure:"terminal"`
		Developer bool `mapstructure:"developer"`
	}

	Config struct {
		Project Project `mapstructure:"project"`
		Package Package `mapstructure:"package"`
		Feature Feature `mapstructure:"feature"`
	}
)
