package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	qb "github.com/go-jet/jet/v2/sqlite"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/table"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/providers"
	"github.com/m4rc3l05/media-follower/internal/storage"
)

type FetchReleasesJob[I any, O any] struct {
	ReleaseProvider providers.IReleaseProvider[I, O]
	DB              *storage.Db
	Log             *slog.Logger
}

// Compile time check that providers implement interface
var (
	_ common.IEntrypoint = FetchReleasesJob[any, any]{}
)

func NewFetchReleasesJob[I any, O any](
	releaseProvider providers.IReleaseProvider[I, O],
	db *storage.Db,
	log *slog.Logger,
) FetchReleasesJob[I, O] {
	return FetchReleasesJob[I, O]{
		ReleaseProvider: releaseProvider,
		DB:              db,
		Log:             log,
	}
}

func (f FetchReleasesJob[I, O]) Run(ctx context.Context) error {
	var inputs []model.Inputs

	stmt := qb.SELECT(
		table.Inputs.AllColumns.Except(table.Inputs.Raw),
		storage.JSONCol(table.Inputs.Raw).AS("inputs.raw"),
	).
		FROM(table.Inputs).
		WHERE(table.Inputs.Provider.EQ(qb.String(f.ReleaseProvider.Name())))

	if err := stmt.QueryContext(ctx, f.DB.DB, &inputs); err != nil {
		return err
	}

	for i, inputPersistance := range inputs {
		if i > 0 && i < len(inputs)-1 {
			common.SleepWithContext(ctx, 5*time.Second)
		}

		f.Log.Info(
			"Handling input",
			slog.Group(
				"inputPersistance",
				slog.String("id", inputPersistance.ID),
				slog.String("provider", inputPersistance.Provider),
			),
		)

		input, err := f.ReleaseProvider.FromPersistanceToInput(inputPersistance)
		if err != nil {
			f.Log.Error("Error getting input from input persistance", slog.Any("err", err))

			continue
		}

		releases, err := f.ReleaseProvider.FetchReleases(*input)
		if err != nil {
			f.Log.Error("Error fetching outputs", slog.Any("err", err))

			continue
		}

		f.Log.Info(fmt.Sprintf("Processing %d releases for input", len(releases)))

		for _, release := range releases {
			persistance, err := f.ReleaseProvider.FromReleaseToPersistance(
				inputPersistance,
				release,
			)
			if err != nil {

				f.Log.Error("Error converting output to persistance", slog.Any("err", err))
				continue
			}

			stmt2 := table.Releases.
				INSERT(
					table.Releases.ID,
					table.Releases.InputID,
					table.Releases.InputProvider,
					table.Releases.Provider,
					table.Releases.ReleasedAt,
					table.Releases.Raw,
				).
				VALUES(
					persistance.ID,
					persistance.InputID,
					persistance.InputProvider,
					persistance.Provider,
					persistance.ReleasedAt,
					storage.JSONB(persistance.Raw),
				).
				ON_CONFLICT(
					table.Releases.ID,
					table.Releases.InputID,
					table.Releases.InputProvider,
					table.Releases.Provider,
				).
				DO_UPDATE(
					qb.SET(
						table.Releases.ReleasedAt.SET(
							qb.String(persistance.ReleasedAt),
						),
						table.Releases.Raw.SET(
							storage.JSONB(persistance.Raw),
						),
					),
				)

			_, err = stmt2.ExecContext(ctx, f.DB.DB)
			if err != nil {
				f.Log.Error(
					"Error storing release persistance in database",
					slog.Group(
						"inputPersistance",
						slog.String("id", inputPersistance.ID),
						slog.String("provider", inputPersistance.Provider),
					),
					slog.Group(
						"releasePersistance",
						slog.String("id", persistance.ID),
						slog.String("inputId", persistance.InputID),
						slog.String("inputProvider", persistance.InputProvider),
						slog.String("raw", string(persistance.Raw)),
					),
					slog.Any("err", err),
				)
				continue
			}
		}

		f.Log.Info(
			fmt.Sprintf("Processed %d outputs for input", len(releases)),
			slog.Group(
				"inputPersistance",
				slog.String("id", inputPersistance.ID),
				slog.String("provider", inputPersistance.Provider),
			),
		)
	}

	return nil
}

func (f FetchReleasesJob[I, O]) Close(ctx context.Context) error {
	f.Log.Info("Closing database")

	defer f.Log.Info("Closed database")

	return f.DB.Close(ctx)
}
