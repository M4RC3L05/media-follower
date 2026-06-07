package testdata

import (
	"fmt"
	"math/rand"

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

func LoadDBOutput(db *storage.Db, input *model.Inputs) model.Outputs {
	if input == nil {
		i := LoadDBInput(db)
		fmt.Printf("i: %v\n", i)
		input = &i
	}

	var output model.Outputs

	stmt := table.Outputs.
		INSERT(
			table.Outputs.ID,
			table.Outputs.InputID,
			table.Outputs.InputProvider,
			table.Outputs.Provider,
			table.Outputs.Raw,
		).
		VALUES(rand.Int(), input.ID, input.Provider, "bar", storage.JSONB([]byte(`{}`))).
		RETURNING(table.Outputs.AllColumns, storage.JSONCol(table.Outputs.Raw).AS("outputs.raw"))

	if err := stmt.Query(db.DB, &output); err != nil {
		panic(err)
	}

	return output
}
