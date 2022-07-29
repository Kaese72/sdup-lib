package devicestoretemplates

type AttributeState struct {
	Boolean *bool    `json:"boolean-state"`
	Numeric *float32 `json:"numeric-state"`
	Text    *string  `json:"string-state"`
}
