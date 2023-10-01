package sduptemplates

import (
	"encoding/json"
	"errors"
)

type Update struct {
	deviceUpdate      DeviceUpdate
	deviceGroupUpdate DeviceGroupUpdate
}

type updateParser struct {
	DeviceUpdate
	DeviceGroupUpdate
}

func (update *Update) UnmarshalJSON(b []byte) error {
	parsed := updateParser{}
	if err := json.Unmarshal(b, &parsed); err != nil {
		return err
	}
	update.deviceUpdate = parsed.DeviceUpdate
	update.deviceGroupUpdate = parsed.DeviceGroupUpdate
	return nil
}

func (update Update) MarshalJSON() ([]byte, error) {
	return json.Marshal(updateParser{
		DeviceUpdate:      update.deviceUpdate,
		DeviceGroupUpdate: update.deviceGroupUpdate,
	})
}

func UpdateFromDeviceUpdate(update DeviceUpdate) Update {
	return Update{deviceUpdate: update}
}

func UpdateFromDeviceGroupUpdate(update DeviceGroupUpdate) Update {
	return Update{deviceGroupUpdate: update}
}

func (update Update) GetDeviceUpdate() (DeviceUpdate, error) {
	var err error
	if update.deviceUpdate.ID == "" {
		err = errors.New("Update is not a device update")
	}
	return update.deviceUpdate, err
}

func (update Update) GetDeviceGroupUpdate() (DeviceGroupUpdate, error) {
	var err error
	if update.deviceGroupUpdate.GroupID == "" {
		err = errors.New("Update is not a device group update")
	}
	return update.deviceGroupUpdate, err
}
