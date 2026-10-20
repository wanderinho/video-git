package auth

import (
	"context"
	"time"
)

type TokenService struct {
	tm *TokenManager
	tr *TokenRepository
	refreshTTL time.Duration
}

func NewTokenService(tm *TokenManager, tr *TokenRepository, refreshTTL time.Duration) *TokenService {
	return &TokenService{
		tm: tm,
		tr: tr,
		refreshTTL: refreshTTL,
	}
}

func (ts *TokenService) GetTokenPair(ctx context.Context, userID int64) (string, string, error) {
	access, err := ts.tm.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refresh, refresh_hash, err := GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	now := time.Now()
	exp := now.Add(ts.refreshTTL)
	if err = ts.tr.CreateToken(ctx, userID, refresh_hash, exp); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}


func (ts *TokenService) UpdateTokenPair(ctx context.Context, refresh string) (string, string, error) {
	refresh_hash := HashRefreshToken(refresh)
	userID, err := ts.tr.PopHash(ctx, refresh_hash)
	if err != nil {
		return "", "", err
	}

	newAccess, newRefresh, err := ts.GetTokenPair(ctx, userID)
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}