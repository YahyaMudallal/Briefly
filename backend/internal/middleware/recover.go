package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// RecoverMiddleware is a middleware function that recovers from panics in the handlers and returns a 500 Internal Server Error response to the client.
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// logging the panic and the stack trace for debugging purposes
				log.Printf("[PANIC RECOVERED] %v\n%s", err, debug.Stack())

				// returning a 500 Internal Server Error response to the client
				http.Error(w, "Internal Server Error ", http.StatusInternalServerError)
			}
		}()

		// calling the next handler in the chain
		next.ServeHTTP(w, r)
	})
}