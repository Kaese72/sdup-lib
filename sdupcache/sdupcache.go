package sdupcache

import (
	"github.com/Kaese72/sdup-lib/sdupcache/filters"
	"github.com/Kaese72/sdup-lib/sduptemplates"
)

type SDUPFilterable interface {
	//Device
	Device(sduptemplates.DeviceID) (sduptemplates.DeviceSpec, error)
	//FIXME searchable
	FilteredDevices(filters.AttributeFilters) ([]sduptemplates.DeviceSpec, error)
	//Attributes
	//DeviceAttributes(sduptemplates.DeviceID) (sduptemplates.AttributeSpecMap, error)
	//DeviceAttribute(sduptemplates.DeviceID, sduptemplates.AttributeKey) (sduptemplates.AttributeSpec, error)

	//Capabilities
	//DeviceCapabilities(sduptemplates.DeviceID) (sduptemplates.CapabilitySpecMap, error)
	//DeviceCapability(sduptemplates.DeviceID, sduptemplates.CapabilityKey) (sduptemplates.CapabilitySpec, error)
}

type SDUPCache interface {
	sduptemplates.SDUPTarget
	SDUPFilterable
}
