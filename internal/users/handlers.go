package users

import (
	"errors"
	"net/http"
	"time"
	"video-git/internal/auth"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service      *UserService
	refreshTTL   time.Duration
	cookieSecure bool
}

func NewUserHandler(service *UserService, refreshTTL time.Duration, cookieSecure bool) *UserHandler {
	return &UserHandler{
		service:      service,
		refreshTTL:   refreshTTL,
		cookieSecure: cookieSecure,
	}
}

func (uh *UserHandler) CreateUser(c echo.Context) error {
	var registerReq UserRegisterRequest
	if err := c.Bind(&registerReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request body"})
	}

	u, err := uh.service.Register(c.Request().Context(), registerReq.Username, registerReq.Email, registerReq.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrLongPassword):
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "password too long"})

		case errors.Is(err, ErrAlreadyExists):
			return c.JSON(http.StatusConflict, echo.Map{"message": "username or email already exists"})

		default:
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
		}
	}

	return c.JSON(http.StatusCreated, UserRegisterResponse{ID: u.id, Username: u.username, Email: u.email, CreatedAt: u.createdAt})
}

func (uh *UserHandler) LoginUser(c echo.Context) error {
	var loginReq UserLoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request body"})
	}

	u, err := uh.service.Login(c.Request().Context(), loginReq.Username, loginReq.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "incorrect username or password"})

		case errors.Is(err, ErrInvalidCredentials):
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "incorrect username or password"})

		default:
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
		}
	}

	auth.SetRefreshCookie(c, u.refreshToken, uh.refreshTTL, uh.cookieSecure)

	return c.JSON(http.StatusOK, UserLoginResponse{ID: u.id, Username: u.username, Email: u.email, CreatedAt: u.createdAt, AccessToken: u.accessToken})
}
