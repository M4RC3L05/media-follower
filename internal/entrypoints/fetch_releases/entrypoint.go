package fetchreleases

import (
	"context"
	"log/slog"
	"time"

	qb "github.com/go-jet/jet/v2/sqlite"
	"github.com/google/uuid"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/table"
	"github.com/m4rc3l05/media-follower/internal/common/logging"
	"github.com/m4rc3l05/media-follower/internal/common/providers"
	"github.com/m4rc3l05/media-follower/internal/common/utils"
	"github.com/m4rc3l05/media-follower/internal/entrypoints"
	"github.com/m4rc3l05/media-follower/internal/storage"
)

var log = logging.New("fetch-releases-entrypoint")

type FetchReleasesEntrypoint struct {
	Provider providers.IReleaseProvider
	DB       *storage.Db
}

var _ entrypoints.IEntrypoint = FetchReleasesEntrypoint{}

func (f FetchReleasesEntrypoint) Run(ctx context.Context) error {
	var inputs []model.Inputs

	err := qb.SELECT(table.Inputs.AllColumns).
		FROM(table.Inputs).
		WHERE(table.Inputs.Provider.EQ(qb.String(string(f.Provider.Name())))).
		QueryContext(ctx, f.DB.DB, &inputs)
	if err != nil {
		return err
	}

	inputsSize := len(inputs)
	for key, input := range inputs {
		if key > 0 {
			log.Info("Waiting 5 seconds before processing next input")

			utils.SleepWithContext(ctx, 5*time.Second)
		}

		log.Info("Processing input", "inputName", input.Name, "n", key+1, "total", inputsSize)

		fetchedReleases, err := f.Provider.FetchReleases(input)
		if err != nil {
			log.Warn(
				"Failed to fetch releases for input, skipping",
				"inputName",
				input.Name,
				slog.Any("error", err),
			)

			continue
		}

		log.Info("Storing releases", "inputName", input.Name, "total", len(fetchedReleases))

		for _, fetchedRelease := range fetchedReleases {
			stmt := table.Releases.INSERT(
				table.Releases.ID,
				table.Releases.Description,
				table.Releases.ImageURL,
				table.Releases.InternalProviderID,
				table.Releases.ReleasedAt,
				table.Releases.Title,
				table.Releases.InputID,
				table.Releases.ExternalLink,
			).
				VALUES(
					qb.String(uuid.Must(uuid.NewRandom()).String()),
					fetchedRelease.Description,
					fetchedRelease.ImageURL,
					fetchedRelease.InternalProviderID,
					fetchedRelease.ReleasedAt,
					fetchedRelease.Title,
					fetchedRelease.InputID,
					fetchedRelease.ExternalLink,
				).
				ON_CONFLICT(table.Releases.InputID, table.Releases.InternalProviderID).
				DO_UPDATE(
					qb.SET(
						table.Releases.ImageURL.SET(table.Releases.EXCLUDED.ImageURL),
						table.Releases.Description.SET(table.Releases.EXCLUDED.Description),
						table.Releases.ReleasedAt.SET(table.Releases.EXCLUDED.ReleasedAt),
						table.Releases.Title.SET(table.Releases.EXCLUDED.Title),
						table.Releases.ExternalLink.SET(table.Releases.EXCLUDED.ExternalLink),
					),
				)

			_, err := stmt.ExecContext(ctx, f.DB.DB)
			if err != nil {
				log.Warn(
					"Error saving release for input, ignoring",
					"inputName",
					input.Name,
					slog.Any("release", fetchedRelease),
					slog.Any("error", err),
				)
			}
		}
	}

	_, err = f.DB.DB.ExecContext(ctx, "PRAGMA wal_checkpoint(TRUNCATE)")
	if err != nil {
		return err
	}

	log.Info("Fetch releases done")

	return nil
}

func (f FetchReleasesEntrypoint) Stop(ctx context.Context) error {
	if f.DB == nil {
		return nil
	}

	log.Info("Closing database")

	if err := f.DB.Close(ctx); err != nil {
		return err
	}

	log.Info("Database closed")

	return nil
}
