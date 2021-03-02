package sduptemplates

//AttributeSpecMap defines the relationship between AttributeKeys and AttributeSpecs
type AttributeSpecMap map[AttributeKey]AttributeSpec

//AttributeSpec defines how an attribute specification is communicated over SDUP
type AttributeSpec struct {
	AttributeState
}
