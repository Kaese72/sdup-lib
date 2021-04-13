package subscription

import "github.com/Kaese72/sdup-lib/sduptemplates"

//Subscription represents one currently listened to subscription
type Subscription interface {
	Updates() chan sduptemplates.DeviceUpdate
}

type subImpl struct {
	updates chan sduptemplates.DeviceUpdate
}

func (sub subImpl) Updates() chan sduptemplates.DeviceUpdate {
	return sub.updates
}

//NewSubscription creates a Subscription with the default implementation
func NewSubscription() Subscription {
	return subImpl{
		updates: make(chan sduptemplates.DeviceUpdate),
	}
}
