package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/m4rc3l05/media-follower/internal/commands"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/store"
	_ "modernc.org/sqlite"
)

func run() (statusCode int) {
	log := common.NewLogger("gen-db-types")
	config, err := common.NewConfig()
	if err != nil {
		log.Error("Unable to load config", slog.Any("error", err))
		return 1
	}
	db, err := store.New(config.Database.Path)

	defer func() {
		if err := db.Close(context.Background()); err != nil {
			log.Error("Error closing database", slog.Any("err", err))
			statusCode = 1

			return
		}
	}()

	if err != nil {
		log.Error("Error creating database", slog.Any("err", err))

		return 1
	}

	if err := commands.GenDbTypes(db); err != nil {
		log.Error("Error generating database types", slog.Any("err", err))

		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
