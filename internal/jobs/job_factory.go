package jobs

import (
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/m4rc3l05/media-follower/internal/providers"
	"github.com/m4rc3l05/media-follower/internal/store"
)

func ResolveJob(
	providerName string,
	db *store.Db,
	log *slog.Logger,
) (IJob, error) {
	validator := validator.New(validator.WithRequiredStructEnabled())

	if providerName == "itunes-album-provider" {
		input := providers.NewItunesArtistProvider(validator)
		output := providers.NewItunesAlbumProvider(validator)

		return NewFetchOutputsJob(
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
