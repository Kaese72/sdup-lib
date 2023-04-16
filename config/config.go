package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Kaese72/huemie-lib/logging"
)

type HuemieStandardConfig interface {
	PopulateExample()
	Validate() error
}

func ReadConfig(conf HuemieStandardConfig) error {
	if _, err := os.Stat("./settings.json"); err == nil {
		file, err := os.Open("./settings.json")
		if err != nil {
			logging.Error(fmt.Sprintf("Unable to open local settings file, %s", err.Error()))
			return err
		}
		if err := json.NewDecoder(file).Decode(&conf); err != nil {
			logging.Error(err.Error())
			return err
		}

	} else {
		if err := json.NewDecoder(os.Stdin).Decode(&conf); err != nil {
			logging.Error(err.Error())
			return err
		}
	}

	return nil
}

func ValidConfigOrPrintExample(conf HuemieStandardConfig) bool {
	if err := conf.Validate(); err != nil {
		logging.Error(err.Error())
		conf.PopulateExample()
		obj, err := json.Marshal(conf)
		if err != nil {
			logging.Error(err.Error())
		}
		_, err = fmt.Fprintf(os.Stdout, "%s\n", obj)
		if err != nil {
			logging.Error(err.Error())
		}
		return false
	}
	return true
}
