package sduptemplates

//SDUPTarget defines the interface required for SDUP to function properly against a target
type SDUPTarget interface {
	DeviceUpdates() chan DeviceUpdate
	Devices() ([]DeviceSpec, error)
	TriggerCapability(DeviceID, CapabilityKey, CapabilityArgument) error
}
