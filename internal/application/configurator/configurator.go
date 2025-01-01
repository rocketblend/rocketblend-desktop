package configurator

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/mitchellh/mapstructure"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	"github.com/rocketblend/rocketblend/pkg/validator"
	"github.com/spf13/viper"
)

const (
	FileName      = "config"
	FileExtension = "json"
)

type (
	Options struct {
		Logger    types.Logger
		Validator types.Validator
		Path      string
		Name      string
		Extension string
	}

	Option func(*Options)

	Configurator struct {
		logger    types.Logger
		validator types.Validator
		viper     *viper.Viper
		path      string
		extension string
		name      string
	}
)

func WithLogger(logger types.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithValidator(validator types.Validator) Option {
	return func(o *Options) {
		o.Validator = validator
	}
}

func WithLocation(path string) Option {
	return func(o *Options) {
		o.Path = path
	}
}

func WithApplication(name string, extenstion string) Option {
	return func(o *Options) {
		o.Name = name
		o.Extension = extenstion
	}
}

func New(opts ...Option) (*Configurator, error) {
	options := &Options{
		Logger:    logger.NoOp(),
		Validator: validator.New(),
		Extension: FileExtension,
		Name:      FileName,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Path == "" {
		return nil, errors.New("path is required")
	}

	viper, err := load(options.Path, options.Name, options.Extension)
	if err != nil {
		return nil, err
	}

	return &Configurator{
		logger:    options.Logger,
		validator: options.Validator,
		path:      options.Path,
		extension: options.Extension,
		name:      options.Name,
		viper:     viper,
	}, nil
}

func (c *Configurator) Get() (*types.Config, error) {
	var config types.Config
	if err := c.viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := c.validator.Validate(config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Configurator) GetAllValues() map[string]interface{} {
	return c.viper.AllSettings()
}

func (c *Configurator) GetValueByString(key string) string {
	return fmt.Sprint(c.viper.Get(key))
}

func (c *Configurator) SetValueByString(key string, value string) error {
	c.viper.Set(key, value)

	_, err := c.Get()
	if err != nil {
		return err
	}

	c.viper.WriteConfig()

	return nil
}

func (c *Configurator) Save(config *types.Config) error {
	if err := c.validator.Validate(config); err != nil {
		return err
	}

	var m map[string]interface{}
	mapstructure.Decode(config, &m)

	c.viper.MergeConfigMap(m)

	return c.viper.WriteConfig()
}

func (c *Configurator) Path() string {
	return fmt.Sprintf("%s.%s", filepath.Join(c.path, c.name), c.extension)
}

func load(path string, name string, extension string) (*viper.Viper, error) {
	v := viper.New()

	v.SetDefault("package.autoPull", true)

	v.SetConfigName(name)
	v.AddConfigPath(path)
	v.SetConfigType(extension)

	v.SafeWriteConfig()

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return v, nil
}
