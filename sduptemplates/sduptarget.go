package sduptemplates

//SDUPTarget defines the interface required for SDUP to function properly against a target
type SDUPTarget interface {
	Initialize() (chan Update, error)
	Devices() ([]DeviceSpec, error)
	Groups() ([]DeviceGroupSpec, error)
	TriggerCapability(DeviceID, CapabilityKey, CapabilityArgument) error
	GTriggerCapability(DeviceGroupID, CapabilityKey, CapabilityArgument) error
}
