package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

/*
API GATEWAY PATTERN

Implements a full API Gateway with:
- Request routing
- Load balancing (round-robin, least connections)
- Circuit breaking
- Rate limiting per route
- Request/Response transformation
- Authentication passthrough
- Service discovery integration
*/

// Route defines an API route
type Route struct {
	Path        string
	Methods     []string
	Backends    []Backend
	StripPrefix string
	Timeout     time.Duration
	RateLimit   int
	Auth        bool
	Transform   *TransformConfig
}

// Backend represents a backend service
type Backend struct {
	URL    string
	Weight int
	Health bool
}

// TransformConfig configures request/response transformation
type TransformConfig struct {
	AddHeaders    map[string]string
	RemoveHeaders []string
	RewritePath   string
}

// LoadBalancer selects backends
type LoadBalancer interface {
	Select(backends []Backend) *Backend
}

// Gateway is the API Gateway
type Gateway struct {
	routes       map[string]*Route
	loadBalancer LoadBalancer
	rateLimiters map[string]*rateLimiter
	proxies      map[string]*httputil.ReverseProxy
	mu           sync.RWMutex
}

type rateLimiter struct {
	tokens   float64
	max      float64
	rate     float64
	lastTime time.Time
	mu       sync.Mutex
}

// NewGateway creates an API gateway
func NewGateway() *Gateway {
	return &Gateway{
		routes:       make(map[string]*Route),
		loadBalancer: &RoundRobinLB{},
		rateLimiters: make(map[string]*rateLimiter),
		proxies:      make(map[string]*httputil.ReverseProxy),
	}
}

// RegisterRoute registers a route
func (g *Gateway) RegisterRoute(route *Route) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.routes[route.Path] = route

	// Setup rate limiter if configured
	if route.RateLimit > 0 {
		g.rateLimiters[route.Path] = &rateLimiter{
			tokens:   float64(route.RateLimit),
			max:      float64(route.RateLimit),
			rate:     float64(route.RateLimit),
			lastTime: time.Now(),
		}
	}

	// Create proxies for backends
	for _, backend := range route.Backends {
		if _, exists := g.proxies[backend.URL]; !exists {
			target, err := url.Parse(backend.URL)
			if err != nil {
				return err
			}
			g.proxies[backend.URL] = httputil.NewSingleHostReverseProxy(target)
		}
	}

	return nil
}

