package testdata

import (
	"math/rand"
	"time"

	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/table"
	"github.com/m4rc3l05/media-follower/internal/storage"
)

func LoadDBInput(db *storage.Db) model.Inputs {
	var input model.Inputs

	stmt := table.Inputs.
		INSERT(table.Inputs.ID, table.Inputs.Provider, table.Inputs.Raw).
		VALUES(rand.Int(), "foo", storage.JSONB([]byte(`{}`))).
		RETURNING(table.Inputs.AllColumns, storage.JSONCol(table.Inputs.Raw).AS("inputs.raw"))

	if err := stmt.Query(db.DB, &input); err != nil {
		panic(err)
	}

	return input
}

func LoadDBRelease(db *storage.Db, input *model.Inputs) model.Releases {
	if input == nil {
		i := LoadDBInput(db)

		input = &i
	}

	var output model.Releases

	stmt := table.Releases.
		INSERT(
			table.Releases.ID,
			table.Releases.InputID,
			table.Releases.InputProvider,
			table.Releases.Provider,
			table.Releases.ReleasedAt,
			table.Releases.Raw,
		).
		VALUES(
			rand.Int(),
			input.ID,
			input.Provider,
			"bar",
			time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00"),
			storage.JSONB([]byte(`{}`)),
		).
		RETURNING(table.Releases.AllColumns, storage.JSONCol(table.Releases.Raw).AS("releases.raw"))

	if err := stmt.Query(db.DB, &output); err != nil {
		panic(err)
	}

	return output
}
