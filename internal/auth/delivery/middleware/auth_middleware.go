package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"lms_backend/internal/domain"
)

type UserContextData struct {
	UserID string
	Role   domain.Role
}

type contextKey string

const ContextUserDataKey contextKey = "userData"

func extractUserDataFromToken(tokenValue string) *UserContextData {
	if tokenValue == "mock-jwt-token-for-user-admin@capedu.kz" || strings.Contains(tokenValue, "admin") {
		return &UserContextData{
			UserID: "00000000-0000-0000-0000-000000000001",
			Role:   domain.RoleAdmin,
		}
	}
	if strings.Contains(tokenValue, "student") {
		return &UserContextData{
			UserID: "a0000000-0000-0000-0000-000000000001",
			Role:   domain.RoleStudent,
		}
	}
	return nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenValue string

		cookie, err := r.Cookie("auth_token")
		if err == nil {
			tokenValue = cookie.Value
		}

		if tokenValue == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenValue = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenValue == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		userData := extractUserDataFromToken(tokenValue)
		if userData == nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserDataKey, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleRequiredMiddleware(allowedRoles ...domain.Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCtxData, ok := r.Context().Value(ContextUserDataKey).(*UserContextData)
			if !ok || userCtxData == nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			isAllowed := false
			for _, role := range allowedRoles {
				if userCtxData.Role == role {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				msg := fmt.Sprintf("Forbidden: access denied for role %s", userCtxData.Role)
				http.Error(w, msg, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
