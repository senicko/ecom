package users_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"shp/internal/users"
	"testing"
)

type TestSrv struct{}

func (s TestSrv) SignIn(ctx context.Context, params *users.UserCreateParams) (*users.User, error) {
	return &users.User{
		ID:        1,
		Firstname: params.Firstname,
		Lastname:  params.Lastname,
		Email:     params.Email,
	}, nil
}

func assertNilError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("expected %v but got %v", nil, err)
	}
}

func assertDeeplyEqual(t testing.TB, a any, b any) {
	t.Helper()

	if !reflect.DeepEqual(a, b) {
		t.Errorf("expected %v, got %v", b, a)
	}
}

func TestPOSTSignIn(t *testing.T) {
	controller := users.NewController(TestSrv{}, nil)

	params := users.UserCreateParams{
		Firstname: "Tom",
		Lastname:  "Thompson",
		Email:     "tomthompson@email.com",
		Password:  "123",
	}

	t.Run("returns data of created user", func(t *testing.T) {
		expected := users.User{
			ID:        1,
			Firstname: params.Firstname,
			Lastname:  params.Lastname,
			Email:     params.Email,
		}

		b, _ := json.Marshal(params)
		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewReader(b))
		res := httptest.NewRecorder()
		controller.SignIn(res, req)

		got := res.Body.Bytes()
		var user users.User
		err := json.Unmarshal(got, &user)

		assertNilError(t, err)
		assertDeeplyEqual(t, user, expected)
	})
}
