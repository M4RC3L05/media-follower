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
	"github.com/m4rc3l05/media-follower/internal/store"
)

type FetchOutputsJob[I any, O any] struct {
	InputProvider  providers.IInputProvider[I]
	OutputProvider providers.IOutputProvider[I, O]
	DB             *store.Db
	Log            *slog.Logger
}

// Compile time check that providers implement interface
var (
	_ IJob = FetchOutputsJob[any, any]{}
)

func NewFetchOutputsJob[I any, O any](
	inputProvider providers.IInputProvider[I],
	outputProvider providers.IOutputProvider[I, O],
	db *store.Db,
	log *slog.Logger,
) FetchOutputsJob[I, O] {
	return FetchOutputsJob[I, O]{
		InputProvider:  inputProvider,
		OutputProvider: outputProvider,
		DB:             db,
		Log:            log,
	}
}

func (f FetchOutputsJob[I, O]) Run(ctx context.Context) error {
	var inputs []model.Inputs

	stmt := qb.SELECT(
		table.Inputs.AllColumns.Except(table.Inputs.Raw),
		store.JSONCol(table.Inputs.Raw).AS("inputs.raw"),
	).
		FROM(table.Inputs).
		WHERE(table.Inputs.Provider.EQ(qb.String(f.InputProvider.Name())))

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

		input, err := f.InputProvider.FromPersistanceToInput(inputPersistance)
		if err != nil {
			f.Log.Error("Error getting input from input persistance", slog.Any("err", err))

			continue
		}

		outputs, err := f.OutputProvider.FetchOutputs(*input)
		if err != nil {
			f.Log.Error("Error fetching outputs", slog.Any("err", err))

			continue
		}

		f.Log.Info(fmt.Sprintf("Processing %d outputs for input", len(outputs)))

		for _, output := range outputs {
			persistance, err := f.OutputProvider.FromOutputToPersistance(inputPersistance, output)
			if err != nil {

				f.Log.Error("Error converting output to persistance", slog.Any("err", err))
				continue
			}

			stmt2 := table.Outputs.
				INSERT(
					table.Outputs.ID,
					table.Outputs.InputID,
					table.Outputs.InputProvider,
					table.Inputs.Provider,
					table.Outputs.Raw,
				).
				VALUES(
					persistance.ID,
					persistance.InputID,
					persistance.InputProvider,
					persistance.Provider,
					store.JSONB(persistance.Raw),
				).
				ON_CONFLICT(
					table.Outputs.ID,
					table.Outputs.InputID,
					table.Outputs.InputProvider,
					table.Outputs.Provider,
				).
				DO_UPDATE(
					qb.SET(
						table.Outputs.Raw.SET(
							store.JSONB(persistance.Raw),
						),
					),
				)

			_, err = stmt2.ExecContext(ctx, f.DB.DB)
			if err != nil {
				f.Log.Error(
					"Error storing output persistance in database",
					slog.Group(
						"inputPersistance",
						slog.String("id", inputPersistance.ID),
						slog.String("provider", inputPersistance.Provider),
					),
					slog.Group(
						"outputPersistance",
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
			fmt.Sprintf("Processed %d outputs for input", len(outputs)),
			slog.Group(
				"inputPersistance",
				slog.String("id", inputPersistance.ID),
				slog.String("provider", inputPersistance.Provider),
			),
		)
	}

	return nil
}
