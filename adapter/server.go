package adapter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Kaese72/huemie-lib/logging"
	"github.com/spf13/viper"
)

type Config struct {
	HttpServer struct {
		Port    int    `mapstructure:"port" doc:"the port to listen on for HTTP requests"`
		Address string `mapstructure:"address" default:"0.0.0.0" doc:"the address to listen on for HTTP requests"`
	} `mapstructure:"http-server"`
	Enroll struct {
		// Store is the URL of the device-store this adapter will send updates to
		Store string `mapstructure:"store"`
		// Token is the JWT token provided by the adapter-attendant that identifies this adapter
		Token string `mapstructure:"token"`
	} `mapstructure:"enroll"`
	DebugLogging bool `mapstructure:"debug-logging" default:"false" doc:"if true, the adapter will log debug information"`
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
	myVip.BindEnv("http-server.address")
	myVip.SetDefault("http-server.address", "0.0.0.0")
	// Listen port
	myVip.BindEnv("http-server.port")
	myVip.SetDefault("http-server.port", 8080)
	// # Enroll Config
	// Device store to send updates to
	myVip.BindEnv("enroll.store")
	// Token used to authenticate towards the device store
	myVip.BindEnv("enroll.token")

	// # Logging
	myVip.BindEnv("debug-logging")
	myVip.SetDefault("debug-logging", false)

	var conf Config
	err := myVip.Unmarshal(&conf)
	if err != nil {
		logging.Error(err.Error())
		return Config{}, err
	}
	//Logger assumed initiated
	logging.SetDebugLogging(conf.DebugLogging)
	logging.Debug("Debug logging enabled")
	if conf.Enroll.Store != "" && conf.Enroll.Token == "" {
		err := fmt.Errorf("enroll.token is required when enroll.store is set")
		logging.Error(err.Error())
		return Config{}, err
	}
	// We allow disabling the updates to the device store by not setting the enroll.store,
	// but if it is set we require a token to be set as well to avoid misconfiguration
	if conf.HttpServer.Port <= 0 || conf.HttpServer.Port > 65535 {
		err := fmt.Errorf("http-server.port must be a valid port number")
		logging.Error(err.Error())
		return Config{}, err
	}
	if conf.HttpServer.Address == "" {
		err := fmt.Errorf("http-server.address is required")
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
	go deviceUpdater(conf.Enroll.Store, conf.Enroll.Token, updates)
	logging.Info("Starting HTTP server", map[string]any{"address": conf.HttpServer.Address, "port": conf.HttpServer.Port})
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.HttpServer.Address, conf.HttpServer.Port), router); err != nil {
		return err
	}
	return nil
}
