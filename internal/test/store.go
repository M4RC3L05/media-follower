package test

import (
	"math/rand"

	"github.com/m4rc3l05/media-follower/.gen/jetdb/model"
	"github.com/m4rc3l05/media-follower/.gen/jetdb/table"
	"github.com/m4rc3l05/media-follower/internal/commands"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/store"
)

func NewDB() *store.Db {
	db, err := store.New(":memory:")
	if err != nil {
		panic(err)
	}

	if err := commands.DBMigrate(db.DB, common.NewLogger("foo")); err != nil {
		panic(err)
	}

	return db
}

func LoadInput(db *store.Db) model.Inputs {
	var input model.Inputs

	stmt := table.Inputs.
		INSERT(table.Inputs.ID, table.Inputs.Provider, table.Inputs.Raw).
		VALUES(rand.Int(), "foo", store.JSONB([]byte(`{}`))).
		RETURNING(table.Inputs.AllColumns, store.JSONCol(table.Inputs.Raw).AS("inputs.raw"))

	if err := stmt.Query(db.DB, &input); err != nil {
		panic(err)
	}

	return input
}
