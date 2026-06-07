package adminserver

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/utils/v2"
	"github.com/m4rc3l05/media-follower/internal/apps"
	"github.com/m4rc3l05/media-follower/internal/apps/admin_server/handlers"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/common/middlewares"
	passwordhashing "github.com/m4rc3l05/media-follower/internal/common/password_hashing"
	"github.com/m4rc3l05/media-follower/internal/store"
)

type AdminServerApp struct {
	Db  *store.Db
	app *fiber.App
}

type structValidator struct {
	validate *validator.Validate
}

//go:embed all:.dist
var dist embed.FS
var log = common.NewLogger("admin-server-app")
var (
	_ apps.IApp             = &AdminServerApp{}
	_ fiber.StructValidator = structValidator{}
)

func (v structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func (a *AdminServerApp) Start(ctx context.Context) error {
	a.app = fiber.New(fiber.Config{
		StructValidator: &structValidator{
			validate: validator.New(validator.WithRequiredStructEnabled()),
		},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			log.Error("Caught an error", slog.Any("error", err))

			return c.Type("text").Status(500).Send([]byte("error"))
		},
	})

	a.app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		log.Info(fmt.Sprintf("Serving on http://%s:%s", listenData.Host, listenData.Port))

		return nil
	})

	a.app.Use(recoverer.New())
	a.app.Use("/", static.New(".dist", static.Config{FS: dist}))
	a.app.Use("/", favicon.New(favicon.Config{FileSystem: dist, File: ".dist/static/favicon.ico"}))
	a.app.Use(session.New(session.Config{
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
	a.app.Use(middlewares.FlashMessages())

	indexHandler := handlers.IndexHandler{}
	authHandler := handlers.AuthHandler{Db: a.Db, Ph: passwordhashing.NewArgon2di()}

	a.app.Get("/", middlewares.Auth, indexHandler.GetIndex).Name("home")

	a.app.Get("/auth/register", middlewares.NotAuth, authHandler.GetRegister).Name("auth-register")
	a.app.Post("/auth/register", middlewares.NotAuth, authHandler.PostRegister)
	a.app.Get("/auth/login", middlewares.NotAuth, authHandler.GetLogin).Name("auth-login")
	a.app.Post("/auth/login", middlewares.NotAuth, authHandler.PostLogin)
	a.app.Post("/auth/logout", middlewares.Auth, authHandler.PostLogout)

	return a.app.Listen(
		"127.0.0.1:4321",
		fiber.ListenConfig{
			GracefulContext:       ctx,
			ShutdownTimeout:       10 * time.Second,
			DisableStartupMessage: true,
		},
	)
}

func (a AdminServerApp) Stop(ctx context.Context) error {
	log.Info("Server is shuting down")

	if err := a.app.ShutdownWithContext(ctx); err != nil {
		return err
	}

	log.Info("Server shutdown")

	return nil
}
