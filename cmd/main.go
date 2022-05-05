package main

import (
	"context"
	"net/http"
	"shp/internal/users"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// NewServer creates a new http server.
func NewServer(db *pgxpool.Pool, log *zap.Logger) *http.Server {
	userRepo := users.NewRepo(db, log)
	userSrv := users.NewSrv(userRepo, log)
	userController := users.NewController(userSrv, log)

	mux := chi.NewMux()
	userController.SetupRoutes(mux)

	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      mux,
	}

	return s
}

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	connString := "postgresql://localhost:5432/shp"
	db, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Panic("database connection failed", zap.Error(err))
	}
	defer db.Close()

	log.Info("starting",
		zap.String("port", ":8080"),
		zap.String("database_conn_url", connString),
	)

	s := NewServer(db, log)
	if err := s.ListenAndServe(); err != nil {
		log.Panic("start failed", zap.Error(err))
	}
}
