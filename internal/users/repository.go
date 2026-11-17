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

func (ur *UserRepository) Create(ctx context.Context, us user) (user, error) {
	err := ur.db.QueryRow(ctx, "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at", us.username, us.email, us.passwordHash).Scan(&us.id, &us.createdAt)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return user{}, ErrAlreadyExists
		}
		return user{}, err
	}

	return us, nil
}

func (ur *UserRepository) FindByUsername(ctx context.Context, username string) (user, error) {
	var us user

	err := ur.db.QueryRow(ctx, "SELECT * FROM users WHERE username = $1", username).Scan(&us.id, &us.username, &us.email, &us.passwordHash, &us.createdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user{}, ErrInvalidCredentials
		}
		return user{}, err
	}

	return us, nil
}
