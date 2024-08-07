package capabilitytriggerer

import (
	"encoding/json"
	"io"
	"net/http"

	devicestoretemplates "github.com/Kaese72/device-store/rest/models"
	"github.com/Kaese72/sdup-lib/sduptemplates"
	"github.com/gorilla/mux"
)

type CapabilityTriggerer interface {
	TriggerCapability(string, string, devicestoretemplates.DeviceCapabilityArgs) error
	GTriggerCapability(string, string, devicestoretemplates.GroupCapabilityArgs) error
}

// HTTPStatusCode crudely translates error into http status code
func HTTPStatusCode(err error) int {
	switch err {
	case sduptemplates.NoSuchAttribute, sduptemplates.NoSuchCapability, sduptemplates.NoSuchDevice:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// InitHTTPMux initializes a HTTP server mux with the appropriate paths
func InitCapabilityTriggerMux(target CapabilityTriggerer) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/devices/{deviceID}/capabilities/{capabilityKey}", func(writer http.ResponseWriter, reader *http.Request) {
		vars := mux.Vars(reader)
		deviceID := vars["deviceID"]
		capabilityKey := vars["capabilityKey"]
		//log.Log(log.Info, "Triggering capability", map[string]string{"device": deviceID, "capability": capabilityKey})
		var args devicestoretemplates.DeviceCapabilityArgs
		err := json.NewDecoder(reader.Body).Decode(&args)
		if err != nil {
			if err == io.EOF {
				//No body was sent. That is fine
				args = devicestoretemplates.DeviceCapabilityArgs{}
			} else {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		err = target.TriggerCapability(deviceID, capabilityKey, args)
		if err != nil {
			http.Error(writer, err.Error(), HTTPStatusCode(err))
			return

		}
		http.Error(writer, "OK", http.StatusOK)
	}).Methods("POST")

	router.HandleFunc("/groups/{groupID}/capabilities/{capabilityKey}", func(writer http.ResponseWriter, reader *http.Request) {
		vars := mux.Vars(reader)
		groupID := vars["groupID"]
		capabilityKey := vars["capabilityKey"]
		//log.Log(log.Info, "Triggering capability", map[string]string{"group": groupID, "capability": capabilityKey})
		var args devicestoretemplates.GroupCapabilityArgs
		err := json.NewDecoder(reader.Body).Decode(&args)
		if err != nil {
			if err == io.EOF {
				//No body was sent. That is fine
				args = devicestoretemplates.GroupCapabilityArgs{}
			} else {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		err = target.GTriggerCapability(groupID, capabilityKey, args)
		if err != nil {
			http.Error(writer, err.Error(), HTTPStatusCode(err))
			return

		}
		http.Error(writer, "OK", http.StatusOK)
	}).Methods("POST")

	return router
}
