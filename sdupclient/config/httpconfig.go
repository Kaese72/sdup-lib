package config

import (
	"errors"
	"fmt"
)

type HTTPConfig struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Port   *int   `json:"port"`
}

func (config *HTTPConfig) PopulateExample() {
	port := 8080
	config.Scheme = "https"
	config.Host = "localhost"
	config.Port = &port
}

func (config HTTPConfig) Validate() error {
	// Scheme validation
	if config.Scheme == "" {
		return errors.New("scheme may not be empty")
	}
	if !(config.Scheme == "https" || config.Scheme == "http") {
		return fmt.Errorf("uknown http scheme, '%s'", config.Scheme)
	}

	// Host validation
	if config.Host == "" {
		return errors.New("host may not be empty")
	}

	// Port validation
	if config.Port != nil {
		if *config.Port < 1 || *config.Port > 65565 {
			return fmt.Errorf("invalid port number, %d", *config.Port)
		}
	}
	return nil
}

func (config HTTPConfig) ResolvePort() (port int, err error) {
	if config.Port != nil {
		port = *config.Port

	} else {
		switch config.Scheme {
		case "http":
			port = 80
		case "https":
			port = 443
		default:
			err = errors.New("could not identify port")
		}
	}
	return
}

func (config *HTTPConfig) URL() (url string, err error) {
	var port int
	if port, err = config.ResolvePort(); err != nil {
		return
	}
	url = fmt.Sprintf("%s://%s:%d", config.Scheme, config.Host, port)
	return
}
