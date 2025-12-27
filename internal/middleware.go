package internal

import (
	"dos/logger"
	"net/http"
)

// FIXME this log message appears AFTER all actions are done
func LogMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		logger.L.Info("http request", "method", r.Method, "path", r.URL.Path)
	})
}
