package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: pool,
	}
}

func (ur *UserRepository) Create(ctx context.Context, user User) (User, error) {
	err := ur.db.QueryRow(ctx, "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at", user.Username, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return User{}, ErrAlreadyExists
		}
		return User{}, err
	}

	return user, nil
}

func (ur *UserRepository) FindByUsername(ctx context.Context, username string) (User, error) {
	var user User

	err := ur.db.QueryRow(ctx, "SELECT * FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrInvalidCredentials
		}
		return User{}, err
	}

	return user, nil
}
