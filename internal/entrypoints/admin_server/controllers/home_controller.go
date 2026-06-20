package controllers

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/m4rc3l05/media-follower/internal/common/middlewares"
	"github.com/m4rc3l05/media-follower/internal/entrypoints/admin_server/views"
)

type HomeController struct{}

func (hc HomeController) IndexPage(c *echo.Context) error {
	token, err := echo.ContextGet[string](c, "csrf")
	if err != nil {
		return err
	}

	flashProvider := middlewares.FlashMessageProviderFromContext(c)
	flashes := flashProvider.Flashes(c)

	return views.Index().Render(
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
