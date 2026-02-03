package middleware

import (
	"context"
	"fmt"
	"lms_backend/internal/domain"
	"net/http"
	"strings"
)

type UserContextData struct {
    UserID string
    Role   domain.Role
}

const ContextUserDataKey = "userData"

func extractUserDataFromToken(tokenValue string) *UserContextData {
   
    if tokenValue == "mock-jwt-token-for-user-admin@capedu.kz" || strings.Contains(tokenValue, "admin") {
        return &UserContextData{
            UserID: "00000000-0000-0000-0000-000000000001", 
            Role: domain.RoleAdmin,
        }
    }
    if strings.Contains(tokenValue, "student") {
        return &UserContextData{
            UserID: "a0000000-0000-0000-0000-000000000001", 
            Role: domain.RoleStudent,
        }
    }
    
    return nil 
}


func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, "Unauthorized: Missing or invalid token in cookie", http.StatusUnauthorized)
			return
		}
		tokenValue := cookie.Value

		userData := extractUserDataFromToken(tokenValue)

		if userData == nil {
			http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
			return
		}
        
        // ПРОКИДЫВАЕМ В CONTEXT
		ctx := context.WithValue(r.Context(), ContextUserDataKey, userData)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleRequiredMiddleware(allowedRoles ...domain.Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
			
			userCtxData, ok := r.Context().Value(ContextUserDataKey).(*UserContextData)
            
			if !ok || userCtxData == nil || userCtxData.UserID == "" {
				http.Error(w, "Forbidden: Authentication context not set.", http.StatusForbidden)
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
				http.Error(w, "Forbidden: Access denied. Required roles: "+strings.Join(strings.Fields(strings.Trim(fmt.Sprintf("%s", allowedRoles), "[]")), ", "), http.StatusForbidden)
				return
			}
            
			next.ServeHTTP(w, r)
		})
	}
}