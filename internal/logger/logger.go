package logger

import (
	"log/slog"
	"os"
	"time"
)

func New() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo, //TODO: вынести реализацию в env
		AddSource: true,
	}
	h := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(h).With(
		slog.String("string", "purchases-service"),
		slog.Time("time", time.Now()),
	)
	return logger
}
