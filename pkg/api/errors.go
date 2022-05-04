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
func JsonError(rw http.ResponseWriter, status int, msg string) {
	res, _ := json.Marshal(map[string]string{
		"message": msg,
	})

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(res)
}

// HandleError sends an error response to the client.
func HandleError(rw http.ResponseWriter, log *zap.Logger, err error) {
	log.Error("error occured", zap.Error(err))

	var e HttpError
	if errors.As(err, &e) {
		JsonError(rw, e.Status, e.Msg)
	} else {
		s := http.StatusInternalServerError
		JsonError(rw, s, http.StatusText(s))
	}
}
