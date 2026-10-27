package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// ключ, под которым userID лежит в контексте запроса
const userIDKey = "user_id"

func (tm *TokenManager) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			token, found := strings.CutPrefix(header, "Bearer ")
			if !found || token == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthorized"})
			}

			userID, err := tm.ParseAccessToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthorized"})
			}

			c.Set(userIDKey, userID)
			return next(c)
		}
	}
}

func UserIDFromContext(c echo.Context) (int64, bool) {
	id, ok := c.Get(userIDKey).(int64)
	return id, ok
}