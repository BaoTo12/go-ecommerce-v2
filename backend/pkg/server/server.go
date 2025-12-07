package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/titan-commerce/backend/pkg/metrics"
	"github.com/titan-commerce/backend/pkg/perf"
	"github.com/titan-commerce/backend/pkg/ratelimit"
	"github.com/titan-commerce/backend/pkg/resilience"
	"github.com/titan-commerce/backend/pkg/security"
	"github.com/titan-commerce/backend/pkg/tracing"
)

// Config holds all server configuration
type Config struct {
	// Server settings
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration

	// Security
	JWTSecret        []byte
	CSRFEnabled      bool
	CORSOrigins      []string
	SecurityHeaders  bool

	// Performance
	CompressionEnabled bool
	CacheEnabled       bool
	CacheTTL           time.Duration

	// Rate limiting
	RateLimitEnabled bool
	RateLimitRate    float64
	RateLimitBurst   int64

	// Tracing
	TracingEnabled bool
	ServiceName    string

	// Metrics
	MetricsEnabled bool
}

// DefaultConfig returns production-ready defaults
func DefaultConfig() *Config {
	return &Config{
		Host:            "0.0.0.0",
		Port:            "8080",
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 30 * time.Second,

		CSRFEnabled:     true,
		CORSOrigins:     []string{},
		SecurityHeaders: true,

		CompressionEnabled: true,
		CacheEnabled:       true,
		CacheTTL:           5 * time.Minute,

		RateLimitEnabled: true,
		RateLimitRate:    100,
		RateLimitBurst:   200,

		TracingEnabled: true,
		ServiceName:    "titan-commerce",

		MetricsEnabled: true,
	}
}

// Server is the unified HTTP server
type Server struct {
	config   *Config
	router   *http.ServeMux
	server   *http.Server
	
	// Integrated components
	jwt         *security.JWTService
	csrf        *security.CSRFProtection
	rateLimiter *ratelimit.TokenBucket
	tracer      *tracing.Tracer
	metrics     *metrics.Registry
	cache       *perf.ResponseCache
	circuitBreaker *resilience.CircuitBreaker

	// Lifecycle
	wg       sync.WaitGroup
	shutdown chan struct{}
}

// New creates a new integrated server
func New(config *Config) *Server {
	if config == nil {
		config = DefaultConfig()
	}

	s := &Server{
		config:   config,
		router:   http.NewServeMux(),
		shutdown: make(chan struct{}),
	}

	// Initialize integrated components
	s.initComponents()

	return s
}

func (s *Server) initComponents() {
	// JWT Service
	if len(s.config.JWTSecret) > 0 {
		s.jwt = security.NewJWTService(security.DefaultJWTConfig(s.config.JWTSecret))
	}

	// CSRF Protection
	if s.config.CSRFEnabled {
		s.csrf = security.NewCSRFProtection(nil)
	}

	// Rate Limiter
	if s.config.RateLimitEnabled {
		s.rateLimiter = ratelimit.NewTokenBucket(s.config.RateLimitBurst, s.config.RateLimitRate)
	}

	// Tracing
	if s.config.TracingEnabled {
		s.tracer = tracing.NewTracer(s.config.ServiceName, &tracing.ConsoleExporter{}, nil)
	}

	// Metrics
	if s.config.MetricsEnabled {
		s.metrics = metrics.NewRegistry()
		metrics.NewRuntimeMetrics(s.metrics)
	}

	// Response Cache
	if s.config.CacheEnabled {
		s.cache = perf.NewResponseCache(s.config.CacheTTL)
	}

	// Circuit Breaker (for external calls)
	s.circuitBreaker = resilience.NewCircuitBreaker(resilience.DefaultCircuitConfig("external-api"))
}

// Handle registers a handler with all middleware applied
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.router.Handle(pattern, s.applyMiddleware(handler))
}

// HandleFunc registers a handler function with all middleware applied
func (s *Server) HandleFunc(pattern string, handler http.HandlerFunc) {
	s.Handle(pattern, handler)
}

