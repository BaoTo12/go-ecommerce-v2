package security

import (
	"net/http"
	"strconv"
	"strings"
)

// SecurityHeaders contains all security header configurations
type SecurityHeaders struct {
	ContentSecurityPolicy   string
	XContentTypeOptions     string
	XFrameOptions           string
	XXSSProtection          string
	StrictTransportSecurity string
	ReferrerPolicy          string
	PermissionsPolicy       string
	CacheControl            string
}

// DefaultSecurityHeaders returns production-ready security headers
func DefaultSecurityHeaders() *SecurityHeaders {
	return &SecurityHeaders{
		ContentSecurityPolicy: strings.Join([]string{
			"default-src 'self'",
			"script-src 'self' 'unsafe-inline'",
			"style-src 'self' 'unsafe-inline'",
			"img-src 'self' data: https:",
			"font-src 'self'",
			"connect-src 'self'",
			"frame-ancestors 'none'",
			"base-uri 'self'",
			"form-action 'self'",
		}, "; "),
		XContentTypeOptions:     "nosniff",
		XFrameOptions:           "DENY",
		XXSSProtection:          "1; mode=block",
		StrictTransportSecurity: "max-age=31536000; includeSubDomains; preload",
		ReferrerPolicy:          "strict-origin-when-cross-origin",
		PermissionsPolicy: strings.Join([]string{
			"accelerometer=()",
			"camera=()",
			"geolocation=()",
			"gyroscope=()",
			"magnetometer=()",
			"microphone=()",
			"payment=()",
			"usb=()",
		}, ", "),
		CacheControl: "no-store, no-cache, must-revalidate, proxy-revalidate",
	}
}

// StrictSecurityHeaders returns stricter headers for sensitive endpoints
func StrictSecurityHeaders() *SecurityHeaders {
	headers := DefaultSecurityHeaders()
	headers.ContentSecurityPolicy = "default-src 'none'; frame-ancestors 'none'"
	headers.CacheControl = "no-store, max-age=0"
	return headers
}

// Middleware creates security headers middleware
func (h *SecurityHeaders) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set security headers
		if h.ContentSecurityPolicy != "" {
			w.Header().Set("Content-Security-Policy", h.ContentSecurityPolicy)
		}
		if h.XContentTypeOptions != "" {
			w.Header().Set("X-Content-Type-Options", h.XContentTypeOptions)
		}
		if h.XFrameOptions != "" {
			w.Header().Set("X-Frame-Options", h.XFrameOptions)
		}
		if h.XXSSProtection != "" {
			w.Header().Set("X-XSS-Protection", h.XXSSProtection)
		}
		if h.StrictTransportSecurity != "" && r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", h.StrictTransportSecurity)
		}
		if h.ReferrerPolicy != "" {
			w.Header().Set("Referrer-Policy", h.ReferrerPolicy)
		}
		if h.PermissionsPolicy != "" {
			w.Header().Set("Permissions-Policy", h.PermissionsPolicy)
		}
		if h.CacheControl != "" {
			w.Header().Set("Cache-Control", h.CacheControl)
		}

		// Additional security
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Del("X-Powered-By")
		w.Header().Del("Server")

		next.ServeHTTP(w, r)
	})
}

// CORSConfig configures CORS
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns secure CORS defaults
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-CSRF-Token"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// CORSMiddleware creates CORS middleware
func CORSMiddleware(config *CORSConfig) func(http.Handler) http.Handler {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, o := range config.AllowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}

			if allowed && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				
				if config.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				
				if len(config.ExposedHeaders) > 0 {
					w.Header().Set("Access-Control-Expose-Headers", 
						strings.Join(config.ExposedHeaders, ", "))
				}
			}

			// Handle preflight
			if r.Method == "OPTIONS" {
				if len(config.AllowedMethods) > 0 {
					w.Header().Set("Access-Control-Allow-Methods", 
						strings.Join(config.AllowedMethods, ", "))
				}
				if len(config.AllowedHeaders) > 0 {
					w.Header().Set("Access-Control-Allow-Headers", 
						strings.Join(config.AllowedHeaders, ", "))
				}
				if config.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", 
						strconv.Itoa(config.MaxAge))
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
