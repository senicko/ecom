package user

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"shp/pkg/api"
	"shp/pkg/models"
)

var (
	ErrEmailTaken = errors.New("")
)

type Service interface {
	// SignIn SingIn sings in a user. In case email provided by the user is already being used returns an error.
	SignIn(ctx context.Context, params *models.UserCreateParams) (*models.User, error)
}

type service struct {
	userRepo Repo
	l        *zap.Logger
}

// NewService creates a new user service.
func NewService(userRepo Repo, l *zap.Logger) *service {
	return &service{
		userRepo: userRepo,
		l:        l,
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

func (s service) SignIn(ctx context.Context, params *models.UserCreateParams) (*models.User, error) {
	if u, err := s.userRepo.FindByEmail(ctx, params.Email); err != nil {
		return nil, err
	} else if u != nil {
		return nil, api.HttpError{Status: http.StatusBadRequest, Msg: "email is already taken", Err: err}
	}

	h, err := hashPassword(params.Password)
	if err != nil {
		return nil, err
	}
	params.Password = h

	u, err := s.userRepo.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	return u, nil
}
