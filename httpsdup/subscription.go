package httpsdup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Kaese72/sdup-lib/sduptemplates"
)

type Subscription struct {
	Reader chan sduptemplates.DeviceUpdate
}

func NewSubscription() *Subscription {
	return &Subscription{
		Reader: make(chan sduptemplates.DeviceUpdate),
	}
}

type Subscriptions struct {
	subscriptions []*Subscription
	subsMutex     sync.Mutex
	cancelChan    chan *Subscription
	EventChan     chan sduptemplates.DeviceUpdate
}

func NewSubscriptions(updates chan sduptemplates.DeviceUpdate) *Subscriptions {
	subs := &Subscriptions{
		subscriptions: []*Subscription{},
		cancelChan:    make(chan *Subscription),
		EventChan:     updates,
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

func (subscriptions *Subscriptions) eventRoutine() {
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

		case event := <-subscriptions.EventChan:
			func() {
				subscriptions.subsMutex.Lock()
				defer subscriptions.subsMutex.Unlock()
				for subscriptionIndex := range subscriptions.subscriptions {
					// FIXME currently blocking
					subscriptions.subscriptions[subscriptionIndex].Reader <- event
				}
			}()
		}
	}
}

//Subscribe returns a channel that sends updates
func (subscriptions *Subscriptions) subscribe() *Subscription {
	//todo go routine that reads and forwards events
	subscriptions.subsMutex.Lock()
	defer subscriptions.subsMutex.Unlock()
	newSub := NewSubscription()
	subscriptions.subscriptions = append(subscriptions.subscriptions, newSub)
	return newSub
}

//UnSubscribe cancels a subscription and closes the associated chan
func (subscriptions *Subscriptions) unSubscribe(subCancel *Subscription) {
	go func() {
		subscriptions.cancelChan <- subCancel
	}()
}

//Subscribe connects the client to an SSE channel that informs the client of device changes
func (subscriptions *Subscriptions) Subscribe(writer http.ResponseWriter, reader *http.Request) {
	//log.Log(log.Info, "Started SSE handler", nil)
	// prepare the header
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, _ := writer.(http.Flusher)

	subscription := subscriptions.subscribe()
	doneChan := reader.Context().Done()
	for {

		select {
		// connection is closed then defer will be executed
		case <-doneChan:
			// Communicate the cancellation of this subscription
			subscriptions.cancelChan <- subscription
			doneChan = nil

		case event, ok := <-subscription.Reader:
			if ok {
				jsonString, err := json.Marshal(event)
				if err != nil {
					//log.Log(log.Error, "Failed to Marshal device update", nil)

				} else {
					fmt.Fprintf(writer, "data: %s\n\n", jsonString)
					flusher.Flush()
				}

			} else {
				return
			}
		}
	}
}
