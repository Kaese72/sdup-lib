package sduptemplates

import (
	"reflect"

	"github.com/Kaese72/device-store/rest/models"
)

type DeviceGroupID string

type DeviceGroupSpec struct {
	ID   DeviceGroupID `json:"id"`
	Name string        `json:"name"`
	// DeviceIDs    []string                                   `json:"deviceids"`
	Capabilities map[models.CapabilityKey]models.Capability `json:"capabilities"`
}

func (spec DeviceGroupSpec) SpecToInitialUpdate() DeviceGroupUpdate {
	return DeviceGroupUpdate{
		GroupID: spec.ID,
		Name:    spec.Name,
		// DeviceIDs:      spec.DeviceIDs,
		CapabilityDiff: spec.Capabilities,
	}
}

func (spec DeviceGroupSpec) Equal(other DeviceGroupSpec) bool {
	return reflect.DeepEqual(spec, other)
}
