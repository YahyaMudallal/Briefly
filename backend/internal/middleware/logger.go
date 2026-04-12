package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware is a middleware function that logs the incoming requests and the time taken to process them.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("[%s] %s - %s (treatment time: %v)", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}
