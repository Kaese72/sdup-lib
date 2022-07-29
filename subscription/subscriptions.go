package subscription

import (
	"github.com/Kaese72/sdup-lib/sduptemplates"
)

//Subscriptions is simply a container of subscriptions
type Subscriptions interface {
	Subscribe() Subscription
	UnSubscribe(Subscription)
}

type subsImpl struct {
	subscriptions []Subscription
	cancelChan    chan Subscription
	eventChan     chan sduptemplates.Update
	subscribeChan chan Subscription
}

//NewSubscriptions creates a Subscriptions container
func NewSubscriptions(updates chan sduptemplates.Update) Subscriptions {
	subs := &subsImpl{
		subscriptions: []Subscription{},
		cancelChan:    make(chan Subscription),
		eventChan:     updates,
		subscribeChan: make(chan Subscription),
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
	//FIXME Separate subscriptions from networking, If a subscriber stops receiving for whatever reason, all other subscribers are affected as well
	for {
		select {
		case subscription := <-subscriptions.cancelChan:
			index := SliceIndex(len(subscriptions.subscriptions), func(i int) bool { return subscriptions.subscriptions[i] == subscription })
			if index < 0 {
				panic("Could not find subscription")
			}
			subscriptions.subscriptions[index].Close()
			subscriptions.subscriptions[index] = subscriptions.subscriptions[len(subscriptions.subscriptions)-1]
			subscriptions.subscriptions = subscriptions.subscriptions[:len(subscriptions.subscriptions)-1]

		case event := <-subscriptions.eventChan:
			for subscriptionIndex := range subscriptions.subscriptions {
				// FIXME currently blocking
				subscriptions.subscriptions[subscriptionIndex].Updates() <- event
			}
		case newSubscription := <-subscriptions.subscribeChan:
			// Register new subscription and feed initial state
			subscriptions.subscriptions = append(subscriptions.subscriptions, newSubscription)
		}
	}
}

//Subscribe returns a channel that sends updates
func (subscriptions *subsImpl) Subscribe() Subscription {
	newSub := NewSubscription()
	subscriptions.subscribeChan <- newSub
	return newSub
}

//UnSubscribe cancels a subscription and closes the associated chan
func (subscriptions *subsImpl) UnSubscribe(subCancel Subscription) {
	go func() {
		subscriptions.cancelChan <- subCancel
	}()
}
