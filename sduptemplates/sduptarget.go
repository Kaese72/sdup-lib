package sduptemplates

//SDUPTarget defines the interface required for SDUP to function properly against a target
type SDUPTarget interface {
	Initialize() ([]DeviceSpec, chan DeviceUpdate, error)
	Devices() ([]DeviceSpec, error)
	TriggerCapability(DeviceID, CapabilityKey, CapabilityArgument) error
}
