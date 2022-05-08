package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"shp/pkg/api"
	"shp/pkg/config"
	"shp/pkg/logger"
	"shp/users"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// InitDatabase initializes pool connection with database.
func InitDatabase(cfg *config.AppConfig) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewServer creates a new http server.
func NewServer(cfg *config.AppConfig, db *pgxpool.Pool, l *zap.Logger) *http.Server {
	mux := chi.NewMux()

	// setup global middlewares
	mux.Use(api.LoggerMiddleware(l))
	mux.Use(middleware.Recoverer)

	// setup controllers
	userRepo := users.NewRepo(db, l)
	userSrv := users.NewSvc(userRepo, l)
	userController := users.NewController(userSrv, l)
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

	l, err := logger.New(cfg)
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	defer func(l *zap.Logger) {
		err := l.Sync()
		if err != nil {
			l.Error("can't sync the logger", zap.Error(err))
		}
	}(l)

	db, err := InitDatabase(cfg)
	if err != nil {
		l.Error("can't initialize database connection")
	}
	defer db.Close()

	s := NewServer(cfg, db, l)

	l.Info("server starting ...",
		zap.Int("port", cfg.Port),
		zap.String("database_url", cfg.DatabaseURL),
	)

	if err := s.ListenAndServe(); err != nil {
		l.Panic("start failed", zap.Error(err))
	}
}
