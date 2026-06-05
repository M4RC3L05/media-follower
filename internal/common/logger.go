package common

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
)

type Logger struct {
	L *slog.Logger
}

func resolveHandler() slog.Handler {
	if testing.Testing() {
		return slog.DiscardHandler
	} else {
		return slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		)
	}
}

func NewLogger(name string) *slog.Logger {
	return slog.New(resolveHandler()).With(slog.String("name", name))
}

func (s Logger) Printf(format string, v ...interface{}) {
	s.L.Info(fmt.Sprintf(strings.ReplaceAll(format, "\n", ""), v...))
}

func (s Logger) Verbose() bool {
	return true
}
