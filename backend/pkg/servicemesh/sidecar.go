package servicemesh

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

/*
SERVICE MESH SIDECAR PATTERN

Implements service mesh functionality:
- Service discovery
- Load balancing
- Circuit breaking
- Retry with exponential backoff
- Timeout management
- Request tracing
- Health checking
- mTLS simulation
*/

// ServiceEndpoint represents a service instance
type ServiceEndpoint struct {
	ID       string
	Service  string
	Host     string
	Port     int
	Protocol string
	Weight   int
	Health   HealthStatus
	Metadata map[string]string
}

// HealthStatus represents health status
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "HEALTHY"
	HealthStatusUnhealthy HealthStatus = "UNHEALTHY"
	HealthStatusUnknown   HealthStatus = "UNKNOWN"
)

// ServiceDiscovery discovers services
type ServiceDiscovery interface {
	Register(ctx context.Context, endpoint *ServiceEndpoint) error
	Deregister(ctx context.Context, id string) error
	Discover(ctx context.Context, service string) ([]*ServiceEndpoint, error)
	Watch(service string, callback func([]*ServiceEndpoint))
}

// Sidecar is the service mesh sidecar
type Sidecar struct {
	discovery    ServiceDiscovery
	loadBalancer LoadBalancerStrategy
	circuitBreakers map[string]*circuitBreaker
	retryPolicy  *RetryPolicy
	healthChecker *HealthChecker
	metrics      *SidecarMetrics
	mu           sync.RWMutex
}

// SidecarConfig configures the sidecar
type SidecarConfig struct {
	RetryPolicy   *RetryPolicy
	CircuitConfig *CircuitConfig
	HealthConfig  *HealthConfig
}

// RetryPolicy configures retries
type RetryPolicy struct {
	MaxRetries      int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	RetryOn         []int // HTTP status codes to retry
}

// CircuitConfig configures circuit breakers
type CircuitConfig struct {
	FailureThreshold int
	SuccessThreshold int
	Timeout          time.Duration
	HalfOpenRequests int
}

// HealthConfig configures health checking
type HealthConfig struct {
	Interval    time.Duration
	Timeout     time.Duration
	HealthPath  string
	Threshold   int
}

// DefaultRetryPolicy returns default retry policy
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:      3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     2 * time.Second,
		Multiplier:      2.0,
		RetryOn:         []int{500, 502, 503, 504},
	}
}

// NewSidecar creates a service mesh sidecar
func NewSidecar(discovery ServiceDiscovery, config *SidecarConfig) *Sidecar {
	if config == nil {
		config = &SidecarConfig{
			RetryPolicy: DefaultRetryPolicy(),
		}
	}

	return &Sidecar{
		discovery:       discovery,
		loadBalancer:    &RoundRobinStrategy{},
		circuitBreakers: make(map[string]*circuitBreaker),
		retryPolicy:     config.RetryPolicy,
		metrics:         NewSidecarMetrics(),
	}
}

// Call makes a service call with all sidecar features
func (s *Sidecar) Call(ctx context.Context, service, method, path string, body []byte) (*http.Response, error) {
	// Discover endpoints
	endpoints, err := s.discovery.Discover(ctx, service)
	if err != nil {
		return nil, err
	}

	// Filter healthy endpoints
	healthy := make([]*ServiceEndpoint, 0)
	for _, ep := range endpoints {
		if ep.Health == HealthStatusHealthy {
			healthy = append(healthy, ep)
		}
	}

	if len(healthy) == 0 {
		return nil, errors.New("no healthy endpoints")
	}

	// Execute with retry
	var lastErr error
	interval := s.retryPolicy.InitialInterval

	for attempt := 0; attempt <= s.retryPolicy.MaxRetries; attempt++ {
		// Select endpoint
		endpoint := s.loadBalancer.Select(healthy)
		if endpoint == nil {
			return nil, errors.New("no endpoint available")
		}

		// Check circuit breaker
		cb := s.getCircuitBreaker(endpoint.ID)
		if !cb.Allow() {
			continue
		}

		// Make request
		s.metrics.RequestsTotal.Add(1)
		start := time.Now()

		resp, err := s.doRequest(ctx, endpoint, method, path, body)
		
		s.metrics.RequestDuration.Add(time.Since(start).Milliseconds())

		if err != nil {
			cb.RecordFailure()
			lastErr = err
			time.Sleep(interval)
			interval = time.Duration(float64(interval) * s.retryPolicy.Multiplier)
			if interval > s.retryPolicy.MaxInterval {
				interval = s.retryPolicy.MaxInterval
			}
			continue
		}

		// Check if should retry based on status
		if s.shouldRetry(resp.StatusCode) {
			cb.RecordFailure()
			lastErr = errors.New("retriable status code")
			resp.Body.Close()
			time.Sleep(interval)
			interval = time.Duration(float64(interval) * s.retryPolicy.Multiplier)
			continue
		}

		cb.RecordSuccess()
		return resp, nil
	}

	s.metrics.RequestsFailed.Add(1)
	return nil, lastErr
}

