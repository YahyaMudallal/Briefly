package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/YahyaMudallal/newsWebSite/internal/auth"
)

// OptionalAuthMiddleware checks for a JWT token.
// If valid, it adds the user ID to the context. If missing, it proceeds as a guest.
func OptionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// If no token is provided, just proceed to the handler as a guest
		if authHeader == "" {
			next(w, r)
			return
		}

		// Check if the header format is "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Bad format, proceed as guest
			next(w, r)
			return
		}
		// or tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		tokenString := parts[1]

		// Parse and validate the token (using your existing auth package)
		userID, err := auth.VerifyToken(tokenString)
		if err != nil {
			// Invalid/expired token. Proceed as guest so they still see the feed.
			next(w, r)
			return
		}

		// Token is valid! Inject the user_id into the context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		reqWithCtx := r.WithContext(ctx)

		// Pass the new context to the handler
		next(w, reqWithCtx)
	}
}
