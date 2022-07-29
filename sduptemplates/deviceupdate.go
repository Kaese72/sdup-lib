package sduptemplates

//DeviceUpdate defines how changes in a device are communicated over SDUP
type DeviceUpdate struct {
	// Device specific fields
	ID             DeviceID          `json:"deviceid"`
	AttributesDiff AttributeStateMap `json:"attributes,omitempty"`
	CapabilityDiff CapabilitySpecMap `json:"capabilities,omitempty"`
}

func (update DeviceUpdate) UpdateToDevice() DeviceSpec {
	return DeviceSpec{
		ID:           update.ID,
		Attributes:   update.AttributesDiff,
		Capabilities: update.CapabilityDiff,
	}
}

func (update DeviceUpdate) Relevant() bool {
	return len(update.AttributesDiff) > 0 || len(update.CapabilityDiff) > 0
}
