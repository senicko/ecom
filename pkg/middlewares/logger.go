package middlewares

import (
	"net/http"

	"go.uber.org/zap"
)

// LoggerMiddleware logs info about incoming requests.
func LoggerMiddleware(l *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Info("request", zap.String("method", r.Method), zap.String("url", r.URL.String()))
			next.ServeHTTP(w, r)
		})
	}
}
