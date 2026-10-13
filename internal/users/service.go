package users

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

type Repository interface {          
    Create(ctx context.Context, user User) (User, error)
    FindByUsername(ctx context.Context, username string) (User, error)
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Register(ctx context.Context, username, email, password string) (User, error) {
	hash, err := hashPassword(password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return User{}, ErrLongPassword
		}
		return User{}, err
	}

	newUser := User{Username: username, Email: email, PasswordHash: hash}
	newUser, err = s.repo.Create(ctx, newUser)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (s *Service) Login(ctx context.Context, username, password string) (User, error) {
	u, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return User{}, err
	}

	if err := checkPassword(u.PasswordHash, password); err != nil {
		return User{}, err
	}

	return u, nil
}