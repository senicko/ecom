package auth

import (
	"fmt"
	"shp/users"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"
)

type Svc interface {
}

type svc struct {
	l *zap.Logger
}

// New creates a new logger service.
func New(l *zap.Logger) *svc {
	return &svc{
		l: l,
	}
}

// CreateTokenPackage creates access token and refresh token.
// TODO: Maybe this could be split into two separate functions.
func (s svc) CreateTokenPackage(user users.User) (*jwt.Token, *jwt.Token, error) {
	refreshToken, err := jwt.NewBuilder().
		IssuedAt(time.Now()).
		Expiration(time.Now().AddDate(0, 1, 0)).
		Build()

	if err != nil {
		return nil, nil, fmt.Errorf("can't build refresh token: %w", err)
	}

	accessToken, err := jwt.NewBuilder().
		Issuer("localhost").
		IssuedAt(time.Now()).
		Build()

	if err != nil {
		return nil, nil, fmt.Errorf("can't build access token: %w", err)
	}

	return &accessToken, &refreshToken, nil
}
