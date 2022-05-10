package auth

import (
	"fmt"
	"shp/pkg/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

type Service interface {
	// NewAccessToken creates a new access token.
	NewAccessToken(user *models.User) (string, error)
	// NewRefreshToken creates a new refresh token.
	NewRefreshToken(user *models.User) (string, error)
}

type authService struct {
	l          *zap.Logger
	accessKey  []byte
	refreshKey []byte
}

// NewService creates a new logger service.
func NewService(l *zap.Logger) *authService {
	return &authService{
		l: l,
		// TODO: Load signing keys from config
		accessKey:  []byte("test_access_key"),
		refreshKey: []byte("test_refresh_key"),
	}
}

// TODO: Refactor NewAccessToken and NewRefreshToken to reduce code repetition

func (s authService) NewAccessToken(user *models.User) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
		Subject:   strconv.Itoa(user.ID),
		Issuer:    "localhost",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.accessKey)
	if err != nil {
		return "", fmt.Errorf("can't sign the token: %w", err)
	}

	return signed, nil
}

func (s authService) NewRefreshToken(user *models.User) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().AddDate(0, 1, 0).Unix(),
		Subject:   strconv.Itoa(user.ID),
		Issuer:    "localhost",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.accessKey)
	if err != nil {
		return "", fmt.Errorf("can't sign the token: %w", err)
	}

	return signed, nil
}
