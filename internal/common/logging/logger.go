package logging

import (
	"log/slog"
	"os"
	"testing"

	"github.com/Marlliton/slogpretty"
)

func resolveLogHandler() slog.Handler {
	if os.Getenv("ENV") == "production" {
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	if testing.Testing() {
		return slog.DiscardHandler
	}

	return slogpretty.New(
		os.Stdout,
		&slogpretty.Options{
			Level:      slog.LevelDebug,
			AddSource:  true,
			Colorful:   true,
			Multiline:  true,
			TimeFormat: slogpretty.DefaultTimeFormat,
		},
	)
}

func New(name string) *slog.Logger {
	return slog.New(resolveLogHandler()).With("name", name)
}
