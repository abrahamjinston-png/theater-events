package logger

import (
	"fmt"
	"log"
	"os"
)

var std = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Llongfile)

// Infof logs an informational message.
func Infof(format string, args ...any) {
	std.Output(2, "[INFO] "+fmt.Sprintf(format, args...))
}

// Debugf logs a debug message.
func Debugf(format string, args ...any) {
	std.Output(2, "[DEBUG] "+fmt.Sprintf(format, args...))
}

// Warnf logs a warning message.
func Warnf(format string, args ...any) {
	std.Output(2, "[WARN] "+fmt.Sprintf(format, args...))
}

// Errorf logs an error message.
func Errorf(format string, args ...any) {
	std.Output(2, "[ERROR] "+fmt.Sprintf(format, args...))
}

// Fatalf logs a fatal message and exits the process.
func Fatalf(format string, args ...any) {
	std.Output(2, "[FATAL] "+fmt.Sprintf(format, args...))
	os.Exit(1)
}
