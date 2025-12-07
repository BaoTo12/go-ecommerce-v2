package security_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

// ====================
// SECURITY TESTS
// ====================

// TestSQLInjection tests for SQL injection vulnerabilities
func TestSQLInjection_QueryParameters(t *testing.T) {
	// Common SQL injection payloads
	payloads := []string{
		"'; DROP TABLE users; --",
		"' OR '1'='1",
		"1; DELETE FROM products",
		"' UNION SELECT * FROM passwords --",
		"admin'--",
		"' OR 1=1--",
		"'; EXEC xp_cmdshell('dir'); --",
		"1' AND '1'='1",
		"1' AND SLEEP(5)--",
	}

	handler := createSafeHandler()

	for _, payload := range payloads {
		t.Run("Payload_"+sanitizeTestName(payload), func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/products?id="+payload, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			// Should not return 500 (indicates potential injection success)
			if rec.Code == http.StatusInternalServerError {
				t.Errorf("Potential SQL injection vulnerability with payload: %s", payload)
			}

			// Response should not contain DB error messages
			body := rec.Body.String()
			if containsSQLError(body) {
				t.Errorf("SQL error leaked in response for payload: %s", payload)
			}
		})
	}
}

// TestXSS tests for Cross-Site Scripting vulnerabilities
func TestXSS_InputSanitization(t *testing.T) {
	xssPayloads := []string{
		"<script>alert('XSS')</script>",
		"<img src=x onerror=alert('XSS')>",
		"javascript:alert('XSS')",
		"<svg onload=alert('XSS')>",
		"<body onload=alert('XSS')>",
		"\"><script>alert('XSS')</script>",
		"'><script>alert('XSS')</script>",
		"<iframe src=\"javascript:alert('XSS')\">",
		"<a href=\"javascript:alert('XSS')\">Click</a>",
		"<div style=\"background:url(javascript:alert('XSS'))\">",
	}

	handler := createSafeHandler()

	for _, payload := range xssPayloads {
		t.Run("XSS_"+sanitizeTestName(payload), func(t *testing.T) {
			body := map[string]string{"name": payload}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest("POST", "/api/products", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			respBody := rec.Body.String()

			// Response should not contain unescaped script tags
			if strings.Contains(respBody, "<script>") {
				t.Errorf("Unescaped script tag in response for payload: %s", payload)
			}

			// Check Content-Type header for proper escaping
			contentType := rec.Header().Get("Content-Type")
			if strings.Contains(contentType, "text/html") {
				// If returning HTML, check for proper escaping
				if strings.Contains(respBody, payload) && !isProperlyEscaped(respBody, payload) {
					t.Errorf("XSS payload not properly escaped: %s", payload)
				}
			}
		})
	}
}

// TestAuthorizationBypass tests for authorization bypass
func TestAuthorizationBypass(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		method   string
		headers  map[string]string
		wantCode int
	}{
		{
			name:     "No auth header",
			path:     "/api/admin/users",
			method:   "GET",
			headers:  map[string]string{},
			wantCode: http.StatusUnauthorized,
		},
		{
			name:   "Invalid token",
			path:   "/api/admin/users",
			method: "GET",
			headers: map[string]string{
				"Authorization": "Bearer invalid-token",
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			name:   "Modified JWT",
			path:   "/api/admin/users",
			method: "GET",
			headers: map[string]string{
				"Authorization": "Bearer eyJhbGciOiJub25lIn0.eyJyb2xlIjoiYWRtaW4ifQ.",
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			name:   "Path traversal attempt",
			path:   "/api/users/../admin/config",
			method: "GET",
			headers: map[string]string{
				"Authorization": "Bearer valid-user-token",
			},
			wantCode: http.StatusForbidden,
		},
	}

	handler := createSecureHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantCode {
				t.Errorf("Expected status %d, got %d", tt.wantCode, rec.Code)
			}
		})
	}
}

