package configuration

import (
	"errors"

	"github.com/spf13/viper"
)

const (
	defaultConfigPath = "."
	defaultConfigName = "config"
	defaultConfigType = "yaml"
)

var current *Configuration // singleton

type Configuration struct {
	Postgres Postgres
}

func newEmptyConfig() *Configuration {
	return &Configuration{}
}

func (c Configuration) Validate() error {
	if err := c.Postgres.Validate(); err != nil {
		return err
	}
	return nil
}

func Get() Configuration {
	if current == nil {
		load()
	}
	return *current
}

func load() {
	c := newEmptyConfig()
	var root *viper.Viper
	var err error
	if root, err = viperSetup(defaultConfigPath, defaultConfigName, defaultConfigType); err != nil {
		panic(errors.Join(ErrorLoadingConfiguration, err))
	}
	if err = setValues(root, c); err != nil {
		panic(errors.Join(ErrorLoadingConfiguration, err))
	}
	current = c
}
