package adapter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Kaese72/huemie-lib/logging"
	"github.com/spf13/viper"
)

type Config struct {
	Huemie struct {
		Server struct {
			Http struct {
				Port    int    `mapstructure:"port" doc:"the port to listen on for HTTP requests"`
				Address string `mapstructure:"address" default:"0.0.0.0" doc:"the address to listen on for HTTP requests"`
			} `mapstructure:"http"`
		} `mapstructure:"server"`
		Enroll struct {
			// Store is the URL of the device-store this adapter will send updates to
			Store string `mapstructure:"store"`
			// Token is the JWT token provided by the adapter-attendant that identifies this adapter
			Token string `mapstructure:"token"`
		} `mapstructure:"enroll"`
		DebugLogging bool `mapstructure:"debug-logging" default:"false" doc:"if true, the adapter will log debug information"`
	} `mapstructure:"huemie"`
}

type InitializableUpdater interface {
	Initialize() (chan Update, error)
}

func readConfig() (Config, error) {
	myVip := viper.New()
	// We have elected to no use AutomaticEnv() because of https://github.com/spf13/viper/issues/584
	// myVip.AutomaticEnv()
	// Set replaces to allow keys like "database.mongodb.connection-string"
	// WARNING. Overriding any of these may hav unintended consequences.
	myVip.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	// # API configuration
	// Listen address
	myVip.BindEnv("huemie.server.http.address")
	myVip.SetDefault("huemie.server.http.address", "0.0.0.0")
	// Listen port
	myVip.BindEnv("huemie.server.http.port")
	myVip.SetDefault("huemie.server.http.port", 8080)
	// # Enroll Config
	// Device store to send updates to
	myVip.BindEnv("huemie.enroll.store")
	// Token used to authenticate towards the device store
	myVip.BindEnv("huemie.enroll.token")

	// # Logging
	myVip.BindEnv("huemie.debug-logging")
	myVip.SetDefault("huemie.debug-logging", false)

	var conf Config
	err := myVip.Unmarshal(&conf)
	if err != nil {
		logging.Error(err.Error())
		return Config{}, err
	}
	//Logger assumed initiated
	logging.SetDebugLogging(conf.Huemie.DebugLogging)
	logging.Debug("Debug logging enabled")
	if conf.Huemie.Enroll.Store != "" && conf.Huemie.Enroll.Token == "" {
		err := fmt.Errorf("huemie.enroll.token is required when huemie.enroll.store is set")
		logging.Error(err.Error())
		return Config{}, err
	}
	// We allow disabling the updates to the device store by not setting the enroll.store,
	// but if it is set we require a token to be set as well to avoid misconfiguration
	if conf.Huemie.Server.Http.Port <= 0 || conf.Huemie.Server.Http.Port > 65535 {
		err := fmt.Errorf("huemie.server.http.port must be a valid port number")
		logging.Error(err.Error())
		return Config{}, err
	}
	if conf.Huemie.Server.Http.Address == "" {
		err := fmt.Errorf("huemie.server.http.address is required")
		logging.Error(err.Error())
		return Config{}, err
	}
	return conf, nil
}

// StartAdapter initiates and starts the adapter by setting up the HTTP server and the device update loop.
// target, in addition to being a InitializableUpdater, can also implement the interfaces in interface.go which allows for different functionality
// to be enabled on the adapter.
func StartAdapter(target InitializableUpdater) error {
	conf, err := readConfig()
	if err != nil {
		logging.Error("Could not get config", map[string]any{"error": err.Error()})
		return err
	}
	router := createAdapterMux(target)
	logging.Info("Starting device store updater")
	updates, err := target.Initialize()
	if err != nil {
		return err
	}
	go deviceUpdater(conf.Huemie.Enroll.Store, conf.Huemie.Enroll.Token, updates)
	logging.Info("Starting HTTP server", map[string]any{"address": conf.Huemie.Server.Http.Address, "port": conf.Huemie.Server.Http.Port})
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Huemie.Server.Http.Address, conf.Huemie.Server.Http.Port), router); err != nil {
		return err
	}
	return nil
}
