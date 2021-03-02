package httpsdup

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Kaese72/sdup-lib/sduptemplates"
	"github.com/gorilla/mux"
)

func InitHTTPMux(target sduptemplates.SDUPTarget) *mux.Router {
	subs := NewSubscriptions(target.DeviceUpdates())
	router := mux.NewRouter()
	router.HandleFunc("/discovery", func(writer http.ResponseWriter, reader *http.Request) {
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
	router.HandleFunc("/subscribe", subs.Subscribe)

	router.HandleFunc("/capability/{deviceID}/{capabilityKey}", func(writer http.ResponseWriter, reader *http.Request) {
		vars := mux.Vars(reader)
		deviceID := vars["deviceID"]
		capabilityKey := vars["capabilityKey"]
		//log.Log(log.Info, "Triggering capability", map[string]string{"device": deviceID, "capability": capabilityKey})
		var args sduptemplates.CapabilityArgument

		err := json.NewDecoder(reader.Body).Decode(&args)
		if err != nil {
			if err == io.EOF {
				//No body was sent. That is fine
				args = sduptemplates.CapabilityArgument{}
			} else {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		err = target.TriggerCapability(sduptemplates.DeviceID(deviceID), sduptemplates.CapabilityKey(capabilityKey), args)
		if err != nil {
			http.Error(writer, err.Error(), HTTPStatusCode(err))
			return

		}
		http.Error(writer, "OK", http.StatusOK)

	}).Methods("POST")
	//router.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir("./ui/"))))
	return router
}
