package middlewares

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
)

type SessionsProvider struct {
	Store sessions.Store
}

func defaultSessionCookieOptions() *sessions.Options {
	return &sessions.Options{
		Path:        "/",
		MaxAge:      60 * 60 * 24,
		Secure:      true,
		HttpOnly:    true,
		Partitioned: true,
		SameSite:    http.SameSiteStrictMode,
	}
}

func (sp SessionsProvider) IsAuthenticated(c *echo.Context) bool {
	sess, err := sp.Store.Get(c.Request(), "_sess")
	return sess != nil && err == nil && sess.Values["uid"] != nil
}

func (sp SessionsProvider) Authenticate(c *echo.Context, uId string) error {
	sess, _ := sp.Store.New(c.Request(), "_sess")
	sess.Values["uid"] = uId
	sess.Options = defaultSessionCookieOptions()

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}

func (sp SessionsProvider) Logout(c *echo.Context) error {
	sess, _ := sp.Store.Get(c.Request(), "_sess")
	sess.Options = defaultSessionCookieOptions()
	sess.Options.MaxAge = -1

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}

func SessionProviderFromContext(c *echo.Context) SessionsProvider {
	m, _ := echo.ContextGet[SessionsProvider](c, "_sess")
	return m
}

func Session(store sessions.Store) echo.MiddlewareFunc {
	provider := SessionsProvider{Store: store}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("_session_store", store)
			c.Set("_sess", provider)

			return next(c)
		}
	}
}
