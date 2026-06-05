package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/m4rc3l05/media-follower/cmd/admin_server/views"
	"github.com/m4rc3l05/media-follower/internal/common"
)

//go:embed all:.dist
var dist embed.FS

func run(ctx context.Context) (exitCode int) {
	log := common.NewLogger("admin-server")
	app := fiber.New(fiber.Config{
		Services: []fiber.Service{},
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

	app.Use("/", static.New(".dist", static.Config{FS: dist}))
	app.Use("/", favicon.New(favicon.Config{FileSystem: dist, File: ".dist/static/favicon.ico"}))

	app.Get("/", func(ctx fiber.Ctx) error {
		ctx.Type("html", "utf-8")
		return views.Hello("foo").Render(ctx.Context(), ctx.Res())
	})

	err := app.Listen(
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
