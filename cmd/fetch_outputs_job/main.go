package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/jobs"
	"github.com/m4rc3l05/media-follower/internal/store"
)

func run(providerName string, ctx context.Context) (exitCode int) {
	log := common.NewLogger("fetch-outputs-job")
	config, err := common.NewConfig()
	if err != nil {
		log.Error("Unable to load config", slog.Any("error", err))
		return 1
	}

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

	job, err := jobs.ResolveJob(providerName, db, log)
	if err != nil {
		log.Error("Unable to create job", slog.Any("err", err))

		return 1
	}

	log.Info("Running job")

	if err := job.Run(ctx); err != nil {
		log.Error("Error running fetch-outputs job", slog.Any("err", err))

		return 1
	}

	log.Info("Job completed")

	return 0
}

func main() {
	provider := flag.String("p", "", "The provider to fetch outputs from")

	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	os.Exit(run(*provider, ctx))
}
