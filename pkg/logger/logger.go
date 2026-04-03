package logger

import (
	"log/slog"
	"os"
)

// New creates a structured JSON logger with the given minimum log level
func New(level string) *slog.Logger {
	var l slog.Level

	switch level {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: l,
	})

	return slog.New(handler)
}
