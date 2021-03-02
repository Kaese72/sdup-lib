package sduptemplates

//DeviceID represents a locally unique ID for a device
type DeviceID string

//DeviceSpec defines how Device specifications are communicated over SDUP
type DeviceSpec struct {
	ID           DeviceID          `json:"id"`
	Attributes   AttributeSpecMap  `json:"attributes"`
	Capabilities CapabilitySpecMap `json:"capabilities"`
}

//DeviceDiff defines how changes in a device are communicated over SDUP
type DeviceUpdate struct {
	ID             DeviceID          `json:"id"`
	AttributesDiff AttributeStateMap `json:"attributes"`
}
