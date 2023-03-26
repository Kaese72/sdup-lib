package logging

import (
	"fmt"
	"runtime"
)

type StandardOutLogger struct {
}

const (
	DEBUG_MSG = "DEBUG %s ; %s:%d ; %s\n"
	INFO_MSG  = "INFO  %s ; %s:%d ; %s\n"
	ERROR_MSG = "ERROR %s ; %s:%d ; %s\n"
)

func (log StandardOutLogger) Log(msg string, priority int, datas ...map[string]string) {
	label, ok := map[int]string{
		7: "DEBUG",
		3: "ERROR",
	}[priority]
	if !ok {
		label = "ERROR"
	}
	_, file, no, _ := runtime.Caller(2)
	fmt.Printf("%s %s ; %s:%d ; %s\n", label, msg, file, no, stringifyData(datas))
}
