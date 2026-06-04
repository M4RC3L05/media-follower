package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/store"
)

func run() (statusCode int) {
	log := common.NewLogger("db-migrate")
	db, err := store.New()

	defer func() {
		if err := db.Close(); err != nil {
			log.Error("Error closing database", slog.Any("err", err))
			statusCode = 1

			return
		}
	}()

	if err != nil {
		log.Error("Error creating database", slog.Any("err", err))

		return -1
	}

	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		log.Error("Error crating migration driver", slog.Any("err", err))

		return -1
	}

	d, err := iofs.New(store.Migrations, "migrations")
	if err != nil {
		log.Error("Error creating migration fs", slog.Any("err", err))

		return -1
	}

	m, err := migrate.NewWithInstance("iofs", d, "sqlite", driver)
	if err != nil {
		log.Error("Error creating migrator instance", slog.Any("err", err))

		return -1
	}

	m.Log = common.Logger{L: log}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Error("Error running migrations", slog.Any("err", err))
			return -1
		}
	}

	return 0
}

func main() {
	os.Exit(run())
}
