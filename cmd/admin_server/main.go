package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-jet/jet/v2/qrm"
	qb "github.com/go-jet/jet/v2/sqlite"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/utils/v2"
	"github.com/m4rc3l05/media-follower/.gen/jetdb/model"
	"github.com/m4rc3l05/media-follower/.gen/jetdb/table"
	"github.com/m4rc3l05/media-follower/cmd/admin_server/views"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/common/middlewares"
	passwordhashing "github.com/m4rc3l05/media-follower/internal/common/password_hashing"
	"github.com/m4rc3l05/media-follower/internal/store"
)

//go:embed all:.dist
var dist embed.FS

type structValidator struct {
	validate *validator.Validate
}

// Validator needs to implement the Validate method
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func authMiddleware(c fiber.Ctx) error {
	sess := session.FromContext(c)
	fmp := middlewares.FlashMessageProviderFromContext(c)

	if sess.Get("userId") == nil {
		fmp.Flash(c, middlewares.FlashMessageKeyError, "You MUST be authenticated to access")
		return c.Redirect().Route("auth-login")
	}

	return c.Next()
}

func notAuthMiddleware(c fiber.Ctx) error {
	sess := session.FromContext(c)
	fmp := middlewares.FlashMessageProviderFromContext(c)

	if sess.Get("userId") != nil {
		fmp.Flash(c, middlewares.FlashMessageKeyError, "You MUST NOT be authenticated to access")
		return c.Redirect().Route("home")
	}

	return c.Next()
}

type RegisterBody struct {
	Username string `form:"username" validate:"required,min=1,max=10"`
	Password string `form:"password" validate:"required,min=8,max=26"`
}

func run(ctx context.Context) (exitCode int) {
	config, err := common.NewConfig()
	if err != nil {
		log.Error("Unable to load config", slog.Any("error", err))
		return 1
	}

	ph := passwordhashing.NewArgon2di()
	log := common.NewLogger("admin-server")

	db, err := store.New(config.Database.Path)
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("Error closing database", slog.Any("err", err))
			exitCode = 1

			return
		}

		log.Info("Database closed")
	}()

	if err != nil {
		log.Error("Error creating database", slog.Any("err", err))

		return 1
	}

	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{
			validate: validator.New(validator.WithRequiredStructEnabled()),
		},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			log.Error("Caught an error", slog.Any("error", err))

			return c.Type("text").Status(500).Send([]byte("error"))
		},
	})

	defer func() {
		log.Info("Server is shuting down")

		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Error("Error shuting down server", slog.Any("error", err))
			exitCode = 1

			return
		}

		log.Info("Server shutdown")
	}()

	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		log.Info(fmt.Sprintf("Serving on http://%s:%s", listenData.Host, listenData.Port))

		return nil
	})

	app.Use(recoverer.New())
	app.Use("/", static.New(".dist", static.Config{FS: dist}))
	app.Use("/", favicon.New(favicon.Config{FileSystem: dist, File: ".dist/static/favicon.ico"}))

	app.Use(session.New(session.Config{
		CookieSecure:      true,
		CookieHTTPOnly:    true,
		CookieSameSite:    "Strict",
		IdleTimeout:       30 * time.Minute,
		AbsoluteTimeout:   24 * time.Hour,
		CookiePath:        "/",
		CookieSessionOnly: false,
		Extractor:         extractors.FromCookie("sid"),
		KeyGenerator:      utils.SecureToken,
		ErrorHandler: func(c fiber.Ctx, err error) {
			log.Error("Error on session middleware", slog.Any("error", err))
		},
	}))
	app.Use(middlewares.FlashMessages())

	app.Get("/", authMiddleware, func(c fiber.Ctx) error {
		c.Type("html", "utf-8")
		return views.Hello("foo").
			Render(context.WithValue(c.RequestCtx(), views.FlashMessagesContextViewKey, middlewares.FlashMessageProviderFromContext(c).Flashes(c)), c.Res())
	}).Name("home")

	app.Get("/auth/register", notAuthMiddleware, func(c fiber.Ctx) error {
		var user []model.Users
		qStmt := qb.SELECT(table.Users.AllColumns).
			FROM(table.Users).
			LIMIT(1)

		if err := qStmt.QueryContext(c.RequestCtx(), db.DB, &user); err != nil {
			return err
		}

		if len(user) > 0 {
			middlewares.FlashMessageProviderFromContext(c).
				Flash(c, middlewares.FlashMessageKeyError, "A user was already setup")
			return c.Redirect().Route("auth-login")
		}

		c.Type("html", "utf-8")
		return views.Register().
			Render(context.WithValue(c.RequestCtx(), views.FlashMessagesContextViewKey, middlewares.FlashMessageProviderFromContext(c).Flashes(c)), c.Res())
	}).Name("auth-register")
	app.Post("/auth/register", notAuthMiddleware, func(c fiber.Ctx) error {
		body := RegisterBody{}
		if err := c.Bind().Form(&body); err != nil {
			return err
		}
		hash := ph.Hash(body.Password)

		var user model.Users
		stmt := table.Users.INSERT(table.Users.ID, table.Users.Username, table.Users.Password).
			VALUES("1", body.Username, hash).
			RETURNING(table.Users.AllColumns)

		if err := stmt.QueryContext(c.RequestCtx(), db.DB, &user); err != nil {
			return err
		}

		return c.Redirect().Route("auth-login")
	})
	app.Get("/auth/login", notAuthMiddleware, func(c fiber.Ctx) error {
		c.Type("html", "utf-8")
		return views.Login().
			Render(context.WithValue(c.RequestCtx(), views.FlashMessagesContextViewKey, middlewares.FlashMessageProviderFromContext(c).Flashes(c)), c.Res())
	}).Name("auth-login")
	app.Post("/auth/login", notAuthMiddleware, func(c fiber.Ctx) error {
		body := RegisterBody{}
		if err := c.Bind().Form(&body); err != nil {
			return err
		}

		var user model.Users
		qStmt := qb.SELECT(table.Users.AllColumns).
			FROM(table.Users).
			WHERE(table.Users.Username.EQ(qb.String(body.Username))).
			LIMIT(1)

		if err := qStmt.QueryContext(c.RequestCtx(), db.DB, &user); err != nil {
			if errors.Is(err, qrm.ErrNoRows) {
				return c.Redirect().Back()
			}

			return err
		}

		if !ph.Compare(user.Password, body.Password) {
			return c.Redirect().Back()
		}

		sess := session.FromContext(c)

		if err := sess.Regenerate(); err != nil {
			return err
		}

		sess.Set("userId", user.ID)

		return c.Redirect().Route("home")
	})
	app.Post("/auth/logout", authMiddleware, func(c fiber.Ctx) error {
		sess := session.FromContext(c)
		if err := sess.Destroy(); err != nil {
			return err
		}

		return c.Redirect().Route("auth-login")
	})

	err = app.Listen(
		"127.0.0.1:4321",
		fiber.ListenConfig{
			GracefulContext:       ctx,
			ShutdownTimeout:       10 * time.Second,
			DisableStartupMessage: true,
		},
	)
	if err != nil {
		log.Error("Unable to start server", slog.Any("error", err))

		return 1
	}

	return 0
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	os.Exit(run(ctx))
}
