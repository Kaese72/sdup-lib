package logging

import (
	"fmt"
)

type StandardOutLogger struct {
}

func (log StandardOutLogger) Log(msg string, priority int, datas ...map[string]string) {
	label, ok := map[int]string{
		7: "DEBUG",
		6: "INFO",
		4: "WARNING",
		// Default: 3: "ERROR",
	}[priority]
	if !ok {
		label = "ERROR"
	}
	fmt.Printf("%s %s ; %s\n", label, msg, stringifyData(mergeMaps(datas), collectData()))
}
