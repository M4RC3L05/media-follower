package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func Auth(c fiber.Ctx) error {
	sess := session.FromContext(c)
	fmp := FlashMessageProviderFromContext(c)

	if sess.Get("userId") == nil {
		fmp.Flash(c, FlashMessageKeyError, "You MUST be authenticated to access")
		return c.Redirect().Route("auth-login")
	}

	return c.Next()
}

func NotAuth(c fiber.Ctx) error {
	sess := session.FromContext(c)
	fmp := FlashMessageProviderFromContext(c)

	if sess.Get("userId") != nil {
		fmp.Flash(c, FlashMessageKeyError, "You MUST NOT be authenticated to access")
		return c.Redirect().Route("home")
	}

	return c.Next()
}
