package users

import "time"

type user struct {
	id           int64
	username     string
	email        string
	passwordHash string
	createdAt    time.Time
}

type UserWithTokens struct {
	user
	accessToken  string
	refreshToken string
}
