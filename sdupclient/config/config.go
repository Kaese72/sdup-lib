package config

import (
	"errors"
)

type Config struct {
	HTTP *HTTPConfig `json:"http,omitempty"`
}

func (config *Config) PopulateExample() {
	config.HTTP = &HTTPConfig{}
	config.HTTP.PopulateExample()
}

func (config Config) Validate() error {
	if config.HTTP != nil {
		return config.HTTP.Validate()

	} else {
		return errors.New("at least one config type must be set")
	}
}
