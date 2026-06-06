package test

import (
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
