package grpcx

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

/*
gRPC INTERCEPTORS - Enterprise gRPC Middleware

Provides common interceptors for gRPC:
- Logging
- Metrics
- Tracing
- Rate limiting
- Authentication
- Recovery
- Timeout
- Retry
*/

// UnaryServerInterceptor is the gRPC server interceptor signature
type UnaryServerInterceptor func(
	ctx context.Context,
	req interface{},
	info *UnaryServerInfo,
	handler UnaryHandler,
) (interface{}, error)

// UnaryHandler is the handler signature
type UnaryHandler func(ctx context.Context, req interface{}) (interface{}, error)

// UnaryServerInfo contains info about the method
type UnaryServerInfo struct {
	FullMethod string
	Server     interface{}
}

// ChainUnaryServer chains multiple interceptors
func ChainUnaryServer(interceptors ...UnaryServerInterceptor) UnaryServerInterceptor {
	n := len(interceptors)
	return func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (interface{}, error) {
		chainer := func(currentInter UnaryServerInterceptor, currentHandler UnaryHandler) UnaryHandler {
			return func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return currentInter(currentCtx, currentReq, info, currentHandler)
			}
		}

		chainedHandler := handler
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(interceptors[i], chainedHandler)
		}

		return chainedHandler(ctx, req)
	}
}

// LoggingInterceptor logs requests
func LoggingInterceptor() UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (interface{}, error) {
		start := time.Now()
		
		resp, err := handler(ctx, req)
		
		duration := time.Since(start)
		status := "OK"
		if err != nil {
			status = "ERROR"
		}
		
		log.Printf("[gRPC] %s | %s | %v", info.FullMethod, status, duration)
		
		return resp, err
	}
}

// MetricsInterceptor collects metrics
type GRPCMetrics struct {
	TotalRequests    *atomic.Int64
	SuccessRequests  *atomic.Int64
	FailedRequests   *atomic.Int64
	TotalLatencyMs   *atomic.Int64
	MethodCounts     sync.Map // map[string]*atomic.Int64
}

func NewGRPCMetrics() *GRPCMetrics {
	return &GRPCMetrics{
		TotalRequests:   new(atomic.Int64),
		SuccessRequests: new(atomic.Int64),
		FailedRequests:  new(atomic.Int64),
		TotalLatencyMs:  new(atomic.Int64),
	}
}

func MetricsInterceptor(metrics *GRPCMetrics) UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (interface{}, error) {
		start := time.Now()
		
		metrics.TotalRequests.Add(1)
		
		// Per-method count
		countI, _ := metrics.MethodCounts.LoadOrStore(info.FullMethod, new(atomic.Int64))
		count := countI.(*atomic.Int64)
		count.Add(1)
		
		resp, err := handler(ctx, req)
		
		duration := time.Since(start)
		metrics.TotalLatencyMs.Add(duration.Milliseconds())
		
		if err != nil {
			metrics.FailedRequests.Add(1)
		} else {
			metrics.SuccessRequests.Add(1)
		}
		
		return resp, err
	}
}

// TracingInterceptor adds tracing
func TracingInterceptor() UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (interface{}, error) {
		// Extract trace ID from context or generate new
		traceID := extractTraceID(ctx)
		if traceID == "" {
			traceID = generateTraceID()
		}
		
		// Add to context
		ctx = context.WithValue(ctx, traceIDKey{}, traceID)
		
		// Create span
		span := &Span{
			TraceID:   traceID,
			SpanID:    generateSpanID(),
			Name:      info.FullMethod,
			StartTime: time.Now(),
		}
		
		resp, err := handler(ctx, req)
		
		span.EndTime = time.Now()
		span.Duration = span.EndTime.Sub(span.StartTime)
		if err != nil {
			span.Error = err.Error()
		}
		
		// Log span (in real implementation, send to tracing backend)
		log.Printf("[TRACE] %s | %s | %v", span.TraceID, span.Name, span.Duration)
		
		return resp, err
	}
}

type traceIDKey struct{}

type Span struct {
	TraceID   string
	SpanID    string
	ParentID  string
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Error     string
}

func extractTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDKey{}).(string); ok {
		return id
	}
	return ""
}

func generateTraceID() string {
	return "trace_" + time.Now().Format("20060102150405.000000")
}

func generateSpanID() string {
	return "span_" + time.Now().Format("150405.000")
}

// RateLimitInterceptor limits requests
type RateLimiter struct {
	rate     float64
	capacity float64
	tokens   float64
	lastTime time.Time
	mu       sync.Mutex
}

func NewRateLimiter(rate float64, capacity float64) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		capacity: capacity,
		tokens:   capacity,
		lastTime: time.Now(),
	}
}

func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(r.lastTime).Seconds()
	r.tokens += elapsed * r.rate
	if r.tokens > r.capacity {
		r.tokens = r.capacity
	}
	r.lastTime = now
	
	if r.tokens >= 1 {
		r.tokens--
		return true
	}
	return false
}

func RateLimitInterceptor(limiter *RateLimiter) UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (interface{}, error) {
		if !limiter.Allow() {
			return nil, ErrRateLimited
		}
		return handler(ctx, req)
	}
}

var ErrRateLimited = &GRPCError{Code: 8, Message: "rate limited"}

// AuthInterceptor validates authentication
type AuthValidator func(ctx context.Context, token string) (userID string, err error)

func AuthInterceptor(validator AuthValidator) UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (interface{}, error) {
		// Extract token from metadata
		token := extractToken(ctx)
		if token == "" {
			return nil, ErrUnauthenticated
		}
		
		userID, err := validator(ctx, token)
		if err != nil {
			return nil, ErrUnauthenticated
		}
		
		// Add user to context
		ctx = context.WithValue(ctx, userIDKey{}, userID)
		
		return handler(ctx, req)
	}
}

type userIDKey struct{}

func extractToken(ctx context.Context) string {
	// In real implementation, extract from gRPC metadata
	return ""
}

var ErrUnauthenticated = &GRPCError{Code: 16, Message: "unauthenticated"}

// RecoveryInterceptor recovers from panics
func RecoveryInterceptor() UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] %s | %v", info.FullMethod, r)
				err = ErrInternal
			}
		}()
		
		return handler(ctx, req)
	}
}

var ErrInternal = &GRPCError{Code: 13, Message: "internal error"}

// TimeoutInterceptor adds timeout
func TimeoutInterceptor(timeout time.Duration) UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		
		return handler(ctx, req)
	}
}

// RetryInterceptor (client-side) - retries failed requests
type RetryConfig struct {
	MaxRetries  int
	InitBackoff time.Duration
	MaxBackoff  time.Duration
	Multiplier  float64
	RetryOn     []int // gRPC status codes
}

func RetryInterceptor(config RetryConfig) UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (interface{}, error) {
		var lastErr error
		backoff := config.InitBackoff
		
		for attempt := 0; attempt <= config.MaxRetries; attempt++ {
			resp, err := handler(ctx, req)
			if err == nil {
				return resp, nil
			}
			
			lastErr = err
			
			// Check if should retry
			if !shouldRetry(err, config.RetryOn) {
				return nil, err
			}
			
			// Wait before retry
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * config.Multiplier)
			if backoff > config.MaxBackoff {
				backoff = config.MaxBackoff
			}
		}
		
		return nil, lastErr
	}
}

func shouldRetry(err error, codes []int) bool {
	if grpcErr, ok := err.(*GRPCError); ok {
		for _, code := range codes {
			if grpcErr.Code == code {
				return true
			}
		}
	}
	return false
}

// GRPCError represents a gRPC error
type GRPCError struct {
	Code    int
	Message string
}

func (e *GRPCError) Error() string {
	return e.Message
}

// StandardInterceptorChain returns a standard chain of interceptors
func StandardInterceptorChain(metrics *GRPCMetrics, rateLimiter *RateLimiter) UnaryServerInterceptor {
	return ChainUnaryServer(
		RecoveryInterceptor(),
		LoggingInterceptor(),
		TracingInterceptor(),
		MetricsInterceptor(metrics),
		RateLimitInterceptor(rateLimiter),
		TimeoutInterceptor(30*time.Second),
	)
}
