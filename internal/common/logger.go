package common

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	L *slog.Logger
}

func NewLogger(name string) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		),
	).With(slog.String("name", name))
}

func (s Logger) Printf(format string, v ...interface{}) {
	s.L.Info(fmt.Sprintf(strings.ReplaceAll(format, "\n", ""), v...))
}

func (s Logger) Verbose() bool {
	return true
}
