package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/huemie-lib/logging"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
)

// AdapterError is a simple struct to represent errors in the adapter with a message and an HTTP status code.
type AdapterError struct {
	Message string `json:"message"`
	// Code maps to HTTP status codes, where 400, 404, and 500 are the allowed ones
	Code int `json:"-"`
}

type AdapterSuccess struct {
	Message string `json:"message"`
}

func (e AdapterError) StatusError() huma.StatusError {
	whitelistedStatusCodes := []int{
		http.StatusBadRequest,
		http.StatusNotFound,
		http.StatusInternalServerError,
	}
	if !slices.Contains(whitelistedStatusCodes, e.Code) {
		e.Code = http.StatusInternalServerError
	}
	return huma.NewError(e.Code, e.Message)
}

type DeviceTriggerCapability interface {
	DeviceTriggerCapability(string, string, ingestmodels.IngestDeviceCapabilityArgs) *AdapterError
}

type GroupTriggerCapability interface {
	GroupTriggerCapability(string, string, ingestmodels.IngestGroupCapabilityArgs) *AdapterError
}

func createAdapterMux(adapter any) *mux.Router {
	router := mux.NewRouter()
	humaConfig := huma.DefaultConfig("sdup-adapter", "1.0.0")
	humaConfig.OpenAPIPath = "/openapi"
	humaConfig.DocsPath = "/docs"
	api := humamux.New(router, humaConfig)
	createAdapterHuma(api, adapter)
	return router
}

func createAdapterHuma(api huma.API, adapter any) {
	abilities := []string{}
	if target, ok := adapter.(DeviceTriggerCapability); ok {
		logging.Info("Adapter IS applicable for DeviceTriggerCapability, creating endpoint")
		abilities = append(abilities, "DeviceTriggerCapability")
		huma.Post(api, "/devices/{deviceID}/capabilities/{capabilityKey}", func(ctx context.Context, input *struct {
			DeviceID      string                                   `path:"deviceID" doc:"the device to trigger the capability for"`
			CapabilityKey string                                   `path:"capabilityKey" doc:"the capability to trigger"`
			Body          *ingestmodels.IngestDeviceCapabilityArgs `body:""`
		}) (*struct {
			Body *AdapterSuccess
		}, error) {
			args := ingestmodels.IngestDeviceCapabilityArgs{}
			if input.Body != nil {
				args = *input.Body
			}
			err := target.DeviceTriggerCapability(input.DeviceID, input.CapabilityKey, args)
			if err != nil {
				return nil, err.StatusError()
			}
			return &struct {
				Body *AdapterSuccess
			}{
				Body: &AdapterSuccess{
					Message: "Capability triggered successfully",
				},
			}, nil
		})
	} else {
		logging.Info("Adapter is NOT applicable for DeviceTriggerCapability, skipping endpoint creation")
	}

	if target, ok := adapter.(GroupTriggerCapability); ok {
		logging.Info("Adapter IS applicable for GroupTriggerCapability, creating endpoint")
		abilities = append(abilities, "GroupTriggerCapability")
		huma.Post(api, "/groups/{groupID}/capabilities/{capabilityKey}", func(ctx context.Context, input *struct {
			GroupID       string                                  `path:"groupID" doc:"the group to trigger the capability for"`
			CapabilityKey string                                  `path:"capabilityKey" doc:"the capability to trigger"`
			Body          *ingestmodels.IngestGroupCapabilityArgs `body:""`
		}) (*struct {
			Body *AdapterSuccess
		}, error) {
			args := ingestmodels.IngestGroupCapabilityArgs{}
			if input.Body != nil {
				args = *input.Body
			}
			err := target.GroupTriggerCapability(input.GroupID, input.CapabilityKey, args)
			if err != nil {
				return nil, err.StatusError()
			}
			return &struct {
				Body *AdapterSuccess
			}{
				Body: &AdapterSuccess{
					Message: "Capability triggered successfully",
				},
			}, nil
		})
	} else {
		logging.Info("Adapter is NOT applicable for GroupTriggerCapability, skipping endpoint creation")
	}
	huma.Get(api, "/abilities", func(ctx context.Context, input *struct{}) (*struct {
		Body *struct {
			Abilities []string `json:"abilities"`
		}
	}, error) {
		return &struct {
			Body *struct {
				Abilities []string `json:"abilities"`
			}
		}{
			Body: &struct {
				Abilities []string `json:"abilities"`
			}{
				Abilities: abilities,
			},
		}, nil
	})
}

type Update struct {
	Device *ingestmodels.IngestDevice
	Group  *ingestmodels.IngestGroup
}

func pushDeviceUpdate(deviceStoreBaseUrl string, token string, device ingestmodels.IngestDevice) error {
	bPayload, err := json.Marshal(device)
	if err != nil {
		logging.Error("Failed to marshal struct to JSON to send to device store", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	logging.Info("Sending blob to device store", map[string]interface{}{"blob": string(bPayload)})
	devicePayload, err := http.NewRequest("POST", fmt.Sprintf("%s/device-ingest/v0/devices", deviceStoreBaseUrl), bytes.NewBuffer(bPayload))
	if err != nil {
		logging.Error("Failed to create request", map[string]interface{}{"error": err.Error()})
		return err
	}
	devicePayload.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
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

func pushGroupUpdate(baseUrl string, token string, group ingestmodels.IngestGroup) error {
	bPayload, err := json.Marshal(group)
	if err != nil {
		logging.Error("Failed to marshal struct to JSON to send to device store", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	logging.Info("Sending blob to device store", map[string]interface{}{"blob": string(bPayload)})
	groupPayload, err := http.NewRequest("POST", fmt.Sprintf("%s/device-ingest/v0/groups", baseUrl), bytes.NewBuffer(bPayload))
	if err != nil {
		logging.Error("Failed to create request", map[string]interface{}{"error": err.Error()})
		return err
	}
	groupPayload.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
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

// deviceUpdater is what takes the updates from the channel and pushes them to the device store
func deviceUpdater(deviceStoreBaseUrl string, token string, updateChan chan Update) {
	for update := range updateChan {
		if deviceStoreBaseUrl == "" {
			logging.Debug("No device store URL configured, skipping update", map[string]any{"update": update})
			continue
		}
		if update.Device != nil {
			if err := pushDeviceUpdate(deviceStoreBaseUrl, token, *update.Device); err != nil {
				logging.Error("Failed to send device update", map[string]any{"error": err.Error()})
			}
		} else if update.Group != nil {
			if err := pushGroupUpdate(deviceStoreBaseUrl, token, *update.Group); err != nil {
				logging.Error("Failed to send group update", map[string]any{"error": err.Error()})
			}
		} else {
			logging.Error("Failed to send device group update", map[string]any{"error": "No device or group in update"})
		}
	}
}
