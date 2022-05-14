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

// AddCookie adds a new HTTP only cookie to request.
func AddCookie(w http.ResponseWriter, name, value string) {
	// TODO: Set expire date
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
	})
}
