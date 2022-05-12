package api

import (
	"fmt"
	"net/http"
)

// JsonResponse sends a json response back to the user.
func JsonResponse(w http.ResponseWriter, status int, payload []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(payload); err != nil {
		return fmt.Errorf("can't write: %w", err)
	}

	return nil
}
