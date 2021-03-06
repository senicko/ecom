package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"shp/pkg/api"
	"shp/pkg/auth"
	"shp/pkg/models"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type controller struct {
	l           *zap.Logger
	authService auth.Service
	userService Service
}

// NewController creates a new users' controller.
func NewController(l *zap.Logger, userService Service, authService auth.Service) *controller {
	return &controller{
		userService: userService,
		authService: authService,
		l:           l,
	}
}

// SetupRoutes registers all users' router routes.
func (c *controller) SetupRoutes(m *chi.Mux) {
	m.Post("/signup", c.SignUp)
}

// SignUp is a http controller that creates a new account.
func (c *controller) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var params *models.UserCreateParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		api.HandleError(w, c.l, &api.HttpError{
			Status: http.StatusBadRequest,
			Msg:    http.StatusText(http.StatusBadRequest),
			Err:    err,
		})
		return
	}

	user, err := c.userService.SignIn(ctx, params)
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			api.HandleError(w, c.l, &api.HttpError{
				Status: http.StatusBadRequest,
				Msg:    err.Error(),
			})
			return
		}

		api.HandleError(w, c.l, err)
		return
	}

	atk, err := c.authService.NewAccessToken(user)
	if err != nil {
		api.HandleError(w, c.l, err)
		return
	}

	rtk, err := c.authService.NewRefreshToken(user)
	if err != nil {
		api.HandleError(w, c.l, err)
		return
	}

	api.AddCookie(w, "atk", atk)
	api.AddCookie(w, "rtk", rtk)

	res, err := json.Marshal(user)
	if err != nil {
		api.HandleError(w, c.l, err)
	}

	if err := api.JsonResponse(w, http.StatusOK, res); err != nil {
		c.l.Error("can't respond", zap.Error(err))
	}
}
