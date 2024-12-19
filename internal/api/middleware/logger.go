package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			slog.Info(fmt.Sprintf("%s %s", r.Method, r.URL.RequestURI()))
			next.ServeHTTP(w, r)
		},
	)
}
