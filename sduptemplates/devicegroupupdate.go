package sduptemplates

type DeviceGroupUpdate struct {
	// Device Group specific fields
	GroupID   DeviceGroupID `json:"groupid"`
	Name      string        `json:"groupname"`
	DeviceIDs []string      `json:"deviceids"`
}

func (update DeviceGroupUpdate) UpdateToDeviceGroup() DeviceGroupSpec {
	return DeviceGroupSpec{
		ID:        update.GroupID,
		Name:      update.Name,
		DeviceIDs: update.DeviceIDs,
	}
}
