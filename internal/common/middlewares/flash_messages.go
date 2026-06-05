package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type (
	contextKey      string
	FlashMessageKey int
)

var flashMessagesContextKey contextKey = "flash-messages-context-key"

const (
	FlashMessageKeyError FlashMessageKey = iota
)

type FlashMessageData struct {
	Error []string
}

type FlashMessagesProvider struct{}

func (fm FlashMessagesProvider) Flash(c fiber.Ctx, key FlashMessageKey, value string) {
	sess := session.FromContext(c)

	switch key {
	case FlashMessageKeyError:
		{
			var flashes []string

			if val := sess.Get("flash:error"); val != nil {
				flashes = val.([]string)
				flashes = append(flashes, value)
				sess.Set("flash:error", flashes)
			} else {
				flashes = []string{value}
				sess.Set("flash:error", flashes)
			}
		}
	}
}

func (fm FlashMessagesProvider) Flashes(c fiber.Ctx) *FlashMessageData {
	sess := session.FromContext(c)

	res := FlashMessageData{}

	{
		if flashes := sess.Get("flash:error"); flashes != nil {
			sess.Delete("flash:error")
			res.Error = flashes.([]string)
		}
	}

	if len(res.Error) <= 0 {
		return nil
	}

	return &res
}

func FlashMessageProviderFromContext(c fiber.Ctx) *FlashMessagesProvider {
	if m, ok := fiber.ValueFromContext[*FlashMessagesProvider](c, flashMessagesContextKey); ok {
		return m
	}

	return nil
}

func FlashMessages() fiber.Handler {
	fmp := FlashMessagesProvider{}

	return func(c fiber.Ctx) error {
		fiber.StoreInContext(c, flashMessagesContextKey, &fmp)

		return c.Next()
	}
}
