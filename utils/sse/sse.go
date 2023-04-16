package sse

// import "astuart.co/go-sse"

//FIXME Credit: https://github.com/andrewstuart/go-sse/blob/master/sse.go
// That repo is missing a valid module. PR once this project is serious

// FIXMEs that should go into the original project
//FIXME context for SSE

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	log "github.com/Kaese72/huemie-lib/logging"
)

// SSE name constants
const (
	eName = "event"
	dName = "data"
)

var (
	//ErrNilChan will be returned by Notify if it is passed a nil channel
	ErrNilChan = fmt.Errorf("nil channel given")
)

// Client is the default client used for requests.
var Client = &http.Client{}

func liveReq(verb, uri string, body io.Reader) (*http.Request, error) {
	req, err := GetReq(verb, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/event-stream")

	return req, nil
}

// Event is a go representation of an http server-sent event
type Event struct {
	URI  string
	Type string
	Data io.Reader
}

// GetReq is a function to return a single request. It will be used by notify to
// get a request and can be replaces if additional configuration is desired on
// the request. The "Accept" header will necessarily be overwritten.
var GetReq = func(verb, uri string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(verb, uri, body)
}

// Notify takes the uri of an SSE stream and channel, and will send an Event
// down the channel when recieved, until the stream is closed. It will then
// close the stream. This is blocking, and so you will likely want to call this
// in a new goroutine (via `go Notify(..)`)
func Notify(uri string, evCh chan<- *Event) error {
	if evCh == nil {
		return ErrNilChan
	}

	req, err := liveReq("GET", uri, nil)
	if err != nil {
		return fmt.Errorf("error getting sse request: %v", err)
	}

	res, err := Client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing request for %s: %v", uri, err)
	}

	log.Info("Connection to SSE endpoint successful", map[string]interface{}{"uri": uri})

	br := bufio.NewReader(res.Body)
	defer res.Body.Close()

	delim := []byte{':', ' '}

	var currEvent *Event

	for {
		bs, err := br.ReadBytes('\n')

		if err != nil {
			if err != io.EOF {
				return err
			}
		}

		if len(bs) <= 1 {
			if err != nil && err == io.EOF {
				break
			} else {
				continue
			}
		}
		//FIXME Split only once
		spl := bytes.Split(bs, delim)

		if len(spl) <= 1 {
			if err != nil && err == io.EOF {
				break
			} else {
				continue
			}
		}

		currEvent = &Event{URI: uri}
		switch string(spl[0]) {
		case eName:
			currEvent.Type = string(bytes.TrimSpace(spl[1]))
		case dName:
			currEvent.Data = bytes.NewBuffer(bytes.TrimSpace(spl[1]))
			evCh <- currEvent
		}
		if err == io.EOF {
			break
		}
	}

	return nil
}

// NotifyReconnect tries to maintain the connection to
func NotifyReconnect(uri string, evCh chan<- *Event) {
	//FIXME Implement cancellation (via contexts ?)
	fallbackCounter := 0
	maxRetries := 5
	baseWaitTime := 5.0
	for {
		err := Notify(uri, evCh)
		if err == nil {
			fallbackCounter = 0

		} else {
			if fallbackCounter < maxRetries {
				waitTime := time.Duration(int(math.Pow(baseWaitTime, float64(fallbackCounter)))) * time.Second
				log.Info("Encountered Unexted close on SSE connector. Executing exponential fallback reconnect", map[string]interface{}{"counter": fallbackCounter, "waittime": int(waitTime.Seconds())})
				time.Sleep(waitTime)

			} else {
				log.Error("Encountered Unexted close on SSE connector. maximum retries exceeded. terminating", map[string]interface{}{"counter": strconv.Itoa(fallbackCounter), "maxretries": strconv.Itoa(maxRetries)})
				break
			}

			fallbackCounter += 1
		}
	}
}
