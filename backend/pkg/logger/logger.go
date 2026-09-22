package logger

import (
	"log/slog"
	"os"
)

// New initializes a JSON slog logger
func New() *slog.Logger {
	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)

	return slog.New(handler)
}
