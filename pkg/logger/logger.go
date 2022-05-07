package logger

import (
	"shp/pkg/config"

	"go.uber.org/zap"
)

// NewLogger creates new structured logger.
func New(cfg *config.AppConfig) (*zap.Logger, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return l, nil
}
