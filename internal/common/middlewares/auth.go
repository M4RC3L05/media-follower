package middlewares

import "github.com/labstack/echo/v5"

func Authenticated() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			sess := SessionProviderFromContext(c)
			flash := FlashMessageProviderFromContext(c)

			if !sess.IsAuthenticated(c) {
				if err := flash.Flash(
					c,
					FlashMessageKeyError,
					"You MUST be logged in to proceed",
				); err != nil {
					return err
				}

				return c.Redirect(302, "/auth/login")
			}

			return next(c)
		}
	}
}

func NotAuthenticated() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			sess := SessionProviderFromContext(c)
			flash := FlashMessageProviderFromContext(c)

			if sess.IsAuthenticated(c) {
				if err := flash.Flash(
					c,
					FlashMessageKeyError,
					"You MUST NOT be logged in to proceed",
				); err != nil {
					return err
				}

				referer := c.Request().Header.Get("Referer")
				if referer == "" {
					referer = "/"
				}

				return c.Redirect(302, referer)
			}

			return next(c)
		}
	}
}
