package logger

import (
	"context"
	"log/slog"
	"os"
)

var std = slog.New(slog.NewJSONHandler(os.Stderr, nil))

// Info logs at INFO.
func Info(msg string, attrs ...slog.Attr) {
	std.LogAttrs(context.Background(), slog.LevelInfo, msg, attrs...)
}

// Error logs at ERROR, recording err under "err". err may be nil.
func Error(msg string, err error, attrs ...slog.Attr) {
	if err != nil {
		attrs = append(attrs, slog.String("err", err.Error()))
	}
	std.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
}

// Fatal logs at ERROR and exits the process.
func Fatal(msg string, err error, attrs ...slog.Attr) {
	Error(msg, err, attrs...)
	os.Exit(1)
}
