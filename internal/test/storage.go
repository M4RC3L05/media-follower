package test

import (
	"github.com/m4rc3l05/media-follower/internal/storage"
)

func NewDB() *storage.Db {
	db, err := storage.New(":memory:")
	if err != nil {
		panic(err)
	}

	contents, err := storage.Schema.ReadFile("schema.sql")
	if err != nil {
		panic(err)
	}

	_, err = db.DB.Exec(string(contents))
	if err != nil {
		panic(err)
	}

	return db
}
