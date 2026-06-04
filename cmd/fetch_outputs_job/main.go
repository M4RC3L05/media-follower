package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/jobs"
	"github.com/m4rc3l05/media-follower/internal/providers"
	"github.com/m4rc3l05/media-follower/internal/store"
)

func createJob(
	providerName string,
	db *store.Db,
	log *slog.Logger,
) (jobs.IJob, error) {
	validator := validator.New(validator.WithRequiredStructEnabled())

	if providerName == "itunes-album-provider" {
		input := providers.NewItunesArtistProvider(validator)
		output := providers.NewItunesAlbumProvider(validator)

		return jobs.NewFetchOutputsJob(
			input,
			output,
			db,
			log.With(
				slog.Group(
					"providers",
					slog.String("input", input.Name()),
					slog.String("output", output.Name()),
				),
			),
		), nil
	}

	return nil, fmt.Errorf("provider \"%s\" not valid", providerName)
}

func run(providerName string, ctx context.Context) (statusCode int) {
	log := common.NewLogger("fetch-outputs-job")

	db, err := store.New()
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("Error closing database", slog.Any("err", err))
			statusCode = 1

			return
		}

		log.Info("Database closed")
	}()

	if err != nil {
		log.Error("Error creating database", slog.Any("err", err))

		return 1
	}

	job, err := createJob(providerName, db, log)
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
