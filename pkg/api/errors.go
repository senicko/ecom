package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

type HttpError struct {
	Status int
	Msg    string
	Err    error
}

func (e HttpError) Error() string {
	if e.Err == nil {
		return e.Msg
	}
	return e.Err.Error()
}

// JsonError sends json response back to the user.
func JsonError(w http.ResponseWriter, l *zap.Logger, status int, msg string) {
	res, _ := json.Marshal(map[string]string{
		"message": msg,
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(res); err != nil {
		l.Error("failed to write response message", zap.Error(err))
	}
}

// HandleError sends an error response to the client.
func HandleError(w http.ResponseWriter, l *zap.Logger, err error) {
	l.Error(err.Error(), zap.Error(err))

	var e HttpError
	if errors.As(err, &e) {
		JsonError(w, l, e.Status, e.Msg)
	} else {
		s := http.StatusInternalServerError
		JsonError(w, l, s, http.StatusText(s))
	}
}
