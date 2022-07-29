package devicestoretemplates

import "github.com/Kaese72/sdup-lib/sduptemplates"

type Device struct {
	Identifier   string                                        `json:"identifier"`
	Attributes   map[sduptemplates.AttributeKey]AttributeState `json:"attributes"`
	Capabilities map[sduptemplates.CapabilityKey]Capability    `json:"capabilities"`
}
