package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)


type Repository interface {          
    Create(ctx context.Context, user User) (User, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: pool,
	}
}

func (pr *PostgresRepository) Create(ctx context.Context, user User) (User, error) {
	var newUser User
	err := pr.db.QueryRow(ctx, "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at", user.Username, user.Email, user.PasswordHash).Scan(&newUser.ID, &newUser.CreatedAt)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return User{}, ErrAlreadyExists
		}
		return User{}, err
	}

	newUser.Username = user.Username
	newUser.Email = user.Email
	newUser.PasswordHash = user.PasswordHash
	return newUser, nil
}

