package logging

import (
	"fmt"
	"runtime"
)

type StandardOutLogger struct {
	DebugLogging bool
}

const (
	DEBUG_MSG = "DEBUG %s ; %s:%d ; %s\n"
	INFO_MSG  = "INFO  %s ; %s:%d ; %s\n"
	ERROR_MSG = "ERROR %s ; %s:%d ; %s\n"
	FATAL_MSG = ERROR_MSG
)

func StringifyData(datas []map[string]string) (dataString string) {
	for _, data := range datas {
		for key, val := range data {
			dataString += fmt.Sprintf("%s=%s ", key, val)
		}
	}
	return
}

func (log StandardOutLogger) Debug(msg string, datas ...map[string]string) {
	_, file, no, _ := runtime.Caller(2)
	fmt.Printf(DEBUG_MSG, msg, file, no, StringifyData(datas))
}

func (log StandardOutLogger) Info(msg string, datas ...map[string]string) {
	_, file, no, _ := runtime.Caller(2)
	fmt.Printf(INFO_MSG, msg, file, no, StringifyData(datas))
}

func (log StandardOutLogger) Error(msg string, datas ...map[string]string) {
	_, file, no, _ := runtime.Caller(2)
	fmt.Printf(ERROR_MSG, msg, file, no, StringifyData(datas))
}

func (log StandardOutLogger) Fatal(msg string, datas ...map[string]string) {
	_, file, no, _ := runtime.Caller(2)
	panic(fmt.Sprintf(FATAL_MSG, msg, file, no, StringifyData(datas)))
}
