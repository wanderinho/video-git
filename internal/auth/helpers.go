package auth

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SetRefreshCookie(c echo.Context, refreshToken string, refreshTTL time.Duration, cookieSecure bool) {
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/auth/refresh",
		MaxAge:   int(refreshTTL.Seconds()),
		HttpOnly: true,
		Secure:   cookieSecure, // локально; в продакшне true
		SameSite: http.SameSiteStrictMode,
	})
}