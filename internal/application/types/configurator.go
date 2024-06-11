package types

type (
	WatcherConfig struct {
		FileExtensions []string `mapstructure:"fileExtensions"`
	}

	ProjectConfig struct {
		Paths   []string      `mapstructure:"paths"`
		Watcher WatcherConfig `mapstructure:"watcher"`
	}

	PackageConfig struct {
		Watcher  WatcherConfig `mapstructure:"watcher"`
		AutoPull bool          `mapstructure:"autoPull"`
	}

	FeatureConfig struct {
		Addon     bool `mapstructure:"addon"`
		Developer bool `mapstructure:"developer"`
	}

	Config struct {
		Project ProjectConfig `mapstructure:"project"`
		Package PackageConfig `mapstructure:"package"`
		Feature FeatureConfig `mapstructure:"feature"`
	}

	Configurator interface {
		Get() (config *Config, err error)
		GetAllValues() map[string]interface{}
		GetValueByString(key string) string
		SetValueByString(key string, value string) error
		Save(config *Config) error
		Path() string
	}
)
