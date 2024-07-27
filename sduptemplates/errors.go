package sduptemplates

type SDUPError int

const (
	//NoSuchDevice indicates the device was not found
	NoSuchDevice SDUPError = iota
	//NoSuchAttribute indicates the found device does not have the requested attribute
	NoSuchAttribute
	//NoSuchCapability indicates the found device does not have the requested capability
	NoSuchCapability
	//NoSuchGroup indicates the group was not found
	NoSuchGroup
)

func (err SDUPError) Error() string {
	switch err {
	case NoSuchDevice:
		return "The device was not found"
	case NoSuchAttribute:
		return "The attribute on the selected device was not found"
	case NoSuchCapability:
		return "The capability on the selected device was not found"
	case NoSuchGroup:
		return "The group was not found"
	default:
		return "An unknown error occured"
	}
}
