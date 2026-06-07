package test

import (
	"github.com/m4rc3l05/media-follower/internal/store"
)

func NewDB() *store.Db {
	db, err := store.New(":memory:")
	if err != nil {
		panic(err)
	}

	contents, err := store.Schema.ReadFile("schema.sql")
	if err != nil {
		panic(err)
	}

	_, err = db.DB.Exec(string(contents))
	if err != nil {
		panic(err)
	}

	return db
}
