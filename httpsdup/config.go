package httpsdup

import (
	"errors"
	"fmt"
)

type Config struct {
	ListenAddress string `json:"listen-address"`
	ListenPort    int    `json:"listen-port"`
}

func (config *Config) PopulateExample() {
	config.ListenAddress = "localhost"
	config.ListenPort = 8080
}

func (config Config) Validate() error {
	//FIXME validate ListenAddress
	if config.ListenAddress == "" {
		return errors.New("empty listen address")
	}

	if config.ListenPort < 0 || config.ListenPort > 65665 {
		return fmt.Errorf("invalid port number, %d", config.ListenPort)
	}
	return nil
}
