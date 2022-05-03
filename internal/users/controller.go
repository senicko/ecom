package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type controller struct {
	srv Srv
	log *zap.Logger
}

// NewController creates a new user controller.
func NewController(srv Srv, log *zap.Logger) *controller {
	return &controller{
		srv: srv,
		log: log,
	}
}

// SetupRoutes registers all user router routes.
func (c *controller) SetupRoutes(m *chi.Mux) {
	m.Post("/signin", c.SignIn)
}

// SingIn is a http controller that creates a new account.
func (c *controller) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var params *UserCreateParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u, err := c.srv.SignIn(ctx, params)
	if errors.Is(err, ErrEmailTaken) {
		http.Error(w, "This email is already taken", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(u)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
