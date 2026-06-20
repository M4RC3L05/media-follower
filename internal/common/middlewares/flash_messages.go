package middlewares

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
)

type (
	FlashMessageKey string
)

const (
	FlashMessageKeyError   string = "error"
	FlashMessageKeySuccess string = "success"
)

type FlashMessageData struct {
	Error   []any
	Success []any
}

type FlashMessagesProvider struct {
	Store sessions.Store
}

func (fm FlashMessagesProvider) Flash(c *echo.Context, key string, value string) error {
	sess, _ := fm.Store.Get(c.Request(), "_fmessages")
	sess.AddFlash(value, key)

	sess.Options = &sessions.Options{
		Path:        "/",
		MaxAge:      0,
		Secure:      true,
		HttpOnly:    true,
		Partitioned: true,
		SameSite:    http.SameSiteStrictMode,
	}

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}

func (fm FlashMessagesProvider) Flashes(c *echo.Context) FlashMessageData {
	sess, _ := fm.Store.Get(c.Request(), "_fmessages")

	res := FlashMessageData{
		Error:   sess.Flashes("error"),
		Success: sess.Flashes("success"),
	}

	sess.Options = &sessions.Options{
		Path:        "/",
		MaxAge:      0,
		Secure:      true,
		HttpOnly:    true,
		Partitioned: true,
		SameSite:    http.SameSiteStrictMode,
	}
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return res
	}

	return res
}

func FlashMessageProviderFromContext(c *echo.Context) FlashMessagesProvider {
	m, _ := echo.ContextGet[FlashMessagesProvider](c, "_flash_messages_provider")
	return m
}

func FlashMessages(store sessions.Store) echo.MiddlewareFunc {
	provider := FlashMessagesProvider{Store: store}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("_flash_messages_provider", provider)

			return next(c)
		}
	}
}
