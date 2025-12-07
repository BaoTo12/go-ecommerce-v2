package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/titan-commerce/backend/pkg/security"
	"github.com/titan-commerce/backend/pkg/server"
)

/*
EXAMPLE: Unified Server Usage

This demonstrates how ALL packages work together:
- JWT authentication
- CSRF protection
- Security headers (CSP, HSTS, etc.)
- Rate limiting
- Response compression
- Distributed tracing
- Metrics collection
- Circuit breaker for external calls
*/

func main() {
	// Get JWT secret from environment
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("development-secret-change-in-production")
	}

	// Create server with config
	config := server.DefaultConfig()
	config.JWTSecret = jwtSecret
	config.CORSOrigins = []string{"http://localhost:3000"}
	config.ServiceName = "example-service"
	config.Port = "8080"

	srv := server.New(config)

	// Register routes with full middleware stack automatically applied
	
	// Public endpoints (no auth required)
	srv.Handle("/api/v1/auth/login", handleLogin(srv))
	srv.Handle("/api/v1/auth/register", handleRegister())
	srv.Handle("/api/v1/products", handleProducts())

	// Protected endpoints (require auth)
	srv.Handle("/api/v1/user/profile", 
		server.AuthMiddleware(srv.JWT())(handleProfile()))

	// Admin endpoints (require admin role)
	srv.Handle("/api/v1/admin/users",
		server.AuthMiddleware(srv.JWT())(
			server.RoleMiddleware("admin")(handleAdminUsers())))

	// Example with validation
	srv.Handle("/api/v1/orders",
		server.AuthMiddleware(srv.JWT())(
			server.ValidateMiddleware(validateOrder)(handleCreateOrder())))

	// Start server (includes /health, /ready, /metrics endpoints)
	log.Println("Starting example server...")
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Request/Response types
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type OrderRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// Handlers

func handleLogin(srv *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Validate input
		v := security.NewValidator()
		v.Required("email", req.Email).Email("email", req.Email)
		v.Required("password", req.Password).MinLength("password", req.Password, 8)
		if !v.Result().Valid {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(v.Result().Errors)
			return
		}

		// TODO: Verify password against database
		// In real implementation, you would:
		// 1. Fetch user from database
		// 2. Verify password with security.VerifyPassword(password, hash)
		
		// Generate tokens
		tokens, err := srv.JWT().GenerateTokenPair(
			"user_123",  // user ID from database
			"user",      // role from database
			map[string]string{"email": req.Email},
		)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Increment metrics
		srv.Metrics().Counter("logins_total").Inc()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresIn:    tokens.ExpiresIn,
		})
	})
}

func handleRegister() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Example with password hashing
		password := "user_password"
		hash, err := security.HashPassword(password, nil) // Uses Argon2id
		if err != nil {
			http.Error(w, "Registration failed", http.StatusInternalServerError)
			return
		}

		// TODO: Save user with hash to database
		_ = hash

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"User registered successfully"}`))
	})
}

func handleProducts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This response will be:
		// - Cached (via ResponseCache)
		// - Compressed (via gzip)
		// - Rate limited
		// - Traced
		// All automatically!

		products := []map[string]interface{}{
			{"id": "1", "name": "Product 1", "price": 100000},
			{"id": "2", "name": "Product 2", "price": 200000},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})
}

func handleProfile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user from context
		userID := server.GetUserID(r.Context())
		role := server.GetUserRole(r.Context())

		profile := map[string]string{
			"user_id": userID,
			"role":    role,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	})
}

func handleAdminUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only accessible by admins (checked by RoleMiddleware)
		
		users := []map[string]string{
			{"id": "1", "email": "admin@example.com"},
			{"id": "2", "email": "user@example.com"},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})
}

func validateOrder(order *OrderRequest) *security.ValidationResult {
	v := security.NewValidator()
	v.Required("product_id", order.ProductID)
	if order.Quantity <= 0 {
		v.Result().AddError("quantity", "must be positive")
	}
	return v.Result()
}

func handleCreateOrder() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get validated body from context
		order := server.GetRequestBody[OrderRequest](r.Context())
		userID := server.GetUserID(r.Context())

		// Create order...
		response := map[string]interface{}{
			"order_id":   "ORD-12345",
			"user_id":    userID,
			"product_id": order.ProductID,
			"quantity":   order.Quantity,
			"status":     "pending",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	})
}
