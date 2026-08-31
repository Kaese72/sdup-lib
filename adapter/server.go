package adapter

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Kaese72/huemie-lib/logging"
	"github.com/Kaese72/sdup-lib/retry"
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
		// Updates tunes the exponential-backoff retry applied to every device
		// and group update pushed to the device store.
		Updates struct {
			RetryMaxAttempts int `mapstructure:"retry-max-attempts" default:"8" doc:"total attempts per update push, including the first; <=0 means retry until the process stops"`
			RetryBaseDelayMs int `mapstructure:"retry-base-delay-ms" default:"500" doc:"backoff before the second attempt, in milliseconds"`
			RetryMaxDelayMs  int `mapstructure:"retry-max-delay-ms" default:"30000" doc:"cap on any single backoff wait, in milliseconds"`
		} `mapstructure:"updates"`
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

	// # Update push retry (exponential backoff towards the device store)
	myVip.BindEnv("huemie.updates.retry-max-attempts")
	myVip.SetDefault("huemie.updates.retry-max-attempts", 8)
	myVip.BindEnv("huemie.updates.retry-base-delay-ms")
	myVip.SetDefault("huemie.updates.retry-base-delay-ms", 500)
	myVip.BindEnv("huemie.updates.retry-max-delay-ms")
	myVip.SetDefault("huemie.updates.retry-max-delay-ms", 30000)

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
	retryCfg := retry.Config{
		MaxAttempts: conf.Huemie.Updates.RetryMaxAttempts,
		BaseDelay:   time.Duration(conf.Huemie.Updates.RetryBaseDelayMs) * time.Millisecond,
		MaxDelay:    time.Duration(conf.Huemie.Updates.RetryMaxDelayMs) * time.Millisecond,
		Multiplier:  2,
		Jitter:      0.2,
	}
	go deviceUpdater(conf.Huemie.Enroll.Store, conf.Huemie.Enroll.Token, retryCfg, updates)
	logging.Info("Starting HTTP server", map[string]any{"address": conf.Huemie.Server.Http.Address, "port": conf.Huemie.Server.Http.Port})
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Huemie.Server.Http.Address, conf.Huemie.Server.Http.Port), router); err != nil {
		return err
	}
	return nil
}
