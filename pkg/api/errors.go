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
	return e.Err.Error()
}

// JsonResponse sends json response back to the user.
func JsonError(w http.ResponseWriter, status int, msg string) {
	res, _ := json.Marshal(map[string]string{
		"message": msg,
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}

// HandleError sends an error response to the client.
func HandleError(w http.ResponseWriter, log *zap.Logger, err error) {
	log.Error("error occured", zap.Error(err))

	var e HttpError
	if errors.As(err, &e) {
		JsonError(w, e.Status, e.Msg)
	} else {
		s := http.StatusInternalServerError
		JsonError(w, s, http.StatusText(s))
	}
}
