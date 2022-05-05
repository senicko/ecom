package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"shp/pkg/api"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type controller struct {
	srv Srv
	log *zap.Logger
}

// NewController creates a new users controller.
func NewController(srv Srv, log *zap.Logger) *controller {
	return &controller{
		srv: srv,
		log: log,
	}
}

// SetupRoutes registers all users router routes.
func (c *controller) SetupRoutes(m *chi.Mux) {
	m.Post("/signin", c.SignIn)
}

// SingIn is a http controller that creates a new account.
func (c *controller) SignIn(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var params *UserCreateParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		api.HandleError(rw, c.log, err)
		return
	}

	u, err := c.srv.SignIn(ctx, params)
	if errors.Is(err, ErrEmailTaken) {
		api.HandleError(rw, c.log, api.HttpError{
			Status: http.StatusBadRequest,
			Msg:    "Email adress is already taken",
			Err:    err,
		})
		return
	} else if err != nil {
		api.HandleError(rw, c.log, err)
		return
	}

	res, err := json.Marshal(u)
	if err != nil {
		api.HandleError(rw, c.log, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(res)
}
