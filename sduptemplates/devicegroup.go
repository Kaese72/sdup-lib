package sduptemplates

import (
	"reflect"
)

type DeviceGroupID string

type DeviceGroupSpec struct {
	ID        DeviceGroupID `json:"id"`
	Name      string        `json:"name"`
	DeviceIDs []DeviceID    `json:"deviceids"`
}

func (spec DeviceGroupSpec) SpecToInitialUpdate() DeviceGroupUpdate {
	return DeviceGroupUpdate{
		GroupID:   spec.ID,
		Name:      spec.Name,
		DeviceIDs: spec.DeviceIDs,
	}
}

func (spec DeviceGroupSpec) Equal(other DeviceGroupSpec) bool {
	return reflect.DeepEqual(spec.DeviceIDs, other.DeviceIDs)
}
