package middleware

import (
	"fmt"
	"net/http"

	"github.com/fyzanshaik/ratelimiting-go/ratelimiter"
)

func RateLimit(limiter *ratelimiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr

			if !limiter.Allow(clientIP) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprintf(w, `{"error": "rate limit exceeded"}`)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
