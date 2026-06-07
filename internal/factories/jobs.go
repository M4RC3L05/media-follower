package factories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/entrypoints/jobs"
	"github.com/m4rc3l05/media-follower/internal/providers"
)

func fetchOutputsItunesAlbumProviderJobFactory(
	ctx context.Context,
) (*jobs.FetchOutputsJob[providers.ItunesArtist, providers.ItunesAlbum], error) {
	cfg, err := common.NewConfig()
	if err != nil {
		return nil, err
	}

	db, err := dbFactory(ctx, cfg)
	if err != nil {
		return nil, err
	}

	log := common.NewLogger("fetch-outputs")
	validator := validator.New(validator.WithRequiredStructEnabled())
	inputProvder := providers.NewItunesArtistProvider(validator)
	outputProvider := providers.NewItunesAlbumProvider(validator)

	job := jobs.FetchOutputsJob[providers.ItunesArtist, providers.ItunesAlbum]{
		InputProvider:  inputProvder,
		OutputProvider: outputProvider,
		DB:             db,
		Log: log.With(
			slog.Group(
				"providers",
				slog.String("input", inputProvder.Name()),
				slog.String("output", outputProvider.Name()),
			),
		),
	}

	return &job, nil
}

func jobFactory(ctx context.Context, name string, args ...string) (common.IEntrypoint, error) {
	if name == "fetch-outputs" {
		if len(args) != 1 {
			return nil, errors.New("provider must be suplied to create a new fetch-outputs job")
		}

		if args[0] == "itunes-album-provider" {
			return fetchOutputsItunesAlbumProviderJobFactory(ctx)
		}

		return nil, fmt.Errorf("job \"%s\" and provider \"%s\" is not supported", name, args[0])
	}

	return nil, fmt.Errorf("job \"%s\" is not supported", name)
}
