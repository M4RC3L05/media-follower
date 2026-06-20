package adminserver

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/m4rc3l05/media-follower/internal/common/config"
	"github.com/m4rc3l05/media-follower/internal/common/logging"
	"github.com/m4rc3l05/media-follower/internal/common/middlewares"
	"github.com/m4rc3l05/media-follower/internal/common/providers"
	"github.com/m4rc3l05/media-follower/internal/common/utils"
	"github.com/m4rc3l05/media-follower/internal/entrypoints"
	"github.com/m4rc3l05/media-follower/internal/entrypoints/admin_server/controllers"
	"github.com/m4rc3l05/media-follower/internal/storage"
	"github.com/matthewhartstonge/argon2"
)

var (
	//go:embed dist/*
	dist embed.FS
	log  = logging.New("admin-server-entrypoint")
)

type AdminServerEntrypoint struct {
	Config                  *config.Config
	DB                      *storage.Db
	Validator               *validator.Validate
	Modifier                *mold.Transformer
	ReleaseProviderResolver func(name providers.ProviderName) providers.IReleaseProvider
	Argon2Config            argon2.Config
}

var _ entrypoints.IEntrypoint = AdminServerEntrypoint{}

func (a AdminServerEntrypoint) Run(ctx context.Context) error {
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

	sessionStore := sessions.NewCookieStore([]byte("secret"))

	app.Use(middleware.Recover())
	app.Use(
		middleware.StaticWithConfig(
			middleware.StaticConfig{Filesystem: dist, HTML5: false, Browse: false, Root: "dist"},
		),
	)
	app.Use(middleware.Secure())
	app.Use(middlewares.Session(sessionStore))
	app.Use(middlewares.FlashMessages(sessionStore))
	app.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		Skipper:        middleware.DefaultSkipper,
		TokenLength:    32,
		TokenLookup:    "form:_csrf",
		CookieName:     "_csrf",
		CookiePath:     "/",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
	}))

	homeController := controllers.HomeController{}
	app.GET("", homeController.IndexPage, middlewares.Authenticated())

	authController := controllers.AuthController{DB: a.DB, Argon2Config: a.Argon2Config}
	authGroup := app.Group("/auth")
	authGroup.GET("/register", authController.RegisterPage, middlewares.NotAuthenticated())
	authGroup.POST("/register", authController.Register, middlewares.NotAuthenticated())
	authGroup.GET("/login", authController.LoginPage, middlewares.NotAuthenticated())
	authGroup.POST("/login", authController.Login, middlewares.NotAuthenticated())
	authGroup.POST("/logout", authController.Logout, middlewares.Authenticated())

	inputsController := controllers.InputsController{
		Db:                      a.DB,
		ReleaseProviderResolver: a.ReleaseProviderResolver,
	}
	inputsGroup := app.Group("/inputs", middlewares.Authenticated())
	inputsGroup.GET("", inputsController.IndexPage)
	inputsGroup.GET("/create", inputsController.CreatePage)
	inputsGroup.POST("/create", inputsController.Create)

	releasesController := controllers.ReleasesController{
		Db:                      a.DB,
		ReleaseProviderResolver: a.ReleaseProviderResolver,
	}
	releasesGroup := app.Group("/releases", middlewares.Authenticated())
	releasesGroup.GET("", releasesController.IndexPage)

	startConf := echo.StartConfig{
		Address: fmt.Sprintf(
			"%s:%d",
			a.Config.Entrypoints.AdminServer.Host,
			a.Config.Entrypoints.AdminServer.Port,
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

func (a AdminServerEntrypoint) Stop(ctx context.Context) error {
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
