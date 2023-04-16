package sdupcache

import (
	"github.com/Kaese72/sdup-lib/sdupcache/filters"
	"github.com/Kaese72/sdup-lib/sduptemplates"
)

type SDUPFilterable interface {
	//Device
	Device(string) (sduptemplates.DeviceSpec, error)
	//FIXME searchable
	FilteredDevices(filters.AttributeFilters) ([]sduptemplates.DeviceSpec, error)
	//Attributes
	//DeviceAttributes(sduptemplates.DeviceID) (sduptemplates.AttributeSpecMap, error)
	//DeviceAttribute(sduptemplates.DeviceID, devicestoretemplates.AttributeKey) (sduptemplates.AttributeSpec, error)

	//Capabilities
	//DeviceCapabilities(sduptemplates.DeviceID) (sduptemplates.CapabilitySpecMap, error)
	//DeviceCapability(sduptemplates.DeviceID, devicestoretemplates.CapabilityKey) (sduptemplates.CapabilitySpec, error)
}

type SDUPCache interface {
	sduptemplates.SDUPTarget
	SDUPFilterable
}
