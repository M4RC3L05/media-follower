package releasesfeedserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	qb "github.com/go-jet/jet/v2/sqlite"
	"github.com/go-playground/mold/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/feeds"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/table"
	"github.com/m4rc3l05/media-follower/internal/common/config"
	"github.com/m4rc3l05/media-follower/internal/common/logging"
	"github.com/m4rc3l05/media-follower/internal/common/utils"
	"github.com/m4rc3l05/media-follower/internal/entrypoints"
	"github.com/m4rc3l05/media-follower/internal/storage"
)

var log = logging.New("releases-feed-server-entrypoint")

type ReleasesFeedEntrypoint struct {
	Config    *config.Config
	DB        *storage.Db
	Validator *validator.Validate
	Modifier  *mold.Transformer
}

var _ entrypoints.IEntrypoint = ReleasesFeedEntrypoint{}

type RequestData struct {
	Provider string `query:"provider" mod:"trim" validate:"required,providerName"`
}

func (a ReleasesFeedEntrypoint) Run(ctx context.Context) error {
	app := echo.NewWithConfig(echo.Config{
		Binder: &utils.CustomBinder{
			Conform:   a.Modifier,
			DefBinder: &echo.DefaultBinder{},
		},
		Validator: &utils.CustomValidator{Validator: a.Validator},
		HTTPErrorHandler: func(c *echo.Context, err error) {
			log.Error("Error from admin server", slog.Any("error", err))

			if err := c.String(http.StatusInternalServerError, "something went wrong"); err != nil {
				log.Error("Could not responde from error handler", slog.Any("error", err))
			}
		},
	})

	app.Use(middleware.Recover())
	app.Use(middleware.Secure())

	app.GET("/", func(c *echo.Context) error {
		var reqData RequestData
		if err := c.Bind(&reqData); err != nil {
			return err
		}

		if err := c.Validate(&reqData); err != nil {
			return err
		}

		var releases []model.Releases
		stmt := qb.SELECT(table.Releases.AllColumns).
			FROM(table.Releases.LEFT_JOIN(table.Inputs, table.Releases.InputID.EQ(table.Inputs.ID))).
			WHERE(
				table.Inputs.Provider.EQ(qb.String(reqData.Provider)).
					AND(
						table.Releases.ReleasedAt.LT_EQ(qb.StringExp(qb.Func("strftime", qb.String("%Y-%m-%dT%H:%M:%fZ"), qb.String("now")))),
					),
			).
			ORDER_BY(table.Releases.ReleasedAt.DESC(), qb.Raw(fmt.Sprintf("%s.rowid", table.Releases.TableName())).DESC()).
			LIMIT(500)

		if err := stmt.QueryContext(c.Request().Context(), a.DB.DB, &releases); err != nil {
			return err
		}

		now := time.Now().UTC()
		feed := &feeds.Feed{
			Title:       fmt.Sprintf("Provider %s releases", reqData.Provider),
			Description: fmt.Sprintf("Latest releases from provider %s", reqData.Provider),
			Id:          fmt.Sprintf("media-follower-%s-releases", reqData.Provider),
			Created:     now,
			Updated:     now,
			Items:       []*feeds.Item{},
		}

		for _, release := range releases {
			item := &feeds.Item{
				Title: release.Title,
				Id:    fmt.Sprintf("%s@%s", reqData.Provider, release.ID),
			}

			tParsed, err := time.Parse(time.RFC3339Nano, release.ReleasedAt)
			if err != nil {
				return err
			}

			item.Created = tParsed

			if release.ExternalLink != nil {
				item.Link = &feeds.Link{Href: *release.ExternalLink}
			}

			if release.Description != nil {
				item.Description = *release.Description
			}

			if release.ImageURL != nil {
				item.Enclosure = &feeds.Enclosure{
					Url:    *release.ImageURL,
					Type:   "image/jpeg",
					Length: "0",
				}
			}

			feed.Items = append(feed.Items, item)
		}

		accepts := c.Request().Header.Get("accept")

		if strings.Contains(accepts, "application/json") ||
			strings.Contains(accepts, "application/feed+json") {
			res, err := feed.ToJSON()
			if err != nil {
				return err
			}

			return c.Blob(200, "application/feed+json", []byte(res))
		}

		if strings.Contains(accepts, "application/atom+xml") {
			res, err := feed.ToAtom()
			if err != nil {
				return err
			}

			return c.Blob(200, "application/atom+xml", []byte(res))
		}

		res, err := feed.ToRss()
		if err != nil {
			return err
		}

		return c.Blob(200, "application/rss+xml", []byte(res))
	})

	startConf := echo.StartConfig{
		Address: fmt.Sprintf(
			"%s:%d",
			a.Config.Entrypoints.ReleasesFeedServer.Host,
			a.Config.Entrypoints.ReleasesFeedServer.Port,
		),
		HideBanner:      true,
		HidePort:        true,
		GracefulTimeout: 10 * time.Second,
		ListenerAddrFunc: func(addr net.Addr) {
			log.Info(fmt.Sprintf("Serving on http://%s", addr.String()))
		},
		OnShutdownError: func(err error) {
			log.Error("Error gracefully shutting down server", slog.Any("error", err))
		},
	}

	if err := startConf.Start(ctx, app); err != nil {
		return err
	}

	return nil
}

func (a ReleasesFeedEntrypoint) Stop(ctx context.Context) error {
	if a.DB == nil {
		return nil
	}

	log.Info("Closing database")

	if err := a.DB.Close(ctx); err != nil {
		return err
	}

	log.Info("Database closed")

	return nil
}
