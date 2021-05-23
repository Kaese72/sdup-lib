package logging

type Logger interface {
	Debug(string, ...map[string]string)
	Info(string, ...map[string]string)
	Error(string, ...map[string]string)
	Fatal(string, ...map[string]string)
}

var logger Logger = StandardOutLogger{DebugLogging: true}
var debugLogging bool = false

func Debug(msg string, data ...map[string]string) {
	if debugLogging {
		logger.Debug(msg, data...)
	}
}

func Info(msg string, data ...map[string]string) {
	logger.Info(msg, data...)
}

func Error(msg string, data ...map[string]string) {
	logger.Error(msg, data...)
}

func Fatal(msg string, data ...map[string]string) {
	logger.Fatal(msg, data...)
}

func SetDebugLogging(flag bool) {
	debugLogging = flag
}
