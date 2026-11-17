package users

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo  Repository
	token TokenIssuer
}

type Repository interface {
	Create(ctx context.Context, us user) (user, error)
	FindByUsername(ctx context.Context, username string) (user, error)
}

type TokenIssuer interface {
	GetTokenPair(ctx context.Context, userID int64) (string, string, error)
	UpdateTokenPair(ctx context.Context, refresh string) (string, string, error)
}

func NewUserService(r Repository, ti TokenIssuer) *UserService {
	return &UserService{
		repo:  r,
		token: ti,
	}
}

func (s *UserService) Register(ctx context.Context, username, email, password string) (user, error) {
	hash, err := hashPassword(password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return user{}, ErrLongPassword
		}
		return user{}, err
	}

	newUser := user{username: username, email: email, passwordHash: hash}
	newUser, err = s.repo.Create(ctx, newUser)
	if err != nil {
		return user{}, err
	}

	return newUser, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (UserWithTokens, error) {
	u, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return UserWithTokens{}, err
	}

	if err := checkPassword(u.passwordHash, password); err != nil {
		return UserWithTokens{}, err
	}

	access, refresh, err := s.token.GetTokenPair(ctx, u.id)
	if err != nil {
		return UserWithTokens{}, err
	}

	return UserWithTokens{u, access, refresh}, nil
}