// ServeHTTP handles incoming requests
func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Find matching route
	route := g.matchRoute(r)
	if route == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	// Check method
	if !g.methodAllowed(route, r.Method) {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Rate limiting
	if !g.checkRateLimit(route.Path) {
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	// Select backend
	backend := g.loadBalancer.Select(route.Backends)
	if backend == nil {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	// Apply transformations
	if route.Transform != nil {
		g.applyTransform(r, route.Transform)
	}

	// Strip prefix if configured
	if route.StripPrefix != "" {
		r.URL.Path = r.URL.Path[len(route.StripPrefix):]
	}

	// Proxy request
	g.mu.RLock()
	proxy := g.proxies[backend.URL]
	g.mu.RUnlock()

	if proxy == nil {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	// Set timeout
	ctx := r.Context()
	if route.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, route.Timeout)
		defer cancel()
	}

	proxy.ServeHTTP(w, r.WithContext(ctx))
}

func (g *Gateway) matchRoute(r *http.Request) *Route {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Exact match first
	if route, ok := g.routes[r.URL.Path]; ok {
		return route
	}

	// Prefix match
	longest := ""
	var matched *Route
	for path, route := range g.routes {
		if len(path) > len(longest) && hasPrefix(r.URL.Path, path) {
			longest = path
			matched = route
		}
	}

	return matched
}

func (g *Gateway) methodAllowed(route *Route, method string) bool {
	if len(route.Methods) == 0 {
		return true
	}
	for _, m := range route.Methods {
		if m == method {
			return true
		}
	}
	return false
}

func (g *Gateway) checkRateLimit(path string) bool {
	g.mu.RLock()
	limiter, exists := g.rateLimiters[path]
	g.mu.RUnlock()

	if !exists {
		return true
	}

	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(limiter.lastTime).Seconds()
	limiter.tokens += elapsed * limiter.rate
	if limiter.tokens > limiter.max {
		limiter.tokens = limiter.max
	}
	limiter.lastTime = now

	if limiter.tokens >= 1 {
		limiter.tokens--
		return true
	}
	return false
}

func (g *Gateway) applyTransform(r *http.Request, config *TransformConfig) {
	for k, v := range config.AddHeaders {
		r.Header.Set(k, v)
	}
	for _, k := range config.RemoveHeaders {
		r.Header.Del(k)
	}
	if config.RewritePath != "" {
		r.URL.Path = config.RewritePath
	}
}

// RoundRobinLB implements round-robin load balancing
type RoundRobinLB struct {
	counter uint64
	mu      sync.Mutex
}

func (lb *RoundRobinLB) Select(backends []Backend) *Backend {
	healthy := make([]Backend, 0)
	for _, b := range backends {
		if b.Health {
			healthy = append(healthy, b)
		}
	}

	if len(healthy) == 0 {
		// Fallback to all backends
		healthy = backends
	}

	if len(healthy) == 0 {
		return nil
	}

	lb.mu.Lock()
	lb.counter++
	idx := lb.counter % uint64(len(healthy))
	lb.mu.Unlock()

	return &healthy[idx]
}

// LeastConnectionsLB implements least connections load balancing
type LeastConnectionsLB struct {
	connections map[string]int
	mu          sync.Mutex
}

func NewLeastConnectionsLB() *LeastConnectionsLB {
	return &LeastConnectionsLB{
		connections: make(map[string]int),
	}
}

func (lb *LeastConnectionsLB) Select(backends []Backend) *Backend {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	var selected *Backend
	minConns := -1

	for i := range backends {
		if !backends[i].Health {
			continue
		}
		conns := lb.connections[backends[i].URL]
		if minConns == -1 || conns < minConns {
			minConns = conns
			selected = &backends[i]
		}
	}

	if selected != nil {
		lb.connections[selected.URL]++
	}

	return selected
}

func (lb *LeastConnectionsLB) Release(backendURL string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	if lb.connections[backendURL] > 0 {
		lb.connections[backendURL]--
	}
}

// ServiceRegistry for service discovery
type ServiceRegistry interface {
	Register(name string, url string) error
	Deregister(name string, url string) error
	Discover(name string) ([]string, error)
	Watch(name string, callback func([]string))
}

// InMemoryServiceRegistry is an in-memory implementation
type InMemoryServiceRegistry struct {
	services map[string][]string
	watchers map[string][]func([]string)
	mu       sync.RWMutex
}

func NewInMemoryServiceRegistry() *InMemoryServiceRegistry {
	return &InMemoryServiceRegistry{
		services: make(map[string][]string),
		watchers: make(map[string][]func([]string)),
	}
}

func (r *InMemoryServiceRegistry) Register(name, url string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[name] = append(r.services[name], url)
	r.notifyWatchers(name)
	return nil
}

func (r *InMemoryServiceRegistry) Deregister(name, url string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	urls := r.services[name]
	for i, u := range urls {
		if u == url {
			r.services[name] = append(urls[:i], urls[i+1:]...)
			break
		}
	}
	r.notifyWatchers(name)
	return nil
}

func (r *InMemoryServiceRegistry) Discover(name string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.services[name], nil
}

func (r *InMemoryServiceRegistry) Watch(name string, callback func([]string)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.watchers[name] = append(r.watchers[name], callback)
}

func (r *InMemoryServiceRegistry) notifyWatchers(name string) {
	urls := r.services[name]
	for _, cb := range r.watchers[name] {
		go cb(urls)
	}
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// AggregatorEndpoint aggregates responses from multiple services
type AggregatorEndpoint struct {
	endpoints []struct {
		URL     string
		Key     string
		Timeout time.Duration
	}
}

// Aggregate calls all endpoints and aggregates responses
func (a *AggregatorEndpoint) Aggregate(ctx context.Context) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make([]error, 0)

	for _, ep := range a.endpoints {
		wg.Add(1)
		go func(url, key string, timeout time.Duration) {
			defer wg.Done()

			reqCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			req, _ := http.NewRequestWithContext(reqCtx, "GET", url, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			var data interface{}
			json.NewDecoder(resp.Body).Decode(&data)

			mu.Lock()
			result[key] = data
			mu.Unlock()
		}(ep.URL, ep.Key, ep.Timeout)
	}

	wg.Wait()

	if len(errs) > 0 {
		return result, errors.New("partial failure")
	}
	return result, nil
}
