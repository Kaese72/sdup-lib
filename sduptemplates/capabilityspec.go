package sduptemplates

//CapabilityMap defines the relationship between CapabilityKeys and CapabilitySpecs
type CapabilitySpecMap map[CapabilityKey]CapabilitySpec

type CapabilityArgument map[string]interface{}

type CapabilityType string

const (
	Preconfigured CapabilityType = "preconfigured"
	KeyVal        CapabilityType = "keyval"
)

//CapabilitySpec defines how capability specifications are communicated over SDUP
type CapabilitySpec struct {
}
