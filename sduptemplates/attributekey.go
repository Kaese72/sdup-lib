package sduptemplates

//AttributeKey is the string identifier of an attribute
type AttributeKey string

const (
	//AttributeActive represents whether the device is currently on or off
	AttributeActive AttributeKey = "active"
	//AttributeColorXY represents the primary color of the device, represented by xy coordinates
	AttributeColorXY AttributeKey = "colorxy"
	//AttributeColorTemp represents the primary color of the device, represented by xy coordinates
	AttributeColorTemp AttributeKey = "colorct"
	//AttributeDescription is a readable description, presentable to a user. Should not be used to identify the device
	AttributeDescription AttributeKey = "description"
	//AttributeUniqueID globally identifes a device across bridges. eg. MAC addresses
	AttributeUniqueID AttributeKey = "uniqueID"
)
