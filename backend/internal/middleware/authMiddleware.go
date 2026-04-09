package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/YahyaMudallal/newsWebSite/internal/auth"
)

// AuthMiddleware is a middleware that checks if the user is authenticated used for protected queries.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		authHeader := r.Header.Get("Authorization")
	
		// check if the token is present
		if authHeader == "" {
			http.Error(w, "Token required", http.StatusUnauthorized)
			return
		} 	

		// get the token string and remove "Bearer" if present
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// verify the JWT token
		userID, err := auth.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// store userID in the context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next(w, r.WithContext(ctx))
	})
}