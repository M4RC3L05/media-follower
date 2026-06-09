package factories

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/entrypoints/jobs"
	"github.com/m4rc3l05/media-follower/internal/providers"
	"github.com/m4rc3l05/media-follower/internal/providers/inputs"
	"github.com/m4rc3l05/media-follower/internal/providers/outputs"
)

func fetchReleasesItunesMusicReleasesProviderJobFactory(
	ctx context.Context,
) (*jobs.FetchReleasesJob[inputs.ItunesArtist, outputs.ItunesAlbum], error) {
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
	releaseProvider := providers.NewItunesMusicReleasesProvider(
		inputs.NewItunesArtistProvider(validator),
		outputs.NewItunesAlbumProvider(validator),
	)

	job := jobs.FetchReleasesJob[inputs.ItunesArtist, outputs.ItunesAlbum]{
		ReleaseProvider: releaseProvider,
		DB:              db,
		Log:             log.With("provider", releaseProvider.Name()),
	}

	return &job, nil
}

func jobFactory(ctx context.Context, name string, args ...string) (common.IEntrypoint, error) {
	if name == "fetch-releases" {
		if len(args) != 1 {
			return nil, errors.New("provider must be suplied to create a new fetch-releases job")
		}

		if args[0] == "itunes-music-releases-provider" {
			return fetchReleasesItunesMusicReleasesProviderJobFactory(ctx)
		}

		return nil, fmt.Errorf("job \"%s\" and provider \"%s\" is not supported", name, args[0])
	}

	return nil, fmt.Errorf("job \"%s\" is not supported", name)
}
