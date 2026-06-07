package handlers

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/m4rc3l05/media-follower/internal/common/middlewares"
	"github.com/m4rc3l05/media-follower/internal/entrypoints/apps/admin_server/views"
)

type IndexHandler struct{}

func (h IndexHandler) GetIndex(c fiber.Ctx) error {
	return views.Hello("foo").
		Render(
			context.WithValue(
				c.Type("html", "utf-8"),
				views.FlashMessagesContextViewKey,
				middlewares.FlashMessageProviderFromContext(c).Flashes(c),
			),
			c.Res(),
		)
}
