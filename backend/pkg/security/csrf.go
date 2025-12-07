package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"sync"
	"time"
)

var (
	ErrCSRFTokenMissing  = errors.New("CSRF token missing")
	ErrCSRFTokenInvalid  = errors.New("CSRF token invalid")
	ErrCSRFTokenExpired  = errors.New("CSRF token expired")
)

// CSRFConfig configures CSRF protection
type CSRFConfig struct {
	TokenLength   int
	CookieName    string
	HeaderName    string
	FormField     string
	Secure        bool
	SameSite      http.SameSite
	MaxAge        int
	ExcludePaths  []string
}

// DefaultCSRFConfig returns secure CSRF defaults
func DefaultCSRFConfig() *CSRFConfig {
	return &CSRFConfig{
		TokenLength:  32,
		CookieName:   "csrf_token",
		HeaderName:   "X-CSRF-Token",
		FormField:    "csrf_token",
		Secure:       true,
		SameSite:     http.SameSiteStrictMode,
		MaxAge:       3600, // 1 hour
		ExcludePaths: []string{"/api/health", "/api/metrics"},
	}
}

// CSRFProtection provides CSRF protection
type CSRFProtection struct {
	config *CSRFConfig
	tokens sync.Map // In-memory store for tokens
}

// NewCSRFProtection creates CSRF protection
func NewCSRFProtection(config *CSRFConfig) *CSRFProtection {
	if config == nil {
		config = DefaultCSRFConfig()
	}
	return &CSRFProtection{config: config}
}

// GenerateToken creates a new CSRF token
func (c *CSRFProtection) GenerateToken() (string, error) {
	bytes := make([]byte, c.config.TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(bytes)
	
	// Store token with expiry
	c.tokens.Store(token, time.Now().Add(time.Duration(c.config.MaxAge)*time.Second))
	
	return token, nil
}

// ValidateToken validates a CSRF token
func (c *CSRFProtection) ValidateToken(token string) error {
	if token == "" {
		return ErrCSRFTokenMissing
	}

	expiry, ok := c.tokens.Load(token)
	if !ok {
		return ErrCSRFTokenInvalid
	}

	if time.Now().After(expiry.(time.Time)) {
		c.tokens.Delete(token)
		return ErrCSRFTokenExpired
	}

	return nil
}

// Middleware creates CSRF protection middleware
func (c *CSRFProtection) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip for excluded paths
		for _, path := range c.config.ExcludePaths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Skip for safe methods
		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
			// Set CSRF token cookie for safe requests
			token, err := c.GenerateToken()
			if err == nil {
				http.SetCookie(w, &http.Cookie{
					Name:     c.config.CookieName,
					Value:    token,
					Path:     "/",
					Secure:   c.config.Secure,
					HttpOnly: false, // JS needs to read it
					SameSite: c.config.SameSite,
					MaxAge:   c.config.MaxAge,
				})
			}
			next.ServeHTTP(w, r)
			return
		}

		// Validate CSRF for unsafe methods
		token := r.Header.Get(c.config.HeaderName)
		if token == "" {
			token = r.FormValue(c.config.FormField)
		}

		// Also check cookie
		cookie, err := r.Cookie(c.config.CookieName)
		if err != nil {
			http.Error(w, "CSRF token missing", http.StatusForbidden)
			return
		}

		// Double submit cookie pattern - compare header token with cookie
		if subtle.ConstantTimeCompare([]byte(token), []byte(cookie.Value)) != 1 {
			http.Error(w, "CSRF token mismatch", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SessionCSRF ties CSRF tokens to sessions
type SessionCSRF struct {
	sessions sync.Map
}

// NewSessionCSRF creates session-based CSRF
func NewSessionCSRF() *SessionCSRF {
	return &SessionCSRF{}
}

// GenerateForSession creates a CSRF token for a session
func (s *SessionCSRF) GenerateForSession(sessionID string) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(bytes)
	
	s.sessions.Store(sessionID, token)
	return token, nil
}

// ValidateForSession validates a CSRF token against a session
func (s *SessionCSRF) ValidateForSession(sessionID, token string) bool {
	expected, ok := s.sessions.Load(sessionID)
	if !ok {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(expected.(string)), []byte(token)) == 1
}
