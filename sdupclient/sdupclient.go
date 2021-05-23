package sdupclient

import (
	"errors"

	"github.com/Kaese72/sdup-lib/sdupclient/config"
	"github.com/Kaese72/sdup-lib/sduptemplates"
)

//NewSDUPClient instansiates a SDUPClient
func NewSDUPClient(config config.Config) (sduptemplates.SDUPTarget, error) {
	//FIXME Start DeviceUpdates Thread
	if config.HTTP != nil {
		return NewSDUPHTTPClient(*config.HTTP)

	} else {
		return nil, errors.New("no client config supplied")
	}
}
