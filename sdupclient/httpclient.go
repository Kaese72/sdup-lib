package sdupclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Kaese72/sdup-lib/logging"
	"github.com/Kaese72/sdup-lib/sdupclient/config"
	"github.com/Kaese72/sdup-lib/sduptemplates"
	"github.com/Kaese72/sdup-lib/utils/sse"
)

//FIXME Contexts ?

//SDUPHTTPClient connects to another SDUP node
type SDUPHTTPClient struct {
	sdupURI string
}

//NewSDUPHTTPClient instansiates a SDUPHTTPClient
func NewSDUPHTTPClient(config config.HTTPConfig) (sduptemplates.SDUPTarget, error) {
	baseURI, err := config.URL()
	if err != nil {
		return nil, err
	}
	return SDUPHTTPClient{
		sdupURI: baseURI,
	}, nil
}

//DeviceUpdates starts fetching device updates from the SDUP server
func (client SDUPHTTPClient) Initialize() (chan sduptemplates.Update, error) {
	//Start SSE connection to other SDUP service
	eventChan := make(chan *sse.Event, 10)
	dUpdateChan := make(chan sduptemplates.Update, 10)
	go sse.NotifyReconnect(fmt.Sprintf("%s/subscribe", client.sdupURI), eventChan)
	go func() {
		for event := range eventChan {
			var dUpdate sduptemplates.Update
			err := json.NewDecoder(event.Data).Decode(&dUpdate)
			if err != nil {
				log.Error("JSON hit the fan")
			}
			if strUpdate, err := json.Marshal(dUpdate); err == nil {
				log.Info(fmt.Sprintf("Client recieved update: %s", strUpdate))
			} else {
				log.Info("Failed to marshal update")
			}
			dUpdateChan <- dUpdate
		}
	}()
	return dUpdateChan, nil
}

//Devices runs a discovery against the SDUP server
func (client SDUPHTTPClient) Devices() ([]sduptemplates.DeviceSpec, error) {
	//GET against /discovery endpoint on SDUP service
	resp, err := http.Get(fmt.Sprintf("%s/devices", client.sdupURI))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received status code: %d", resp.StatusCode)
	}
	var templates []sduptemplates.DeviceSpec
	json.NewDecoder(resp.Body).Decode(&templates)
	return templates, nil
}

//Devices runs a discovery against the SDUP server
func (client SDUPHTTPClient) Groups() ([]sduptemplates.DeviceGroupSpec, error) {
	//GET against /discovery endpoint on SDUP service
	resp, err := http.Get(fmt.Sprintf("%s/groups", client.sdupURI))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received status code: %d", resp.StatusCode)
	}
	var templates []sduptemplates.DeviceGroupSpec
	json.NewDecoder(resp.Body).Decode(&templates)
	return templates, nil
}

//TriggerCapability triggers a device capability on the SDUP server
func (client SDUPHTTPClient) TriggerCapability(id sduptemplates.DeviceID, cap sduptemplates.CapabilityKey, arg sduptemplates.CapabilityArgument) error {
	jsonVal, err := json.Marshal(arg)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/devices/%s/capabilities/%s", client.sdupURI, id, cap), "application/json", bytes.NewBuffer(jsonVal))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client SDUPHTTPClient) GTriggerCapability(id sduptemplates.DeviceGroupID, cap sduptemplates.CapabilityKey, arg sduptemplates.CapabilityArgument) error {
	//FIXME Post to the "/capability/{deviceID}/{capabilityKey}" endpoint
	jsonVal, err := json.Marshal(arg)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/groups/%s/capabilities/%s", client.sdupURI, id, cap), "application/json", bytes.NewBuffer(jsonVal))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
