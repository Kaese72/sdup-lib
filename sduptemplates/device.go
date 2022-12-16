package sduptemplates

import "github.com/Kaese72/sdup-lib/devicestoretemplates"

// DeviceSpec defines how Device specifications are communicated over SDUP
type DeviceSpec struct {
	ID           string                                                                    `json:"identifier"`
	Attributes   map[devicestoretemplates.AttributeKey]devicestoretemplates.AttributeState `json:"attributes,omitempty"`
	Capabilities map[devicestoretemplates.CapabilityKey]devicestoretemplates.Capability    `json:"capabilities,omitempty"`
}

func (spec DeviceSpec) SpecToInitialUpdate() DeviceUpdate {
	return DeviceUpdate{
		ID:             spec.ID,
		AttributesDiff: spec.Attributes,
		CapabilityDiff: spec.Capabilities,
	}
}

// FIXME This function would be a good place to identify exatly what updates are relevant
func (spec DeviceSpec) ApplyUpdate(update DeviceUpdate) (DeviceSpec, DeviceUpdate) {
	relevantUpdate := DeviceUpdate{
		ID:             update.ID,
		AttributesDiff: map[devicestoretemplates.AttributeKey]devicestoretemplates.AttributeState{},
		CapabilityDiff: map[devicestoretemplates.CapabilityKey]devicestoretemplates.Capability{},
	}
	for attrKey, updateAttrValue := range update.AttributesDiff {
		if specAttr, ok := spec.Attributes[attrKey]; ok {
			if !specAttr.Equal(updateAttrValue) {
				relevantUpdate.AttributesDiff[attrKey] = updateAttrValue
			}

		} else {
			relevantUpdate.AttributesDiff[attrKey] = updateAttrValue
		}
		spec.Attributes[attrKey] = updateAttrValue
	}

	for capKey, capValue := range update.CapabilityDiff {
		if _, ok := spec.Capabilities[capKey]; !ok {
			relevantUpdate.CapabilityDiff[capKey] = capValue

		}
		spec.Capabilities[capKey] = capValue
	}
	return spec, relevantUpdate
}
