package httpsdup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Kaese72/sdup-lib/devicestoretemplates"
	"github.com/Kaese72/sdup-lib/sduptemplates"
	"github.com/Kaese72/sdup-lib/subscription"
	"github.com/gorilla/mux"
)

// InitHTTPMux initializes a HTTP server mux with the appropriate paths
func InitHTTPMux(target sduptemplates.SDUPTarget) (*mux.Router, subscription.Subscriptions) {
	channel, err := target.Initialize()
	if err != nil {
		//FIXME No reason to panic
		panic(err)
	}
	subs := subscription.NewSubscriptions(channel)
	router := mux.NewRouter()
	router.HandleFunc("/devices", func(writer http.ResponseWriter, reader *http.Request) {
		devices, err := target.Devices()
		if err != nil {
			//log.Log(log.Error, err.Error(), nil)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonEncoded, err := json.MarshalIndent(devices, "", "   ")
		if err != nil {
			//log.Log(log.Error, err.Error(), nil)
			http.Error(writer, "Failed to JSON encode SDUPDevices", http.StatusInternalServerError)
		}
		writer.Write(jsonEncoded)
	})

	router.HandleFunc("/groups", func(writer http.ResponseWriter, reader *http.Request) {
		groups, err := target.Groups()
		if err != nil {
			//log.Log(log.Error, err.Error(), nil)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonEncoded, err := json.MarshalIndent(groups, "", "   ")
		if err != nil {
			//log.Log(log.Error, err.Error(), nil)
			http.Error(writer, "Failed to JSON encode SDUP groups", http.StatusInternalServerError)
		}
		writer.Write(jsonEncoded)
	})

	router.HandleFunc("/subscribe", func(writer http.ResponseWriter, reader *http.Request) {
		//log.Log(log.Info, "Started SSE handler", nil)
		// prepare the header
		writer.Header().Set("Content-Type", "text/event-stream")
		writer.Header().Set("Cache-Control", "no-cache")
		writer.Header().Set("Connection", "keep-alive")
		writer.Header().Set("Access-Control-Allow-Origin", "*")

		flusher, _ := writer.(http.Flusher)

		subscription := subs.Subscribe()
		defer subs.UnSubscribe(subscription)

		devices, err := target.Devices()
		//FIXME Small race condition where we may get updates to a state that have already been sent by the initializer. However, we are pretty much guaranteed to end up in the correct state
		if err != nil {
			return
		}
		for _, device := range devices {
			sendEvent(writer, flusher, sduptemplates.UpdateFromDeviceUpdate(device.SpecToInitialUpdate()))
		}

		doneChan := reader.Context().Done()
		for {

			select {
			// connection is closed then defer will be executed
			case <-doneChan:
				// Communicate the cancellation of this subscription
				doneChan = nil

			case event, ok := <-subscription.Updates():
				if ok {
					sendEvent(writer, flusher, event)
				} else {
					return
				}
			}
		}
	})

	router.HandleFunc("/devices/{deviceID}/capabilities/{capabilityKey}", func(writer http.ResponseWriter, reader *http.Request) {
		vars := mux.Vars(reader)
		deviceID := vars["deviceID"]
		capabilityKey := vars["capabilityKey"]
		//log.Log(log.Info, "Triggering capability", map[string]string{"device": deviceID, "capability": capabilityKey})
		var args devicestoretemplates.CapabilityArgs

		err := json.NewDecoder(reader.Body).Decode(&args)
		if err != nil {
			if err == io.EOF {
				//No body was sent. That is fine
				args = devicestoretemplates.CapabilityArgs{}
			} else {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		err = target.TriggerCapability(deviceID, devicestoretemplates.CapabilityKey(capabilityKey), args)
		if err != nil {
			http.Error(writer, err.Error(), HTTPStatusCode(err))
			return

		}
		http.Error(writer, "OK", http.StatusOK)

	}).Methods("POST")

	router.HandleFunc("/groups/{groupId}/capabilities/{capabilityKey}", func(writer http.ResponseWriter, reader *http.Request) {
		vars := mux.Vars(reader)
		groupId := vars["groupId"]
		capabilityKey := vars["capabilityKey"]
		//log.Log(log.Info, "Triggering capability", map[string]string{"device": deviceID, "capability": capabilityKey})
		var args devicestoretemplates.CapabilityArgs
		if err := json.NewDecoder(reader.Body).Decode(&args); err != nil {
			if err == io.EOF {
				// EOF was reached. Let validators later handle that
				args = devicestoretemplates.CapabilityArgs{}
			} else {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		err = target.GTriggerCapability(sduptemplates.DeviceGroupID(groupId), devicestoretemplates.CapabilityKey(capabilityKey), args)
		if err != nil {
			http.Error(writer, err.Error(), HTTPStatusCode(err))
			return
		}
		http.Error(writer, "OK", http.StatusOK)

	}).Methods("POST")

	router.HandleFunc("/healthcheck", func(writer http.ResponseWriter, reader *http.Request) {
		jsonEncoded, err := json.MarshalIndent(devicestoretemplates.HealthCheck{Ok: true}, "", "   ")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Write(jsonEncoded)
	})
	return router, subs
}

func sendEvent(writer http.ResponseWriter, flusher http.Flusher, event sduptemplates.Update) {
	jsonString, err := json.Marshal(event)
	if err != nil {
		//log.Log(log.Error, "Failed to Marshal device update", nil)

	} else {
		fmt.Fprintf(writer, "data: %s\n\n", jsonString)
		flusher.Flush()
	}

}
