package sduptemplates

import (
	"reflect"
)

//KeyValContainer contains an unknown set keys with an unknown value type
// Exposes methods that collect and convert appropriately
type KeyValContainer map[string]interface{}

//AttributeStateMap defines the relationship between AttributeKeys and AttributeStates
type AttributeStateMap map[AttributeKey]AttributeState

//AttributeState defines how an attribute state is communicated over SDUP
type AttributeState struct {
	Boolean *bool            `json:"boolean-state,omitempty"`
	Numeric *float32         `json:"numeric-state,omitempty"`
	KeyVal  *KeyValContainer `json:"keyval-state,omitempty"`
	Text    *string          `json:"string-state,omitempty"`
}

//Equivalent checks whether two states are equivalent and should therefore trigger a device update
func (state AttributeState) Equivalent(other AttributeState) bool {
	if state.Boolean != nil && other.Boolean != nil {
		return *state.Boolean == *other.Boolean

	} else if state.KeyVal != nil && other.KeyVal != nil {
		return reflect.DeepEqual(*state.KeyVal, *other.KeyVal)

	} else if state.Numeric != nil && other.Numeric != nil {
		return *state.Numeric == *other.Numeric

	} else if state.Text != nil && other.Text != nil {
		return *state.Text == *other.Text
	}

	//log.Log(log.Error, "Could not find common state", nil)
	//FIXME return error
	return false
}