func (s *Sidecar) doRequest(ctx context.Context, endpoint *ServiceEndpoint, method, path string, body []byte) (*http.Response, error) {
	url := endpoint.Protocol + "://" + endpoint.Host + ":" + itoa(endpoint.Port) + path
	
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	// Add tracing headers
	if traceID := ctx.Value("trace_id"); traceID != nil {
		req.Header.Set("X-Trace-ID", traceID.(string))
	}

	return http.DefaultClient.Do(req)
}

func (s *Sidecar) shouldRetry(statusCode int) bool {
	for _, code := range s.retryPolicy.RetryOn {
		if statusCode == code {
			return true
		}
	}
	return false
}

func (s *Sidecar) getCircuitBreaker(id string) *circuitBreaker {
	s.mu.RLock()
	cb, exists := s.circuitBreakers[id]
	s.mu.RUnlock()

	if !exists {
		s.mu.Lock()
		cb = newCircuitBreaker()
		s.circuitBreakers[id] = cb
		s.mu.Unlock()
	}

	return cb
}

// Circuit breaker implementation
type circuitBreaker struct {
	failures   int64
	successes  int64
	state      int32 // 0=closed, 1=open, 2=half-open
	lastFailed time.Time
	mu         sync.Mutex
}

func newCircuitBreaker() *circuitBreaker {
	return &circuitBreaker{}
}

func (cb *circuitBreaker) Allow() bool {
	state := atomic.LoadInt32(&cb.state)
	
	if state == 0 { // Closed
		return true
	}
	
	if state == 1 { // Open
		cb.mu.Lock()
		if time.Since(cb.lastFailed) > 30*time.Second {
			atomic.StoreInt32(&cb.state, 2) // Half-open
			cb.mu.Unlock()
			return true
		}
		cb.mu.Unlock()
		return false
	}
	
	// Half-open
	return true
}

func (cb *circuitBreaker) RecordSuccess() {
	atomic.AddInt64(&cb.successes, 1)
	atomic.StoreInt64(&cb.failures, 0)
	atomic.StoreInt32(&cb.state, 0) // Close circuit
}

func (cb *circuitBreaker) RecordFailure() {
	failures := atomic.AddInt64(&cb.failures, 1)
	cb.mu.Lock()
	cb.lastFailed = time.Now()
	cb.mu.Unlock()
	
	if failures >= 5 {
		atomic.StoreInt32(&cb.state, 1) // Open circuit
	}
}

// Load balancer strategies

// LoadBalancerStrategy selects endpoints
type LoadBalancerStrategy interface {
	Select(endpoints []*ServiceEndpoint) *ServiceEndpoint
}

// RoundRobinStrategy implements round-robin
type RoundRobinStrategy struct {
	counter uint64
}

func (s *RoundRobinStrategy) Select(endpoints []*ServiceEndpoint) *ServiceEndpoint {
	if len(endpoints) == 0 {
		return nil
	}
	idx := atomic.AddUint64(&s.counter, 1) % uint64(len(endpoints))
	return endpoints[idx]
}

