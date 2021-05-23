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
//TODO Allow multiple targets
//	TODO if multiple, prefix with target bridge ID
//TODO Consolidate different IDs based on

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
func (client SDUPHTTPClient) Initialize() ([]sduptemplates.DeviceSpec, chan sduptemplates.DeviceUpdate, error) {
	//Start SSE connection to other SDUP service
	eventChan := make(chan *sse.Event, 10)
	dUpdateChan := make(chan sduptemplates.DeviceUpdate, 10)
	go sse.Notify(fmt.Sprintf("%s/subscribe", client.sdupURI), eventChan)
	go func() {
		var dUpdate sduptemplates.DeviceUpdate
		for event := range eventChan {
			err := json.NewDecoder(event.Data).Decode(&dUpdate)
			if err != nil {
				log.Error("JSON hit the fan")
			}
			dUpdateChan <- dUpdate
		}
	}()
	devices, err := client.Devices()
	return devices, dUpdateChan, err
}

//Devices runs a discovery against the SDUP server
func (client SDUPHTTPClient) Devices() ([]sduptemplates.DeviceSpec, error) {
	//GET against /discovery endpoint on SDUP service
	resp, err := http.Get(fmt.Sprintf("%s/discovery", client.sdupURI))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Received status code: %d", resp.StatusCode)
	}
	var templates []sduptemplates.DeviceSpec
	json.NewDecoder(resp.Body).Decode(&templates)
	return templates, nil
}

//TriggerCapability triggers a device capability on the SDUP server
func (client SDUPHTTPClient) TriggerCapability(id sduptemplates.DeviceID, cap sduptemplates.CapabilityKey, arg sduptemplates.CapabilityArgument) error {
	//FIXME Post to the "/capability/{deviceID}/{capabilityKey}" endpoint
	jsonVal, err := json.Marshal(arg)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/capability/%s/%s", client.sdupURI, id, cap), "application/json", bytes.NewBuffer(jsonVal))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
