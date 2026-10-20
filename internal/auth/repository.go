package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{
		db: pool,
	}
}

func (tr *TokenRepository) CreateToken(ctx context.Context, userID int64, refresh_hash string, refreshTTL time.Time) error {
	_, err := tr.db.Exec(ctx, "INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)", userID, refresh_hash, refreshTTL)
	
	return err
}

func (tr *TokenRepository) PopHash(ctx context.Context, hash string) (int64, error) {
	var userID int64
	err := tr.db.QueryRow(ctx, "DELETE FROM refresh_tokens WHERE token_hash = $1 AND expires_at > now() RETURNING user_id", hash).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrTokenNotFound
		}
		return 0, err
	}

	return userID, nil
}