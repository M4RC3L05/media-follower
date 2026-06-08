package handlers

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/m4rc3l05/media-follower/internal/common/middlewares"
	"github.com/m4rc3l05/media-follower/internal/entrypoints/apps/admin_server/views"
)

type IndexHandler struct{}

func (h IndexHandler) GetIndex(c fiber.Ctx) error {
	csrfToken := csrf.TokenFromContext(c)

	return views.Index("foo").
		Render(
			context.WithValue(
				context.WithValue(
					c.Type("html", "utf-8"),
					views.FlashMessagesContextViewKey,
					middlewares.FlashMessageProviderFromContext(c).Flashes(c),
				),
				views.CSRFTokenCOntextViewKey,
				csrfToken,
			),
			c.Res(),
		)
}
