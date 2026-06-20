package controllers

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	qb "github.com/go-jet/jet/v2/sqlite"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/table"
	"github.com/m4rc3l05/media-follower/internal/common/middlewares"
	"github.com/m4rc3l05/media-follower/internal/common/providers"
	"github.com/m4rc3l05/media-follower/internal/entrypoints/admin_server/views"
	"github.com/m4rc3l05/media-follower/internal/storage"
)

type InputsController struct {
	Db                      *storage.Db
	ReleaseProviderResolver func(name providers.ProviderName) providers.IReleaseProvider
}

type InputsIndexPageRequest struct {
	Provider *string `query:"provider" mod:"trim"             validate:"omitempty,providerName"`
	Page     int64   `query:"page"     mod:"default=0,min=0"  validate:"min=0"`
	Limit    int64   `query:"limit"    mod:"default=10,min=0" validate:"min=0"`
}

func (ic InputsController) IndexPage(c *echo.Context) error {
	flashProvider := middlewares.FlashMessageProviderFromContext(c)

	var reqData InputsIndexPageRequest
	if err := c.Bind(&reqData); err != nil {
		if err2 := flashProvider.Flash(
			c,
			middlewares.FlashMessageKeyError,
			"Could not validate body",
		); err2 != nil {
			return errors.Join(err, err2)
		}

		return c.Redirect(302, "/inputs")
	}

	if err := c.Validate(&reqData); err != nil {
		if err2 := flashProvider.Flash(
			c,
			middlewares.FlashMessageKeyError,
			"Invalid data provided",
		); err2 != nil {
			return errors.Join(err, err2)
		}

		return c.Redirect(302, "/inputs")
	}

	stmt := qb.SELECT(table.Inputs.AllColumns).
		FROM(table.Inputs).
		WHERE(table.Inputs.Provider.NOT_EQ(qb.String("__internal_deleted_input__")))

	if reqData.Provider != nil {
		stmt.WHERE(table.Inputs.Provider.EQ(qb.String(*reqData.Provider)))
	}

	stmt.LIMIT(reqData.Limit).OFFSET(reqData.Limit * reqData.Page)

	var inputs []model.Inputs
	if err := stmt.QueryContext(c.Request().Context(), ic.Db.DB, &inputs); err != nil {
		return err
	}

	flashes := flashProvider.Flashes(c)
	qp := c.QueryParams()
	qp.Set("page", strconv.FormatInt(max(reqData.Page-1, 0), 10))
	previousPage := fmt.Sprintf("?%s", qp.Encode())
	qp.Set("page", strconv.FormatInt(reqData.Page+1, 10))
	nextPage := fmt.Sprintf("?%s", qp.Encode())

	return views.
		Inputs(views.InputsPageArgs{
			Inputs:        inputs,
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

func (ic InputsController) CreatePage(c *echo.Context) error {
	flashProvider := middlewares.FlashMessageProviderFromContext(c)
	flashes := flashProvider.Flashes(c)

	return views.
		InputsCreate(views.InputsCreatePageArgs{
			ProviderNames: providers.PROVIDERS,
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

type InputsCreatePageRequest struct {
	Provider string `form:"provider" mod:"trim" validate:"required,providerName"`
	Term     string `form:"term"     mod:"trim" validate:"required,min=1"`
}

func (ic InputsController) Create(c *echo.Context) error {
	flashProvider := middlewares.FlashMessageProviderFromContext(c)

	var reqData InputsCreatePageRequest
	if err := c.Bind(&reqData); err != nil {
		if err2 := flashProvider.Flash(
			c,
			middlewares.FlashMessageKeyError,
			"Could not validate body",
		); err2 != nil {
			return errors.Join(err, err2)
		}

		return c.Redirect(302, "/inputs/create")
	}

	if err := c.Validate(&reqData); err != nil {
		if err2 := flashProvider.Flash(
			c,
			middlewares.FlashMessageKeyError,
			"Invalid data provided",
		); err2 != nil {
			return errors.Join(err, err2)
		}

		return c.Redirect(302, "/inputs/create")
	}

	provider := ic.ReleaseProviderResolver(providers.ProviderName(reqData.Provider))
	input, err := provider.LookupReleaseInput(reqData.Term)
	if err != nil {
		if err2 := flashProvider.Flash(
			c,
			middlewares.FlashMessageKeyError,
			"Error fetching for input",
		); err2 != nil {
			return errors.Join(err, err2)
		}

		return c.Redirect(302, "/inputs/create")
	}

	stmt := table.Inputs.INSERT(
		table.Inputs.ID,
		table.Inputs.Description,
		table.Inputs.InternalProviderID,
		table.Inputs.Name,
		table.Inputs.ImageURL,
		table.Inputs.Provider,
		table.Inputs.ExternalLink,
	).VALUES(
		uuid.New().String(),
		input.Description,
		input.InternalProviderID,
		input.Name,
		input.ImageURL,
		input.Provider,
		input.ExternalLink,
	).RETURNING(table.Inputs.AllColumns)

	var inserted model.Inputs
	if err := stmt.QueryContext(c.Request().Context(), ic.Db.DB, &inserted); err != nil {
		return err
	}

	if err := flashProvider.Flash(
		c,
		middlewares.FlashMessageKeySuccess,
		fmt.Sprintf("Input %s (%s) created", inserted.ID, inserted.Name),
	); err != nil {
		return errors.Join(err, err)
	}

	return c.Redirect(302, "/inputs")
}