// applyMiddleware chains all middleware together
func (s *Server) applyMiddleware(handler http.Handler) http.Handler {
	// Apply in reverse order (first applied = outermost)
	h := handler

	// 1. Recovery (innermost - catches panics)
	h = s.recoveryMiddleware(h)

	// 2. Request ID
	h = s.requestIDMiddleware(h)

	// 3. Tracing
	if s.config.TracingEnabled && s.tracer != nil {
		h = tracing.HTTPMiddleware(s.tracer)(h)
	}

	// 4. Metrics
	if s.config.MetricsEnabled && s.metrics != nil {
		h = metrics.HTTPMiddleware(s.metrics)(h)
	}

	// 5. Rate Limiting
	if s.config.RateLimitEnabled && s.rateLimiter != nil {
		h = ratelimit.HTTPMiddleware(s.rateLimiter)(h)
	}

	// 6. Compression
	if s.config.CompressionEnabled {
		h = perf.CompressionMiddleware(1024)(h)
	}

	// 7. Security Headers
	if s.config.SecurityHeaders {
		headers := security.DefaultSecurityHeaders()
		h = headers.Middleware(h)
	}

	// 8. CORS
	if len(s.config.CORSOrigins) > 0 {
		corsConfig := security.DefaultCORSConfig()
		corsConfig.AllowedOrigins = s.config.CORSOrigins
		h = security.CORSMiddleware(corsConfig)(h)
	}

	// 9. CSRF Protection
	if s.config.CSRFEnabled && s.csrf != nil {
		h = s.csrf.Middleware(h)
	}

	// 10. Logging (outermost)
	h = s.loggingMiddleware(h)

	return h
}

// Middleware implementations

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type requestIDKey struct{}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		next.ServeHTTP(w, r)
		
		log.Printf("[%s] %s %s %v",
			r.Context().Value(requestIDKey{}),
			r.Method,
			r.URL.Path,
			time.Since(start))
	})
}

func generateRequestID() string {
	token, _ := security.GenerateSecureToken(8)
	return token
}

// Built-in endpoints

func (s *Server) registerBuiltinEndpoints() {
	// Health check
	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Readiness check
	s.router.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ready":true}`))
	})

	// Metrics endpoint
	if s.config.MetricsEnabled {
		s.router.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			snapshot := s.metrics.Snapshot()
			// Simple JSON output
			w.Write([]byte("{"))
			first := true
			for k, v := range snapshot {
				if !first {
					w.Write([]byte(","))
				}
				first = false
				w.Write([]byte(`"` + k + `":`))
				switch val := v.(type) {
				case uint64:
					w.Write([]byte(formatUint64(val)))
				case int64:
					w.Write([]byte(formatInt64(val)))
				default:
					w.Write([]byte(`"` + formatAny(v) + `"`))
				}
			}
			w.Write([]byte("}"))
		})
	}

	// Debug/profiling (only in development)
	s.router.Handle("/debug/profile", perf.ProfileHandler())
}

// Start starts the server and blocks until shutdown
func (s *Server) Start() error {
	s.registerBuiltinEndpoints()

	addr := s.config.Host + ":" + s.config.Port
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	// Graceful shutdown handler
	go s.handleShutdown()

	log.Printf("Server starting on %s", addr)
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) handleShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	close(s.shutdown)
	s.wg.Wait()
	log.Println("Server stopped")
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// JWT returns the JWT service
func (s *Server) JWT() *security.JWTService {
	return s.jwt
}

// Metrics returns the metrics registry
func (s *Server) Metrics() *metrics.Registry {
	return s.metrics
}

// CircuitBreaker returns the circuit breaker for external calls
func (s *Server) CircuitBreaker() *resilience.CircuitBreaker {
	return s.circuitBreaker
}

// Helper formatters
func formatUint64(v uint64) string {
	if v == 0 {
		return "0"
	}
	buf := make([]byte, 20)
	pos := len(buf)
	for v > 0 {
		pos--
		buf[pos] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[pos:])
}

func formatInt64(v int64) string {
	if v < 0 {
		return "-" + formatUint64(uint64(-v))
	}
	return formatUint64(uint64(v))
}

func formatAny(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		return "unknown"
	}
}
