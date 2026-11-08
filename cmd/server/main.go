package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"video-git/internal/auth"
	"video-git/internal/database"
	"video-git/internal/users"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// загрузка переменных окружения в процесс
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("загрузка .env: %v", err)
	}

	// получение необходимых переменных
	access := os.Getenv("ACCESS_TOKEN_TTL")
	accessTTL, err := time.ParseDuration(access)
	if err != nil {
		log.Fatalf("разбор ACCESS_TOKEN_TTL: %v", err)
	}

	refresh := os.Getenv("REFRESH_TOKEN_TTL")
	refreshTTL, err := time.ParseDuration(refresh)
	if err != nil {
		log.Fatalf("разбор REFRESH_TOKEN_TTL: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatalf("JWT_SECRET не задан или короче 32 символов")
	}

	cookieSecure := os.Getenv("COOKIE_SECURE")
	var cookieSecurebool bool
	if cookieSecure == "" {
		cookieSecurebool = true
	} else {
		cookieSecurebool, err = strconv.ParseBool(cookieSecure)
		if err != nil {
			log.Fatalf("разбор COOKIE_SECURE: %v", err)
		}
	}

	// создание пула
	ctx := context.Background()
	pool, err := database.NewPool(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("создание пула: %v", err)
	}
	defer pool.Close()

	// создание объектов
	tManager := auth.NewTokenManager(jwtSecret, accessTTL)
	tRepo := auth.NewTokenRepository(pool)
	tService := auth.NewTokenService(tManager, tRepo, refreshTTL)
	tHandler := auth.NewTokenHandler(tService, refreshTTL, cookieSecurebool)

	userRepo := users.NewPostgresRepository(pool)
	userService := users.NewService(userRepo, tService)
	userHandler := users.NewUserHandler(userService, refreshTTL, cookieSecurebool)

	// создание роутера и регистрация маршрутов
	e := echo.New()
	e.POST("/auth/register", userHandler.CreateUser)
	e.POST("/auth/login", userHandler.LoginUser)
	e.POST("/auth/refresh", tHandler.RefreshTokens)

	api := e.Group("/api", tManager.Middleware())
	api.GET("/me", func(c echo.Context) error {
		id, ok := auth.UserIDFromContext(c)
		if !ok {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, echo.Map{"user_id": id})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
