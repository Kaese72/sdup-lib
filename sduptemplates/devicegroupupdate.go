package sduptemplates

import devicestoretemplates "github.com/Kaese72/device-store/rest/models"

type DeviceGroupUpdate struct {
	// Device Group specific fields
	GroupID DeviceGroupID `json:"groupid"`
	Name    string        `json:"groupname"`
	// DeviceIDs      []string                                                               `json:"deviceids"`
	CapabilityDiff map[devicestoretemplates.CapabilityKey]devicestoretemplates.Capability `json:"groupcapabilities,omitempty"`
}

func (update DeviceGroupUpdate) UpdateToDeviceGroup() DeviceGroupSpec {
	return DeviceGroupSpec{
		ID:   update.GroupID,
		Name: update.Name,
		// DeviceIDs:    update.DeviceIDs,
		Capabilities: update.CapabilityDiff,
	}
}

func (update DeviceGroupUpdate) DeviceStorePatch() devicestoretemplates.Group {
	return devicestoretemplates.Group{
		Identifier:   string(update.GroupID),
		Capabilities: update.CapabilityDiff,
		Name:         update.Name,
		// DeviceIds: update.DeviceIDs // FIXME this is relevant
	}
}
