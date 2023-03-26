package logging

import (
	"fmt"
	"runtime"
	"strconv"
)

type Logger interface {
	Log(string, int, ...map[string]string)
}

type Config struct {
	StdOut bool        `json:"stdout"`
	HTTP   *HTTPConfig `json:"http"`
}

var logger Logger = StandardOutLogger{}

func SetLogger(newLogger Logger) {
	logger = newLogger
}

func InitLoggers(conf Config) error {
	loggers := []Logger{}
	if conf.StdOut {
		loggers = append(loggers, StandardOutLogger{})
	}
	if conf.HTTP != nil {
		httpLogger, err := conf.HTTP.HTTPLogger()
		if err != nil {
			return err
		}
		loggers = append(loggers, httpLogger)
	}
	SetLogger(MultiLogger{Loggers: loggers})
	return nil
}

// backupLogger is used in emergency situations where primary logger fails
var backupLogger Logger = StandardOutLogger{}
var debugLogging bool = false

func Debug(msg string, data ...map[string]string) {
	if debugLogging {
		logger.Log(msg, 7, data...)
	}
}

func Info(msg string, data ...map[string]string) {
	logger.Log(msg, 6, data...)
}

func Error(msg string, data ...map[string]string) {
	logger.Log(msg, 3, data...)
}

func Fatal(msg string, data ...map[string]string) {
	logger.Log(msg, 1, data...)
}

func SetDebugLogging(flag bool) {
	debugLogging = flag
}

func stringifyData(datas ...map[string]string) (dataString string) {
	for _, data := range datas {
		for key, val := range data {
			dataString += fmt.Sprintf("%s=%s ", key, val)
		}
	}
	return
}

func mergeMaps(datas []map[string]string) map[string]string {
	merged := make(map[string]string)
	for _, data := range datas {
		for key, value := range data {
			merged[key] = value
		}
	}
	return merged
}

func collectData() map[string]string {
	_, file, no, _ := runtime.Caller(3)
	return map[string]string{
		"FILE": file,
		"LINE": strconv.Itoa(no),
	}
}
