package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/titan-commerce/backend/pkg/security"
)

// AuthMiddleware creates JWT authentication middleware
func AuthMiddleware(jwt *security.JWTService, excludePaths ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check excluded paths
			for _, path := range excludePaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Extract token from Authorization header
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "Authorization required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(auth, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			// Verify token
			claims, err := jwt.Verify(parts[1])
			if err != nil {
				if err == security.ErrTokenExpired {
					http.Error(w, "Token expired", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add claims to context
			ctx := SetClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RoleMiddleware checks for required roles
func RoleMiddleware(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r.Context())
			if claims == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			hasRole := false
			for _, role := range requiredRoles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ValidateMiddleware validates request body
func ValidateMiddleware[T any](validate func(*T) *security.ValidationResult) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			var body T
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			if validate != nil {
				result := validate(&body)
				if !result.Valid {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"errors": result.Errors,
					})
					return
				}
			}

			// Store validated body in context
			ctx := SetRequestBody(r.Context(), &body)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TimeoutMiddleware adds request timeout
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, timeout, "Request timeout")
	}
}
