package sduptemplates

import "github.com/Kaese72/sdup-lib/devicestoretemplates"

// SDUPTarget defines the interface required for SDUP to function properly against a target
type SDUPTarget interface {
	Initialize() (chan Update, error)
	Devices() ([]DeviceSpec, error)
	Groups() ([]DeviceGroupSpec, error)
	TriggerCapability(string, devicestoretemplates.CapabilityKey, devicestoretemplates.CapabilityArgs) error
	GTriggerCapability(DeviceGroupID, devicestoretemplates.CapabilityKey, devicestoretemplates.CapabilityArgs) error
}
