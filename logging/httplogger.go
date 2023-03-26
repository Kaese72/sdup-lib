package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type HTTPLogger struct {
	Url url.URL
}

type HTTPLog struct {
	Message  string            `json:"message"`
	Priority int               `json:"priority"`
	Data     map[string]string `json:"data"`
}

type HTTPConfig struct {
	URL string `json:"url"`
}

func (config HTTPConfig) HTTPLogger() (HTTPLogger, error) {
	url, err := url.Parse(config.URL)
	if err != nil {
		return HTTPLogger{}, err
	}
	return HTTPLogger{
		Url: *url,
	}, nil
}

func (logger HTTPLogger) postLog(log HTTPLog) {
	encoded, err := json.Marshal(log)
	if err != nil {
		backupLogger.Log("Failed to log", 4, map[string]string{"reason": err.Error()})
		return
	}
	// We spin of a goroutine because we do not want to wait for this.
	// If logs arrive in the wrong order thats fine for now
	// FIXME Include timestamp in log format
	go func() {
		resp, err := http.Post(logger.Url.String(), "application/json", bytes.NewBuffer(encoded))
		if err != nil {
			backupLogger.Log("Failed to log", 4, map[string]string{"reason": err.Error()})
			return
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			backupLogger.Log("Failed to log", 4, map[string]string{"reason": fmt.Sprintf("unexpected response code while logging, %d", resp.StatusCode)})
			return
		}
	}()
}

func (logger HTTPLogger) Log(msg string, priority int, datas ...map[string]string) {
	logger.postLog(HTTPLog{
		Message:  msg,
		Priority: priority,
		Data:     mergeMaps(append(datas, collectData())),
	})
}