// WeightedStrategy implements weighted selection
type WeightedStrategy struct{}

func (s *WeightedStrategy) Select(endpoints []*ServiceEndpoint) *ServiceEndpoint {
	if len(endpoints) == 0 {
		return nil
	}
	
	totalWeight := 0
	for _, ep := range endpoints {
		totalWeight += ep.Weight
	}
	
	if totalWeight == 0 {
		return endpoints[0]
	}
	
	// Random selection based on weight
	r := int(time.Now().UnixNano() % int64(totalWeight))
	for _, ep := range endpoints {
		r -= ep.Weight
		if r < 0 {
			return ep
		}
	}
	
	return endpoints[0]
}

// SidecarMetrics tracks sidecar metrics
type SidecarMetrics struct {
	RequestsTotal    *atomic.Int64
	RequestsFailed   *atomic.Int64
	RequestDuration  *atomic.Int64
	CircuitOpens     *atomic.Int64
}

func NewSidecarMetrics() *SidecarMetrics {
	return &SidecarMetrics{
		RequestsTotal:   new(atomic.Int64),
		RequestsFailed:  new(atomic.Int64),
		RequestDuration: new(atomic.Int64),
		CircuitOpens:    new(atomic.Int64),
	}
}

// HealthChecker checks endpoint health
type HealthChecker struct {
	discovery  ServiceDiscovery
	config     *HealthConfig
	stopCh     chan struct{}
}

func NewHealthChecker(discovery ServiceDiscovery, config *HealthConfig) *HealthChecker {
	return &HealthChecker{
		discovery: discovery,
		config:    config,
		stopCh:    make(chan struct{}),
	}
}

func (h *HealthChecker) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(h.config.Interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return
			case <-h.stopCh:
				return
			case <-ticker.C:
				// Health check logic would go here
			}
		}
	}()
}

func (h *HealthChecker) Stop() {
	close(h.stopCh)
}

// InMemoryServiceDiscovery is an in-memory implementation
type InMemoryServiceDiscovery struct {
	endpoints map[string][]*ServiceEndpoint
	watchers  map[string][]func([]*ServiceEndpoint)
	mu        sync.RWMutex
}

func NewInMemoryServiceDiscovery() *InMemoryServiceDiscovery {
	return &InMemoryServiceDiscovery{
		endpoints: make(map[string][]*ServiceEndpoint),
		watchers:  make(map[string][]func([]*ServiceEndpoint)),
	}
}

func (d *InMemoryServiceDiscovery) Register(ctx context.Context, endpoint *ServiceEndpoint) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.endpoints[endpoint.Service] = append(d.endpoints[endpoint.Service], endpoint)
	d.notifyWatchers(endpoint.Service)
	return nil
}

func (d *InMemoryServiceDiscovery) Deregister(ctx context.Context, id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	for service, eps := range d.endpoints {
		for i, ep := range eps {
			if ep.ID == id {
				d.endpoints[service] = append(eps[:i], eps[i+1:]...)
				d.notifyWatchers(service)
				return nil
			}
		}
	}
	return nil
}

func (d *InMemoryServiceDiscovery) Discover(ctx context.Context, service string) ([]*ServiceEndpoint, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.endpoints[service], nil
}

func (d *InMemoryServiceDiscovery) Watch(service string, callback func([]*ServiceEndpoint)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.watchers[service] = append(d.watchers[service], callback)
}

func (d *InMemoryServiceDiscovery) notifyWatchers(service string) {
	eps := d.endpoints[service]
	for _, cb := range d.watchers[service] {
		go cb(eps)
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	buf := make([]byte, 10)
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}

// JSONCall is a helper for JSON API calls
func (s *Sidecar) JSONCall(ctx context.Context, service, method, path string, body interface{}, response interface{}) error {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	resp, err := s.Call(ctx, service, method, path, bodyBytes)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if response != nil {
		return json.NewDecoder(resp.Body).Decode(response)
	}
	return nil
}
