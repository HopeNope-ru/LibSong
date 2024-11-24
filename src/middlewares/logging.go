package middlewares

import "net/http"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Релизация логгирования с уровнем DEBUG
		// v := mux.Vars(r)
		// Реализаци логгирования с уровнем INFO
		next.ServeHTTP(w, r)
	})
}
