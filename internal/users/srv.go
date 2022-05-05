package users

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken = errors.New("email is already taken")
)

type Srv interface {
	// SingIn sings in a user. In case email provided by the user is already being used returns an error.
	SignIn(ctx context.Context, params *UserCreateParams) (*User, error)
}

type srv struct {
	repo Repo
	log  *zap.Logger
}

// NewSrv creates a new user service.
func NewSrv(repo Repo, log *zap.Logger) *srv {
	return &srv{
		repo: repo,
		log:  log,
	}
}

// hashPassword hashes user's password.
func hashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing failed: %w", err)
	}

	return string(h), nil
}

func (s *srv) SignIn(ctx context.Context, params *UserCreateParams) (*User, error) {
	if u, err := s.repo.FindByEmail(ctx, params.Email); err != nil {
		return nil, err
	} else if u != nil {
		return nil, ErrEmailTaken
	}

	h, err := hashPassword(params.Password)
	if err != nil {
		return nil, err
	}
	params.Password = h

	u, err := s.repo.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	return u, nil
}
