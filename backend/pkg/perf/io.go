package perf

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

// BufferedWriter wraps a writer with buffering
type BufferedWriter struct {
	w   *bufio.Writer
	buf []byte
}

// NewBufferedWriter creates a buffered writer
func NewBufferedWriter(w io.Writer, size int) *BufferedWriter {
	return &BufferedWriter{
		w: bufio.NewWriterSize(w, size),
	}
}

func (bw *BufferedWriter) Write(p []byte) (int, error) {
	return bw.w.Write(p)
}

func (bw *BufferedWriter) Flush() error {
	return bw.w.Flush()
}

// CompressedResponseWriter wraps http.ResponseWriter with gzip
type CompressedResponseWriter struct {
	http.ResponseWriter
	gw          *gzip.Writer
	wroteHeader bool
}

// NewCompressedResponseWriter creates a gzip response writer
func NewCompressedResponseWriter(w http.ResponseWriter) *CompressedResponseWriter {
	return &CompressedResponseWriter{
		ResponseWriter: w,
		gw:             gzip.NewWriter(w),
	}
}

func (crw *CompressedResponseWriter) Write(p []byte) (int, error) {
	if !crw.wroteHeader {
		crw.ResponseWriter.Header().Set("Content-Encoding", "gzip")
		crw.ResponseWriter.Header().Del("Content-Length")
		crw.wroteHeader = true
	}
	return crw.gw.Write(p)
}

func (crw *CompressedResponseWriter) Close() error {
	return crw.gw.Close()
}

// CompressionMiddleware adds gzip compression
func CompressionMiddleware(minSize int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if client accepts gzip
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			crw := NewCompressedResponseWriter(w)
			defer crw.Close()

			next.ServeHTTP(crw, r)
		})
	}
}

// Profiler provides runtime profiling
type Profiler struct {
	enabled bool
	mu      sync.Mutex
}

// NewProfiler creates a profiler
func NewProfiler(enabled bool) *Profiler {
	return &Profiler{enabled: enabled}
}

// StartCPUProfile starts CPU profiling
func (p *Profiler) StartCPUProfile(w io.Writer) error {
	if !p.enabled {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return pprof.StartCPUProfile(w)
}

// StopCPUProfile stops CPU profiling
func (p *Profiler) StopCPUProfile() {
	p.mu.Lock()
	defer p.mu.Unlock()
	pprof.StopCPUProfile()
}

// WriteHeapProfile writes heap profile
func (p *Profiler) WriteHeapProfile(w io.Writer) error {
	if !p.enabled {
		return nil
	}
	runtime.GC()
	return pprof.WriteHeapProfile(w)
}

// MemStats returns memory statistics
func (p *Profiler) MemStats() *runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return &m
}

// ProfileHandler creates an HTTP handler for profiling
func ProfileHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"alloc_mb":%.2f,"total_alloc_mb":%.2f,"sys_mb":%.2f,"num_gc":%d,"goroutines":%d}`,
			float64(m.Alloc)/1024/1024,
			float64(m.TotalAlloc)/1024/1024,
			float64(m.Sys)/1024/1024,
			m.NumGC,
			runtime.NumGoroutine())
	})
}

// ConnectionStats tracks connection statistics
type ConnectionStats struct {
	ActiveConnections int64
	TotalConnections  int64
	BytesRead         int64
	BytesWritten      int64
	RequestCount      int64
	ErrorCount        int64
	AvgResponseTime   time.Duration
	responseTimeSum   int64
	mu                sync.Mutex
}

// NewConnectionStats creates connection stats
func NewConnectionStats() *ConnectionStats {
	return &ConnectionStats{}
}

// RecordRequest records a request
func (cs *ConnectionStats) RecordRequest(duration time.Duration, err error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.RequestCount++
	cs.responseTimeSum += int64(duration)
	cs.AvgResponseTime = time.Duration(cs.responseTimeSum / cs.RequestCount)

	if err != nil {
		cs.ErrorCount++
	}
}

// Snapshot returns a copy of stats
func (cs *ConnectionStats) Snapshot() ConnectionStats {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return ConnectionStats{
		ActiveConnections: cs.ActiveConnections,
		TotalConnections:  cs.TotalConnections,
		BytesRead:         cs.BytesRead,
		BytesWritten:      cs.BytesWritten,
		RequestCount:      cs.RequestCount,
		ErrorCount:        cs.ErrorCount,
		AvgResponseTime:   cs.AvgResponseTime,
	}
}

// ResponseCache caches HTTP responses
type ResponseCache struct {
	cache sync.Map
	ttl   time.Duration
}

type cachedResponse struct {
	body      []byte
	status    int
	headers   http.Header
	expiresAt time.Time
}

// NewResponseCache creates a response cache
func NewResponseCache(ttl time.Duration) *ResponseCache {
	return &ResponseCache{ttl: ttl}
}

// Get retrieves a cached response
func (rc *ResponseCache) Get(key string) (*cachedResponse, bool) {
	if v, ok := rc.cache.Load(key); ok {
		cr := v.(*cachedResponse)
		if time.Now().Before(cr.expiresAt) {
			return cr, true
		}
		rc.cache.Delete(key)
	}
	return nil, false
}

// Set caches a response
func (rc *ResponseCache) Set(key string, status int, headers http.Header, body []byte) {
	rc.cache.Store(key, &cachedResponse{
		body:      body,
		status:    status,
		headers:   headers.Clone(),
		expiresAt: time.Now().Add(rc.ttl),
	})
}

// CacheMiddleware creates caching middleware
func CacheMiddleware(cache *ResponseCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only cache GET requests
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			key := r.URL.String()

			// Check cache
			if cached, ok := cache.Get(key); ok {
				for k, v := range cached.headers {
					w.Header()[k] = v
				}
				w.Header().Set("X-Cache", "HIT")
				w.WriteHeader(cached.status)
				w.Write(cached.body)
				return
			}

			// Cache miss - capture response
			rec := &responseRecorder{
				ResponseWriter: w,
				body:           &strings.Builder{},
			}

			next.ServeHTTP(rec, r)

			// Cache successful responses
			if rec.status >= 200 && rec.status < 300 {
				cache.Set(key, rec.status, w.Header(), []byte(rec.body.String()))
			}
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	body   *strings.Builder
}

func (rr *responseRecorder) WriteHeader(status int) {
	rr.status = status
	rr.ResponseWriter.WriteHeader(status)
}

func (rr *responseRecorder) Write(p []byte) (int, error) {
	rr.body.Write(p)
	return rr.ResponseWriter.Write(p)
}
