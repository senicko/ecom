package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

type HttpError struct {
	Status int    `json:"-"`
	Msg    string `json:"message"`
	Err    error  `json:"-"`
}

func (e HttpError) Error() string {
	return e.Err.Error()
}

// HandleError sends an error response to the client.
func HandleError(rw http.ResponseWriter, log *zap.Logger, err error) {
	log.Error("error occured", zap.Error(err))

	var status int
	var res []byte

	var e HttpError
	if errors.As(err, &e) {
		status = e.Status
		res, _ = json.Marshal(e)
	} else {
		status = http.StatusInternalServerError
		res, _ = json.Marshal(map[string]string{
			"message": http.StatusText(status),
		})
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(res)
}
