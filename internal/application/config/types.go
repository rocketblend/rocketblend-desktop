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

	Feature struct {
		Addons        bool `mapstructure:"addons"`
		Terminal      bool `mapstructure:"terminal"`
		DeveloperMode bool `mapstructure:"developerMode"`
	}

	Config struct {
		Project Project `mapstructure:"project"`
		Package Package `mapstructure:"package"`
		Feature Feature `mapstructure:"feature"`
	}
)
