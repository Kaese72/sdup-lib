package sduptemplates

//CapabilityKey is the string identifier of a capability
type CapabilityKey string

const (
	//CapabilityActivate means the associated attribute can be activated
	CapabilityActivate CapabilityKey = "activate"
	//CapabilityDeactivate means the associated attribute can be deactivated
	CapabilityDeactivate CapabilityKey = "deactivate"
	//CapabilitySetColorXY means that you can change the x and y coordinates in color mode
	CapabilitySetColorXY CapabilityKey = "setcolorxy"
	//CapabilitySetColorTemp means that you can change the temperature in color mode
	CapabilitySetColorTemp CapabilityKey = "setcolorct"
)
