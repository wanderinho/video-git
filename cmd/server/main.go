package main

import (
	"context"
	"os"
	"video-git/internal/database"
	"video-git/internal/users"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}

	ctx := context.Background()
	pool, err := database.NewPool(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		pool.Close()
		return
	}
	defer pool.Close()

	pgRepo := users.NewPostgresRepository(pool)
	service := users.NewService(pgRepo)
	handler := users.NewUserHandler(service)

	e := echo.New()
	e.POST("auth/register/", handler.CreateUser)
}