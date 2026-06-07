package handlers

import (
	"context"
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	qb "github.com/go-jet/jet/v2/sqlite"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/model"
	"github.com/m4rc3l05/media-follower/.gen/go-jet/table"
	"github.com/m4rc3l05/media-follower/internal/apps/admin_server/views"
	"github.com/m4rc3l05/media-follower/internal/common/middlewares"
	passwordhashing "github.com/m4rc3l05/media-follower/internal/common/password_hashing"
	"github.com/m4rc3l05/media-follower/internal/store"
)

type RegisterRequestBody struct {
	Username string `form:"username" validate:"required,min=1,max=10"`
	Password string `form:"password" validate:"required,min=8,max=26"`
}

type LoginRequestBody struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

type AuthHandler struct {
	Db *store.Db
	Ph passwordhashing.IPasswordHashing
}

func (h AuthHandler) GetLogin(c fiber.Ctx) error {
	return views.Login().
		Render(
			context.WithValue(
				c.Type("html", "utf-8").RequestCtx(),
				views.FlashMessagesContextViewKey,
				middlewares.FlashMessageProviderFromContext(c).Flashes(c),
			),
			c.Res(),
		)
}

func (h AuthHandler) PostLogin(c fiber.Ctx) error {
	body := LoginRequestBody{}
	if err := c.Bind().Form(&body); err != nil {
		return err
	}

	var user model.Users
	qStmt := qb.SELECT(table.Users.AllColumns).
		FROM(table.Users).
		WHERE(table.Users.Username.EQ(qb.String(body.Username))).
		LIMIT(1)

	if err := qStmt.QueryContext(c.RequestCtx(), h.Db.DB, &user); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return c.Redirect().Back()
		}

		return err
	}

	if !h.Ph.Compare(user.Password, body.Password) {
		return c.Redirect().Back()
	}

	sess := session.FromContext(c)

	if err := sess.Regenerate(); err != nil {
		return err
	}

	sess.Set("userId", user.ID)

	return c.Redirect().Route("home")
}

func (h AuthHandler) GetRegister(c fiber.Ctx) error {
	var user []model.Users
	qStmt := qb.SELECT(table.Users.AllColumns).
		FROM(table.Users).
		LIMIT(1)

	if err := qStmt.QueryContext(c.RequestCtx(), h.Db.DB, &user); err != nil {
		return err
	}

	if len(user) > 0 {
		middlewares.
			FlashMessageProviderFromContext(c).
			Flash(c, middlewares.FlashMessageKeyError, "A user was already setup")

		return c.Redirect().Route("auth-login")
	}

	return views.Register().
		Render(
			context.WithValue(
				c.Type("html", "utf-8").RequestCtx(),
				views.FlashMessagesContextViewKey,
				middlewares.FlashMessageProviderFromContext(c).Flashes(c),
			),
			c.Res(),
		)
}

func (h AuthHandler) PostRegister(c fiber.Ctx) error {
	body := RegisterRequestBody{}
	if err := c.Bind().Form(&body); err != nil {
		return err
	}
	hash := h.Ph.Hash(body.Password)

	var user model.Users
	stmt := table.Users.INSERT(table.Users.ID, table.Users.Username, table.Users.Password).
		VALUES("1", body.Username, hash).
		RETURNING(table.Users.AllColumns)

	if err := stmt.QueryContext(c.RequestCtx(), h.Db.DB, &user); err != nil {
		return err
	}

	return c.Redirect().Route("auth-login")
}

func (h AuthHandler) PostLogout(c fiber.Ctx) error {
	sess := session.FromContext(c)
	if err := sess.Destroy(); err != nil {
		return err
	}

	return c.Redirect().Route("auth-login")
}
