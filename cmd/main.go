package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"shp/pkg/config"
	"shp/users"
	"time"

	"github.com/go-chi/chi/v5"
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
func NewServer(cfg *config.AppConfig, db *pgxpool.Pool, log *zap.Logger) *http.Server {
	mux := chi.NewMux()

	userRepo := users.NewRepo(db, log)
	userSrv := users.NewSvc(userRepo, log)
	userController := users.NewController(userSrv, log)
	userController.SetupRoutes(mux)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      mux,
	}
}

func main() {
	cfg := config.FromFlags()

	l, err := NewLogger(cfg)
	if err != nil {
		log.Panic("failed to create structured logger")
	}
	defer l.Sync()

	db, err := InitDatabase(cfg)
	if err != nil {
		l.Error("couldn't initialize database connection")
	}
	defer db.Close()

	s := NewServer(cfg, db, l)

	l.Info("server starting ...",
		zap.Int("port", cfg.Port),
		zap.String("database_url", cfg.DatabaseURL),
	)

	if err := s.ListenAndServe(); err != nil {
		log.Panic("start failed", zap.Error(err))
	}
}
