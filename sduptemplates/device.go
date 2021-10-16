package sduptemplates

//DeviceID represents a locally unique ID for a device
type DeviceID string

//DeviceSpec defines how Device specifications are communicated over SDUP
type DeviceSpec struct {
	ID           DeviceID          `json:"id"`
	Attributes   AttributeSpecMap  `json:"attributes"`
	Capabilities CapabilitySpecMap `json:"capabilities"`
}

func (spec DeviceSpec) SpecToInitialUpdate() DeviceUpdate {
	diff := AttributeStateMap{}
	for key, val := range spec.Attributes {
		diff[key] = val.AttributeState
	}

	return DeviceUpdate{
		ID:             spec.ID,
		AttributesDiff: diff,
	}
}

//DeviceDiff defines how changes in a device are communicated over SDUP
type DeviceUpdate struct {
	ID             DeviceID          `json:"id"`
	AttributesDiff AttributeStateMap `json:"attributes"`
	//Lost indicates that the device is no longer available. This indicates permanence, and does not include temporary disconnects
	Lost bool `json:"lost"`
}
