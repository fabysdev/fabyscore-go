package middleware

import (
	"context"
	"net/http"
	"time"
)

// Timeout creates a TimeoutHandler and updates the request context with the given timeout.
func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.TimeoutHandler(http.HandlerFunc(fn), timeout, "Request Timeout")
	}
}
