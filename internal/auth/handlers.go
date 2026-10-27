package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)


type TokenHandler struct {
	ts *TokenService
	refreshTTL time.Duration
	cookieSecure bool
}

func NewTokenHandler(ts *TokenService, refreshTTL time.Duration, cookieSecure bool) *TokenHandler {
	return &TokenHandler{
		ts: ts,
		refreshTTL: refreshTTL,
		cookieSecure: cookieSecure,
	}
}

func (th *TokenHandler) RefreshTokens(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthorized"})
	}
	refreshToken := cookie.Value

	access, refresh, err := th.ts.UpdateTokenPair(c.Request().Context(), refreshToken)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthorized"})
		} else {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
		}
	}

	SetRefreshCookie(c, refresh, th.refreshTTL, th.cookieSecure)
	
	return c.JSON(http.StatusOK, echo.Map{"access_token": access})
}