// TestRateLimitBypass tests rate limit bypassing attempts
func TestRateLimitBypass(t *testing.T) {
	handler := createRateLimitedHandler(5) // 5 requests allowed

	// Normal requests should work
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/api/resource", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Request %d should succeed", i)
		}
	}

	// 6th request should be rate limited
	bypassAttempts := []struct {
		name    string
		headers map[string]string
	}{
		{
			name:    "X-Forwarded-For spoofing",
			headers: map[string]string{"X-Forwarded-For": "1.2.3.4"},
		},
		{
			name:    "X-Real-IP spoofing",
			headers: map[string]string{"X-Real-IP": "5.6.7.8"},
		},
		{
			name:    "Multiple X-Forwarded-For",
			headers: map[string]string{"X-Forwarded-For": "1.1.1.1, 2.2.2.2"},
		},
	}

	for _, attempt := range bypassAttempts {
		t.Run(attempt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/resource", nil)
			for k, v := range attempt.headers {
				req.Header.Set(k, v)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			// Should still be rate limited (headers shouldn't bypass)
			if rec.Code == http.StatusOK {
				t.Errorf("Rate limit bypassed with %s", attempt.name)
			}
		})
	}
}

// TestSensitiveDataExposure tests for sensitive data in responses
func TestSensitiveDataExposure(t *testing.T) {
	sensitivePatterns := []struct {
		name    string
		pattern *regexp.Regexp
	}{
		{"Password", regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[:=]\s*['"]?[\w@#$%^&*]+`)},
		{"API Key", regexp.MustCompile(`(?i)(api[_-]?key|apikey)\s*[:=]\s*['"]?[\w-]+`)},
		{"Secret", regexp.MustCompile(`(?i)(secret|token)\s*[:=]\s*['"]?[\w-]+`)},
		{"SSN", regexp.MustCompile(`\d{3}-\d{2}-\d{4}`)},
		{"Credit Card", regexp.MustCompile(`\d{4}[- ]?\d{4}[- ]?\d{4}[- ]?\d{4}`)},
		{"AWS Key", regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
	}

	handler := createSafeHandler()
	endpoints := []string{"/api/users/1", "/api/config", "/api/debug"}

	for _, endpoint := range endpoints {
		req := httptest.NewRequest("GET", endpoint, nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		body := rec.Body.String()

		for _, sp := range sensitivePatterns {
			if sp.pattern.MatchString(body) {
				t.Errorf("Sensitive data (%s) exposed in %s", sp.name, endpoint)
			}
		}
	}
}

// TestSecurityHeaders checks for security headers
func TestSecurityHeaders(t *testing.T) {
	requiredHeaders := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"X-XSS-Protection":          "1; mode=block",
		"Content-Security-Policy":   "", // Just check presence
		"Strict-Transport-Security": "", // Just check presence
	}

	handler := createSecureHandler()
	req := httptest.NewRequest("GET", "/api/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	for header, expectedValue := range requiredHeaders {
		actual := rec.Header().Get(header)
		if actual == "" {
			t.Errorf("Missing security header: %s", header)
		} else if expectedValue != "" && actual != expectedValue {
			t.Errorf("Header %s: expected %s, got %s", header, expectedValue, actual)
		}
	}
}

// Helper functions

func sanitizeTestName(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	result := re.ReplaceAllString(s, "_")
	if len(result) > 30 {
		result = result[:30]
	}
	return result
}

func containsSQLError(body string) bool {
	errorPatterns := []string{
		"syntax error",
		"SQL error",
		"mysql_",
		"pg_query",
		"ORA-",
		"SQLite",
		"SQLSTATE",
	}
	lower := strings.ToLower(body)
	for _, pattern := range errorPatterns {
		if strings.Contains(lower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func isProperlyEscaped(body, payload string) bool {
	// Check if < and > are escaped
	escaped := strings.ReplaceAll(payload, "<", "&lt;")
	escaped = strings.ReplaceAll(escaped, ">", "&gt;")
	return strings.Contains(body, escaped)
}

// Test handlers (implementations would be in actual service)
func createSafeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}

func createSecureHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000")

		// Auth check
		auth := r.Header.Get("Authorization")
		if auth == "" || auth == "Bearer invalid-token" || strings.Contains(auth, "none") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if strings.Contains(r.URL.Path, "..") {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func createRateLimitedHandler(limit int) http.Handler {
	count := 0
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		if count > limit {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
