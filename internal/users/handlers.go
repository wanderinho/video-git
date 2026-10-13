package users

import (
	"errors"
	"net/http"
	"video-git/internal/users/dto"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service *Service
}

func NewUserHandler(service *Service) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) CreateUser(c echo.Context) error {
	var registerReq dto.UserRegisterRequest
	if err := c.Bind(&registerReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request body"})
	}

	u, err := uh.service.Register(c.Request().Context(), registerReq.Username, registerReq.Email, registerReq.Password)
	if err != nil {
		if errors.Is(err, ErrLongPassword) {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "password too long"})
		} else if errors.Is(err, ErrAlreadyExists) {
			return c.JSON(http.StatusConflict, echo.Map{"message": "username or email already exists"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
	}

	return c.JSON(http.StatusCreated, dto.UserRegisterResponse{ID: u.ID, Username: u.Username, Email: u.Email, CreatedAt: u.CreatedAt})
}

func (uh *UserHandler) LoginUser(c echo.Context) error {
	var loginReq dto.UserLoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request body"})
	}

	u, err := uh.service.Login(c.Request().Context(), loginReq.Username, loginReq.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "incorrect username or password"})
		} else if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "incorrect username or password"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
	}

	return c.JSON(http.StatusOK, dto.UserLoginResponse{ID: u.ID, Username: u.Username, Email: u.Email, CreatedAt: u.CreatedAt})
}