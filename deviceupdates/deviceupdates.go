package deviceupdates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Kaese72/device-store/rest/models"
	"github.com/Kaese72/huemie-lib/logging"
)

type Update struct {
	Device *models.Device
	Group  *models.Group
}

type DeviceUpdater interface {
	Initialize() (chan Update, error)
}

type StoreEnrollmentConfig struct {
	StoreURL   string `mapstructure:"store"`
	AdapterKey string `mapstructure:"adapter-key"`
}

func pushDeviceUpdate(config StoreEnrollmentConfig, device models.Device) error {
	bPayload, err := json.Marshal(device)
	if err != nil {
		logging.Error("Failed to marshal struct to JSON to send to device store", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	logging.Info("Sending blob to device store", map[string]interface{}{"blob": string(bPayload)})
	devicePayload, err := http.NewRequest("POST", fmt.Sprintf("%s/device-store/v0/devices", config.StoreURL), bytes.NewBuffer(bPayload))
	if err != nil {
		logging.Error("Failed to create request", map[string]interface{}{"error": err.Error()})
		return err
	}
	devicePayload.Header.Set("Bridge-Key", config.AdapterKey)
	resp, err := http.DefaultClient.Do(
		devicePayload,
	)
	if err != nil {
		logging.Error("Failed to http.Do request", map[string]interface{}{"error": err.Error()})
		return err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error("Failed to read response body on response", map[string]interface{}{"error": err.Error()})
		return err
	}

	logging.Info("Sent payload to device store", map[string]interface{}{"Response Code": resp.Status, "Response Body": string(respBody)})
	return nil
}

func pushGroupUpdate(config StoreEnrollmentConfig, group models.Group) error {
	bPayload, err := json.Marshal(group)
	if err != nil {
		logging.Error("Failed to marshal struct to JSON to send to device store", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	logging.Info("Sending blob to device store", map[string]interface{}{"blob": string(bPayload)})
	groupPayload, err := http.NewRequest("POST", fmt.Sprintf("%s/device-store/v0/groups", config.StoreURL), bytes.NewBuffer(bPayload))
	if err != nil {
		logging.Error("Failed to create request", map[string]interface{}{"error": err.Error()})
		return err
	}
	groupPayload.Header.Set("Bridge-Key", config.AdapterKey)
	resp, err := http.DefaultClient.Do(
		groupPayload,
	)
	if err != nil {
		logging.Error("Failed to http.Do request", map[string]interface{}{"error": err.Error()})
		return err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error("Failed to read response body on response", map[string]interface{}{"error": err.Error()})
		return err
	}

	logging.Info("Sent payload to device store", map[string]interface{}{"Response Code": resp.Status, "Response Body": string(respBody)})
	return nil
}

func InitDeviceUpdater(config StoreEnrollmentConfig, updater DeviceUpdater) error {
	logging.Info("Starting device store updater")
	updates, err := updater.Initialize()
	if err != nil {
		return err
	}
	go func() {
		for update := range updates {
			if update.Device != nil {
				if err := pushDeviceUpdate(config, *update.Device); err != nil {
					logging.Error("Failed to send device group update", map[string]interface{}{"error": err.Error()})
				}
			} else if update.Group != nil {
				if err := pushGroupUpdate(config, *update.Group); err != nil {
					logging.Error("Failed to send device group update", map[string]interface{}{"error": err.Error()})
				}
			} else {
				logging.Error("Failed to send device group update", map[string]interface{}{"error": "No device or group in update"})
			}
		}
	}()
	return nil
}
