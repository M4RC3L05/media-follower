package testing

import "github.com/m4rc3l05/media-follower/internal/storage"

func NewTestDatabase() *storage.Db {
	db, err := storage.New(":memory:")
	if err != nil {
		panic(err)
	}

	_, err = db.DB.Exec(string(storage.DBSchema))
	if err != nil {
		panic(err)
	}

	return db
}
