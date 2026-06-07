package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/factories"
)

var log = common.NewLogger("main")

func run(ctx context.Context, t string, name string, args ...string) (exitCode int) {
	entrypoint, err := factories.EntrypointFactory(ctx, t, name, args...)

	defer func() {
		if entrypoint == nil {
			return
		}

		log.Info("Closing entrypoint")
		closeCtx, done := context.WithTimeout(context.Background(), 10*time.Second)
		defer done()

		if err := entrypoint.Close(closeCtx); err != nil {
			exitCode = 1

			log.Error("Error closing entrypoint", slog.Any("error", err))

			return
		}

		log.Info("Closed entrypoint")
	}()

	if err != nil {
		log.Error("Error creating entrypoint", slog.Any("error", err))

		return 1
	}

	log.Info("Attempting to run entrypoint", "type", t, "name", name, "args", args)

	if err := entrypoint.Run(ctx); err != nil {
		log.Error("Error running entrypoint", slog.Any("error", err))

		return 1
	}

	return 0
}

func main() {
	a := os.Args[1:]
	t := a[0]
	n := a[1]
	rest := a[2:]

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	os.Exit(run(ctx, t, n, rest...))
}
