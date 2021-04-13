package subscription

import (
	"sync"

	"github.com/Kaese72/sdup-lib/sduptemplates"
)

//Subscriptions is simply a container of subscriptions
type Subscriptions interface {
	Subscribe() Subscription
	UnSubscribe(Subscription)
}

type subsImpl struct {
	subscriptions []Subscription
	subsMutex     sync.Mutex
	cancelChan    chan Subscription
	eventChan     chan sduptemplates.DeviceUpdate
}

//NewSubscriptions creates a Subscriptions container
func NewSubscriptions(updates chan sduptemplates.DeviceUpdate) Subscriptions {
	subs := &subsImpl{
		subscriptions: []Subscription{},
		cancelChan:    make(chan Subscription),
		eventChan:     updates,
	}
	go subs.eventRoutine()
	return subs
}

//SliceIndex Searched for the index of the object matching the `predicate` function,
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func (subscriptions *subsImpl) eventRoutine() {
	for {
		select {
		case subscription := <-subscriptions.cancelChan:
			index := SliceIndex(len(subscriptions.subscriptions), func(i int) bool { return subscriptions.subscriptions[i] == subscription })
			if index < 0 {
				panic("Could not find subscription")
			}
			func() {
				subscriptions.subsMutex.Lock()
				defer subscriptions.subsMutex.Unlock()
				subscriptions.subscriptions[index] = subscriptions.subscriptions[len(subscriptions.subscriptions)-1]
				subscriptions.subscriptions = subscriptions.subscriptions[:len(subscriptions.subscriptions)-1]
			}()

		case event := <-subscriptions.eventChan:
			func() {
				subscriptions.subsMutex.Lock()
				defer subscriptions.subsMutex.Unlock()
				for subscriptionIndex := range subscriptions.subscriptions {
					// FIXME currently blocking
					subscriptions.subscriptions[subscriptionIndex].Updates() <- event
				}
			}()
		}
	}
}

//Subscribe returns a channel that sends updates
func (subscriptions *subsImpl) Subscribe() Subscription {
	//todo go routine that reads and forwards events
	subscriptions.subsMutex.Lock()
	defer subscriptions.subsMutex.Unlock()
	newSub := NewSubscription()
	subscriptions.subscriptions = append(subscriptions.subscriptions, newSub)
	return newSub
}

//UnSubscribe cancels a subscription and closes the associated chan
func (subscriptions *subsImpl) UnSubscribe(subCancel Subscription) {
	go func() {
		subscriptions.cancelChan <- subCancel
	}()
}
