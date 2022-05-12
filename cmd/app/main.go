package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"shp/pkg/auth"
	"shp/pkg/config"
	"shp/pkg/middlewares"
	"shp/pkg/user"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// NewLogger creates new structured logger.
func NewLogger(cfg *config.AppConfig) (*zap.Logger, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return l, nil
}

// InitDatabase initializes pool connection with database.
func InitDatabase(cfg *config.AppConfig) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewServer creates a new http server.
func NewServer(cfg *config.AppConfig, db *pgxpool.Pool, l *zap.Logger) (*http.Server, error) {
	mux := chi.NewMux()

	// setup global middlewares
	mux.Use(middlewares.LoggerMiddleware(l))
	mux.Use(middleware.Recoverer)

	authService := auth.NewService(l)

	// setup controllers
	userRepo := user.NewRepo(db, l)
	userService := user.NewService(userRepo, l)
	userController := user.NewController(l, userService, authService)
	userController.SetupRoutes(mux)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      mux,
	}, nil
}

func main() {
	cfg := config.FromFlags()

	l, err := NewLogger(cfg)
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	defer func(l *zap.Logger) {
		if err := l.Sync(); err != nil {
			l.Error("can't sync the logger", zap.Error(err))
		}
	}(l)

	db, err := InitDatabase(cfg)
	if err != nil {
		l.Error("can't initialize database connection")
	}
	defer db.Close()

	s, err := NewServer(cfg, db, l)
	if err != nil {
		l.Fatal("can't create a server", zap.Error(err))
	}

	l.Info("server starting ...",
		zap.Int("port", cfg.Port),
		zap.String("database_url", cfg.DatabaseURL),
	)

	if err := s.ListenAndServe(); err != nil {
		l.Fatal("start failed", zap.Error(err))
	}
}
