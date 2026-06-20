package controllers

import (
	"context"
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	qb "github.com/go-jet/jet/v2/sqlite"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/table"
	"github.com/m4rc3l05/media-follower/internal/common/middlewares"
	"github.com/m4rc3l05/media-follower/internal/entrypoints/admin_server/views"
	"github.com/m4rc3l05/media-follower/internal/storage"
	"github.com/matthewhartstonge/argon2"
)

type AuthController struct {
	Validator    *validator.Validate
	DB           *storage.Db
	Argon2Config argon2.Config
}

func (ac AuthController) RegisterPage(c *echo.Context) error {
	token, err := echo.ContextGet[string](c, "csrf")
	if err != nil {
		return err
	}

	flashProvider := middlewares.FlashMessageProviderFromContext(c)
	flashes := flashProvider.Flashes(c)

	return views.Register().Render(
		context.WithValue(
			c.Request().Context(),
			views.GlobalViewVarsContentViewKey,
			views.GlobalViewVars{
				CSRFToken:     &token,
				FlashMessages: &flashes,
			},
		),
		c.Response(),
	)
}

type RegisterRequest struct {
	Username string `form:"username" validate:"required,min=3,max=10"`
	Password string `form:"password" validate:"required,min=8,max=26"`
}

func (ac AuthController) Register(c *echo.Context) error {
	var reqData RegisterRequest

	if err := c.Bind(&reqData); err != nil {
		return err
	}

	if err := c.Validate(reqData); err != nil {
		return err
	}

	var user model.Users
	stmt := qb.SELECT(table.Users.ID).FROM(table.Users).LIMIT(1)
	if err := stmt.QueryContext(c.Request().Context(), ac.DB.DB, &user); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	flashProvider := middlewares.FlashMessageProviderFromContext(c)

	// We have a user with an id
	if user.ID != "" {
		err := flashProvider.Flash(
			c,
			middlewares.FlashMessageKeyError,
			"A user already exists, please login",
		)
		if err != nil {
			return err
		}

		return c.Redirect(302, "/auth/login")
	}

	hash, err := ac.Argon2Config.HashEncoded([]byte(reqData.Password))
	if err != nil {
		err := flashProvider.Flash(
			c,
			middlewares.FlashMessageKeyError,
			"Could not create a user",
		)
		if err != nil {
			return err
		}

		return c.Redirect(302, "/auth/login")
	}

	stmt2 := table.Users.INSERT(table.Users.ID, table.Users.Username, table.Users.Password).
		VALUES("1", reqData.Username, string(hash))

	if _, err := stmt2.ExecContext(c.Request().Context(), ac.DB.DB); err != nil {
		return err
	}

	return c.Redirect(302, "/auth/login")
}

func (ac AuthController) LoginPage(c *echo.Context) error {
	token, err := echo.ContextGet[string](c, "csrf")
	if err != nil {
		return err
	}

	flashProvider := middlewares.FlashMessageProviderFromContext(c)
	flashes := flashProvider.Flashes(c)

	return views.Login().Render(
		context.WithValue(
			c.Request().Context(),
			views.GlobalViewVarsContentViewKey,
			views.GlobalViewVars{
				CSRFToken:     &token,
				FlashMessages: &flashes,
			},
		),
		c.Response(),
	)
}

type LoginRequest struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

func (ac AuthController) Login(c *echo.Context) error {
	var reqData LoginRequest

	if err := c.Bind(&reqData); err != nil {
		return err
	}

	if err := c.Validate(reqData); err != nil {
		return err
	}

	sess := middlewares.SessionProviderFromContext(c)
	flashProvider := middlewares.FlashMessageProviderFromContext(c)

	var user model.Users
	stmt := qb.SELECT(table.Users.AllColumns).FROM(table.Users).LIMIT(1)
	if err := stmt.QueryContext(c.Request().Context(), ac.DB.DB, &user); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			err2 := flashProvider.Flash(
				c,
				middlewares.FlashMessageKeyError,
				"No user created, please register a user",
			)

			if err2 != nil {
				return errors.Join(err, err2)
			}

			return c.Redirect(302, "/auth/register")
		}

		return err
	}

	r, err := argon2.Decode([]byte(user.Password))
	if err != nil {
		return err
	}

	same, err := r.Verify([]byte(reqData.Password))
	if err != nil {
		return err
	}

	if !same {
		err := flashProvider.Flash(
			c,
			middlewares.FlashMessageKeyError,
			"Could not login, confirm username & password",
		)
		if err != nil {
			return err
		}

		return c.Redirect(302, "/auth/login")
	}

	if err := sess.Authenticate(c, user.ID); err != nil {
		return err
	}

	if err := flashProvider.Flash(
		c,
		middlewares.FlashMessageKeySuccess,
		"Login successfull",
	); err != nil {
		return err
	}

	return c.Redirect(302, "/")
}

func (ac AuthController) Logout(c *echo.Context) error {
	sess := middlewares.SessionProviderFromContext(c)
	flashProvider := middlewares.FlashMessageProviderFromContext(c)

	if err := sess.Logout(c); err != nil {
		return err
	}

	if err := flashProvider.Flash(
		c,
		middlewares.FlashMessageKeySuccess,
		"Logout successfull",
	); err != nil {
		return err
	}

	return c.Redirect(302, "/auth/login")
}
