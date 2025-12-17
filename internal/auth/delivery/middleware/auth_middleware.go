package middleware

import (
	"context"
	"net/http"
)

type UserIDKey string

const ContextUserIDKey UserIDKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("auth_token")
		if err != nil {

			http.Error(w, "Unauthorized: Missing or invalid token in cookie", http.StatusUnauthorized)
			return
		}

		tokenValue := cookie.Value

		if !isTokenValid(tokenValue) {
			http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
			return
		}

		mockUserID := extractUserIDFromMockToken(tokenValue)

		if mockUserID == "" {
			http.Error(w, "Unauthorized: Cannot determine user identity", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserIDKey, mockUserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isTokenValid(token string) bool {

	return token != "" && len(token) > 10
}

func extractUserIDFromMockToken(token string) string {

	return "a0000000-0000-0000-0000-000000000001"
}
