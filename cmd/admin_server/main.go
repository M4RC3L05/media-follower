package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3/log"
	adminserver "github.com/m4rc3l05/media-follower/internal/apps/admin_server"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/store"
)

func run(ctx context.Context) (exitCode int) {
	config, err := common.NewConfig()
	if err != nil {
		log.Error("Unable to load config", slog.Any("error", err))
		return 1
	}

	log := common.NewLogger("admin-server")

	db, err := store.New(config.Database.Path)
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("Error closing database", slog.Any("err", err))
			exitCode = 1

			return
		}

		log.Info("Database closed")
	}()

	if err != nil {
		log.Error("Error creating database", slog.Any("err", err))

		return 1
	}

	app := adminserver.AdminServerApp{Db: db}

	defer func() {
		closeCtx, done := context.WithTimeout(context.Background(), 10*time.Second)
		defer done()

		if err := app.Stop(closeCtx); err != nil {
			log.Error("Error stoping app", slog.Any("err", err))

			exitCode = 1
			return
		}
	}()

	if err := app.Start(ctx); err != nil {
		log.Error("Error Starting admin server", slog.Any("err", err))

		return 1
	}

	return 0
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	os.Exit(run(ctx))
}
