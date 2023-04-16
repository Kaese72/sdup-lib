package devicestoretemplates

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type BridgeKey string

type Bridge struct {
	Identifier BridgeKey `json:"identifier"`
	URI        string    `json:"uri"`
}

type HealthCheck struct {
	Ok bool `json:"ok"`
}

func (bridge Bridge) HealthCheck() error {
	resp, err := http.Get(fmt.Sprintf("%s/healthcheck", bridge.URI))
	if err != nil {
		return err
	}
	healthCheckRest := HealthCheck{}
	err = json.NewDecoder(resp.Body).Decode(&healthCheckRest)
	if err != nil {
		return err
	}
	if !healthCheckRest.Ok {
		return errors.New("health check failed")
	}
	return nil
}
