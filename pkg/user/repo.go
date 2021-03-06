package user

import (
	"context"
	"shp/pkg/models"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Repo interface {
	// FindByEmail finds a user with the specified email.
	FindByEmail(ctx context.Context, e string) (*models.User, error)
	// Create inserts a new user to the database.
	Create(ctx context.Context, p *models.UserCreateParams) (*models.User, error)
}

type repo struct {
	db *pgxpool.Pool
	l  *zap.Logger
}

// NewRepo creates a new users' repository.
func NewRepo(db *pgxpool.Pool, l *zap.Logger) *repo {
	return &repo{
		db: db,
		l:  l,
	}
}

func (r repo) FindByEmail(ctx context.Context, e string) (*models.User, error) {
	q := "SELECT * FROM users WHERE email = $1"
	row := r.db.QueryRow(ctx, q, e)

	var u models.User
	if err := u.Scan(row); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, nil
		default:
			return nil, err
		}
	}

	return &u, nil
}

func (r repo) Create(ctx context.Context, p *models.UserCreateParams) (*models.User, error) {
	q := "INSERT INTO users (email, firstname, lastname, password) VALUES ($1, $2, $3, $4) RETURNING *"
	row := r.db.QueryRow(ctx, q, p.Email, p.Firstname, p.Lastname, p.Password)

	var u models.User
	if err := u.Scan(row); err != nil {
		return nil, err
	}

	return &u, nil
}
