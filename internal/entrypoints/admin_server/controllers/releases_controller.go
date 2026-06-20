package controllers

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	qb "github.com/go-jet/jet/v2/sqlite"
	"github.com/labstack/echo/v5"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/table"
	"github.com/m4rc3l05/media-follower/internal/common/middlewares"
	"github.com/m4rc3l05/media-follower/internal/common/providers"
	"github.com/m4rc3l05/media-follower/internal/entrypoints/admin_server/views"
	"github.com/m4rc3l05/media-follower/internal/storage"
)

type ReleasesController struct {
	Db                      *storage.Db
	ReleaseProviderResolver func(name providers.ProviderName) providers.IReleaseProvider
}

type ReleasesIndexPageRequest struct {
	Provider *string `query:"provider" mod:"trim"             validate:"omitempty,providerName"`
	Page     int64   `query:"page"     mod:"default=0,min=0"  validate:"min=0"`
	Limit    int64   `query:"limit"    mod:"default=10,min=0" validate:"min=0"`
}

func (ic ReleasesController) IndexPage(c *echo.Context) error {
	flashProvider := middlewares.FlashMessageProviderFromContext(c)

	var reqData ReleasesIndexPageRequest
	if err := c.Bind(&reqData); err != nil {
		if err2 := flashProvider.Flash(
			c,
			middlewares.FlashMessageKeyError,
			"Could not validate body",
		); err2 != nil {
			return errors.Join(err, err2)
		}

		return c.Redirect(302, "/Releases")
	}

	if err := c.Validate(&reqData); err != nil {
		if err2 := flashProvider.Flash(
			c,
			middlewares.FlashMessageKeyError,
			"Invalid data provided",
		); err2 != nil {
			return errors.Join(err, err2)
		}

		return c.Redirect(302, "/Releases")
	}

	stmt := qb.SELECT(table.Releases.AllColumns, table.Inputs.AllColumns).
		FROM(table.Releases.LEFT_JOIN(table.Inputs, table.Releases.InputID.EQ(table.Inputs.ID)))

	if reqData.Provider != nil {
		stmt.WHERE(
			table.Inputs.Provider.EQ(qb.String(*reqData.Provider)),
		)
	}

	stmt.ORDER_BY(table.Releases.ReleasedAt.DESC(), qb.Raw(fmt.Sprintf("%s.rowid", table.Releases.TableName())).DESC()).
		LIMIT(reqData.Limit).
		OFFSET(reqData.Limit * reqData.Page)

	var releases []struct {
		model.Releases
		Input *model.Inputs
	}
	if err := stmt.QueryContext(c.Request().Context(), ic.Db.DB, &releases); err != nil {
		return err
	}

	flashes := flashProvider.Flashes(c)
	qp := c.QueryParams()
	qp.Set("page", strconv.FormatInt(max(reqData.Page-1, 0), 10))
	previousPage := fmt.Sprintf("?%s", qp.Encode())
	qp.Set("page", strconv.FormatInt(reqData.Page+1, 10))
	nextPage := fmt.Sprintf("?%s", qp.Encode())

	return views.
		Releases(views.ReleasesPageArgs{
			Releases:      releases,
			ProviderNames: providers.PROVIDERS,
			Links: struct {
				NextPage     string
				PreviousPage string
			}{
				NextPage:     nextPage,
				PreviousPage: previousPage,
			},
		}).
		Render(
			context.WithValue(
				c.Request().Context(),
				views.GlobalViewVarsContentViewKey,
				views.GlobalViewVars{
					FlashMessages: &flashes,
				},
			),
			c.Response(),
		)
}
