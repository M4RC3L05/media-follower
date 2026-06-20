package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/m4rc3l05/media-follower/internal/common/config"
	"github.com/m4rc3l05/media-follower/internal/common/factories"
	"github.com/m4rc3l05/media-follower/internal/common/logging"
	"github.com/m4rc3l05/media-follower/internal/common/providers"
	"github.com/m4rc3l05/media-follower/internal/common/utils"
	"github.com/m4rc3l05/media-follower/internal/entrypoints"
	adminserver "github.com/m4rc3l05/media-follower/internal/entrypoints/admin_server"
	fetchreleases "github.com/m4rc3l05/media-follower/internal/entrypoints/fetch_releases"
	releasesfeedserver "github.com/m4rc3l05/media-follower/internal/entrypoints/releases_feed_server"
	"github.com/m4rc3l05/media-follower/internal/storage"
	"github.com/matthewhartstonge/argon2"
	_ "golang.org/x/crypto/x509roots/fallback"
)

func entrypointFactory(entrypoint string, args ...string) (entrypoints.IEntrypoint, error) {
	if entrypoint == "admin-server" {
		cfg, err := config.New()
		if err != nil {
			return nil, err
		}

		db, err := storage.New(cfg.Database.Path)
		if err != nil {
			return nil, err
		}

		return adminserver.AdminServerEntrypoint{
			Config:                  cfg,
			DB:                      db,
			Validator:               utils.NewValidator(),
			Modifier:                utils.NewModifier(),
			ReleaseProviderResolver: factories.ProviderFactory,
			Argon2Config:            argon2.DefaultConfig(),
		}, nil
	}

	if entrypoint == "fetch-releases" {
		cfg, err := config.New()
		if err != nil {
			return nil, err
		}

		provider := factories.ProviderFactory(providers.ProviderName(args[0]))

		db, err := storage.New(cfg.Database.Path)
		if err != nil {
			return nil, err
		}

		return fetchreleases.FetchReleasesEntrypoint{Provider: provider, DB: db}, nil
	}

	if entrypoint == "releases-feed-server" {
		cfg, err := config.New()
		if err != nil {
			return nil, err
		}

		db, err := storage.New(cfg.Database.Path)
		if err != nil {
			return nil, err
		}

		return releasesfeedserver.ReleasesFeedEntrypoint{
			Config:    cfg,
			DB:        db,
			Validator: utils.NewValidator(),
			Modifier:  utils.NewModifier(),
		}, nil
	}

	return nil, fmt.Errorf("entrypoint \"%s\" does not exists", entrypoint)
}

func run(ctx context.Context, entrypoint string, args ...string) (exitCode int) {
	log := logging.New("main")

	defer func() {
		if x := recover(); x != nil {

			log.Error(
				"Panic caught",
				slog.Any("panic", x),
				slog.String("stack", string(debug.Stack())),
			)

			exitCode = 1
		}
	}()

	e, err := entrypointFactory(entrypoint, args...)

	defer func() {
		if e == nil {
			return
		}

		closeContext, done := context.WithTimeout(context.Background(), 10*time.Second)
		defer done()

		log.Info("Closing entrypoint")

		if err := e.Stop(closeContext); err != nil {
			log.Error("Unable to properly close entrypoint", slog.Any("error", err))
		}

		log.Info("Entrypoint closed")
	}()

	if err != nil {
		log.Error("Unable to resolve entrypoint", slog.Any("error", err))

		return 1
	}

	log.Info("Running entrypoint", "entrypoint", entrypoint)

	if err := e.Run(ctx); err != nil {
		log.Error("Error starting server", slog.Any("error", err))

		return 1
	}

	return 0
}

func getValOrEmpty(x []string, index int) string {
	if index < len(x) {
		return x[index]
	}

	return ""
}

func main() {
	args := os.Args[1:]
	entrypoint := getValOrEmpty(args, 0)
	rest := args[1:]

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	os.Exit(run(ctx, entrypoint, rest...))
}
