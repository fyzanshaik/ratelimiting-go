package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fyzanshaik/ratelimiting-go/middleware"
	"github.com/fyzanshaik/ratelimiting-go/ratelimiter"
)

func main() {
	limiter := ratelimiter.NewRateLimiter(10, 2)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Hello, World!"}`)
	})

	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"data": [1, 2, 3, 4, 5]}`)
	})

	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		clientIP := middleware.GetClientIP(r)
		capacity, rate, tokens := limiter.GetMetrics(clientIP)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"capacity": %.2f, "rate": %.2f, "available_tokens": %.2f}`, capacity, rate, tokens)
	})

	handler := middleware.RateLimit(limiter)(mux)

	fmt.Println("Server starting on :8080")
	fmt.Println("Rate limit: 10 requests capacity, 2 requests per second refill")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
