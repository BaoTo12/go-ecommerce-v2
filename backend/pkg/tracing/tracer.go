package tracing

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Span represents a single operation in a trace
type Span struct {
	TraceID      string            `json:"trace_id"`
	SpanID       string            `json:"span_id"`
	ParentSpanID string            `json:"parent_span_id,omitempty"`
	OperationName string           `json:"operation_name"`
	ServiceName  string            `json:"service_name"`
	StartTime    time.Time         `json:"start_time"`
	Duration     time.Duration     `json:"duration"`
	Status       SpanStatus        `json:"status"`
	Tags         map[string]string `json:"tags"`
	Logs         []SpanLog         `json:"logs"`
	mu           sync.Mutex
}

// SpanStatus represents span completion status
type SpanStatus string

const (
	StatusOK    SpanStatus = "OK"
	StatusError SpanStatus = "ERROR"
)

// SpanLog represents a log entry in a span
type SpanLog struct {
	Timestamp time.Time         `json:"timestamp"`
	Fields    map[string]string `json:"fields"`
}

// SpanContext carries trace context across boundaries
type SpanContext struct {
	TraceID      string
	SpanID       string
	ParentSpanID string
	Sampled      bool
}

// Tracer creates and manages spans
type Tracer struct {
	serviceName string
	exporter    SpanExporter
	sampler     Sampler
}

// SpanExporter exports completed spans
type SpanExporter interface {
	Export(span *Span) error
}

// Sampler decides if a trace should be sampled
type Sampler interface {
	ShouldSample(traceID string) bool
}

// AlwaysSampler samples all traces
type AlwaysSampler struct{}

func (s *AlwaysSampler) ShouldSample(traceID string) bool { return true }

// ProbabilitySampler samples based on probability
type ProbabilitySampler struct {
	rate float64
}

func NewProbabilitySampler(rate float64) *ProbabilitySampler {
	return &ProbabilitySampler{rate: rate}
}

func (s *ProbabilitySampler) ShouldSample(traceID string) bool {
	return rand.Float64() < s.rate
}

// NewTracer creates a new tracer
func NewTracer(serviceName string, exporter SpanExporter, sampler Sampler) *Tracer {
	if sampler == nil {
		sampler = &AlwaysSampler{}
	}
	return &Tracer{
		serviceName: serviceName,
		exporter:    exporter,
		sampler:     sampler,
	}
}

// StartSpan creates a new span
func (t *Tracer) StartSpan(ctx context.Context, operationName string) (context.Context, *Span) {
	var parentSpanID string
	var traceID string

	// Check for parent span in context
	if parentSpan := SpanFromContext(ctx); parentSpan != nil {
		traceID = parentSpan.TraceID
		parentSpanID = parentSpan.SpanID
	} else {
		traceID = generateID()
	}

	span := &Span{
		TraceID:       traceID,
		SpanID:        generateID(),
		ParentSpanID:  parentSpanID,
		OperationName: operationName,
		ServiceName:   t.serviceName,
		StartTime:     time.Now(),
		Status:        StatusOK,
		Tags:          make(map[string]string),
		Logs:          make([]SpanLog, 0),
	}

	return context.WithValue(ctx, spanContextKey{}, span), span
}

// FinishSpan completes and exports a span
func (t *Tracer) FinishSpan(span *Span) {
	span.Duration = time.Since(span.StartTime)
	if t.exporter != nil && t.sampler.ShouldSample(span.TraceID) {
		t.exporter.Export(span)
	}
}

// SetTag adds a tag to the span
func (s *Span) SetTag(key, value string) *Span {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Tags[key] = value
	return s
}

// LogFields adds a log entry to the span
func (s *Span) LogFields(fields map[string]string) *Span {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Logs = append(s.Logs, SpanLog{
		Timestamp: time.Now(),
		Fields:    fields,
	})
	return s
}

// SetError marks the span as errored
func (s *Span) SetError(err error) *Span {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Status = StatusError
	s.Tags["error"] = "true"
	s.Tags["error.message"] = err.Error()
	return s
}

// Context key for span
type spanContextKey struct{}

// SpanFromContext retrieves span from context
func SpanFromContext(ctx context.Context) *Span {
	if span, ok := ctx.Value(spanContextKey{}).(*Span); ok {
		return span
	}
	return nil
}

// HTTPMiddleware creates tracing middleware
func HTTPMiddleware(tracer *Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracer.StartSpan(r.Context(), r.Method+" "+r.URL.Path)
			defer tracer.FinishSpan(span)

			span.SetTag("http.method", r.Method)
			span.SetTag("http.url", r.URL.String())
			span.SetTag("http.user_agent", r.UserAgent())

			// Wrap response writer to capture status
			sw := &statusWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(sw, r.WithContext(ctx))

			span.SetTag("http.status_code", fmt.Sprintf("%d", sw.status))
			if sw.status >= 400 {
				span.Status = StatusError
			}
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// InjectHTTP injects trace context into HTTP headers
func InjectHTTP(ctx context.Context, req *http.Request) {
	if span := SpanFromContext(ctx); span != nil {
		req.Header.Set("X-Trace-ID", span.TraceID)
		req.Header.Set("X-Span-ID", span.SpanID)
	}
}

// ExtractHTTP extracts trace context from HTTP headers
func ExtractHTTP(r *http.Request) *SpanContext {
	traceID := r.Header.Get("X-Trace-ID")
	spanID := r.Header.Get("X-Span-ID")
	if traceID == "" {
		return nil
	}
	return &SpanContext{
		TraceID:      traceID,
		ParentSpanID: spanID,
	}
}

// ConsoleExporter exports spans to console (for dev)
type ConsoleExporter struct{}

func (e *ConsoleExporter) Export(span *Span) error {
	fmt.Printf("[TRACE] %s | %s | %s | %v | %s\n",
		span.TraceID[:8], span.ServiceName, span.OperationName,
		span.Duration, span.Status)
	return nil
}

func generateID() string {
	return fmt.Sprintf("%x", rand.Int63())
}
