package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type (
	Service interface {
		Get() (config *Config, err error)
		GetAllValues() map[string]interface{}
		GetValueByString(key string) string
		SetValueByString(key string, value string) error
		Save(config *Config) error
	}

	service struct {
		viper     *viper.Viper
		validator *validator.Validate
		rootPath  string
	}
)

func New(rootPath string) (Service, error) {
	v, err := load(rootPath)
	if err != nil {
		return nil, err
	}

	return &service{
		viper:     v,
		validator: validator.New(),
		rootPath:  rootPath,
	}, nil
}

func (srv *service) Get() (config *Config, err error) {
	err = srv.viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	err = srv.validate(config)
	if err != nil {
		return nil, err
	}

	return config, err
}

func (srv *service) GetAllValues() map[string]interface{} {
	return srv.viper.AllSettings()
}

func (srv *service) GetValueByString(key string) string {
	return fmt.Sprint(srv.viper.Get(key))
}

func (srv *service) SetValueByString(key string, value string) error {
	srv.viper.Set(key, value)

	_, err := srv.Get()
	if err != nil {
		return err
	}

	srv.viper.WriteConfig()

	return nil
}

func (srv *service) Save(config *Config) error {
	err := srv.validate(config)
	if err != nil {
		return err
	}

	var m map[string]interface{}
	mapstructure.Decode(config, &m)

	srv.viper.MergeConfigMap(m)

	return srv.viper.WriteConfig()
}

func (srv *service) validate(config *Config) error {
	if err := srv.validator.Struct(config); err != nil {
		return err
	}

	return nil
}

func load(rootPath string) (*viper.Viper, error) {
	v := viper.New()

	v.SetDefault("project.watcher.fileExtensions", []string{".blend", ".yaml"})

	v.SetConfigName("settings")
	v.AddConfigPath(rootPath)
	v.SetConfigType("json")

	v.SafeWriteConfig()

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return v, nil
}
