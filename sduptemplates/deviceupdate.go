package sduptemplates

import devicestoretemplates "github.com/Kaese72/device-store/rest/models"

//DeviceUpdate defines how changes in a device are communicated over SDUP
type DeviceUpdate struct {
	// Device specific fields
	ID             string                                                                    `json:"identifier"`
	AttributesDiff map[devicestoretemplates.AttributeKey]devicestoretemplates.AttributeState `json:"attributes,omitempty"`
	CapabilityDiff map[devicestoretemplates.CapabilityKey]devicestoretemplates.Capability    `json:"capabilities,omitempty"`
}

func (update DeviceUpdate) UpdateToDevice() DeviceSpec {
	return DeviceSpec{
		ID:           update.ID,
		Attributes:   update.AttributesDiff,
		Capabilities: update.CapabilityDiff,
	}
}

func (update DeviceUpdate) DeviceStorePatch() devicestoretemplates.Device {
	return devicestoretemplates.Device{
		Identifier:   string(update.ID),
		Attributes:   update.AttributesDiff,
		Capabilities: update.CapabilityDiff,
	}
}

func (update DeviceUpdate) Relevant() bool {
	return len(update.AttributesDiff) > 0 || len(update.CapabilityDiff) > 0
}
