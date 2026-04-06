package middleware

import (
	"context"
	"net/http"

	"github.com/YahyaMudallal/newsWebSite/internal/auth"
)

// AuthMiddleware is a middleware that checks if the user is authenticated used for protected queries.
func AuthMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		tokenString := r.Header.Get("Authorization")
	
		// verify the JWT token
		userID, err := auth.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Unothonrized", http.StatusUnauthorized)
		}

		// store userID in the context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next(w, r.WithContext(ctx))
	})
}