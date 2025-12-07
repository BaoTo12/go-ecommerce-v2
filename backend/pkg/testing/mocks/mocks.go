package mocks

import (
	"sync"
	"time"
)

// MockAuctionRepository is a mock implementation for testing
type MockAuctionRepository struct {
	auctions   map[string]*MockAuction
	bids       []*MockBid
	mu         sync.RWMutex

	// Call tracking for verification
	Calls      []MethodCall
	callsMu    sync.Mutex

	// Error injection
	GetActiveError   error
	GetByIDError     error
	SaveAuctionError error
	SaveBidError     error
}

// MockAuction represents an auction for testing
type MockAuction struct {
	ID          string
	ProductID   string
	ProductName string
	CurrentBid  float64
	Status      string
	EndTime     time.Time
}

// MockBid represents a bid for testing
type MockBid struct {
	ID        string
	AuctionID string
	UserID    string
	Amount    float64
	CreatedAt time.Time
}

// MethodCall tracks method invocations
type MethodCall struct {
	Method string
	Args   []interface{}
	Time   time.Time
}

// NewMockAuctionRepository creates a mock repository
func NewMockAuctionRepository() *MockAuctionRepository {
	return &MockAuctionRepository{
		auctions: make(map[string]*MockAuction),
		bids:     make([]*MockBid, 0),
		Calls:    make([]MethodCall, 0),
	}
}

func (m *MockAuctionRepository) recordCall(method string, args ...interface{}) {
	m.callsMu.Lock()
	defer m.callsMu.Unlock()
	m.Calls = append(m.Calls, MethodCall{
		Method: method,
		Args:   args,
		Time:   time.Now(),
	})
}

// GetCallCount returns number of times a method was called
func (m *MockAuctionRepository) GetCallCount(method string) int {
	m.callsMu.Lock()
	defer m.callsMu.Unlock()
	count := 0
	for _, c := range m.Calls {
		if c.Method == method {
			count++
		}
	}
	return count
}

// ResetCalls clears call history
func (m *MockAuctionRepository) ResetCalls() {
	m.callsMu.Lock()
	defer m.callsMu.Unlock()
	m.Calls = make([]MethodCall, 0)
}

// SetupAuction adds an auction for testing
func (m *MockAuctionRepository) SetupAuction(auction *MockAuction) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.auctions[auction.ID] = auction
}

// GetActiveAuctions returns all active auctions
func (m *MockAuctionRepository) GetActiveAuctions() ([]*MockAuction, error) {
	m.recordCall("GetActiveAuctions")
	if m.GetActiveError != nil {
		return nil, m.GetActiveError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*MockAuction, 0, len(m.auctions))
	for _, a := range m.auctions {
		result = append(result, a)
	}
	return result, nil
}

// GetAuctionByID returns an auction by ID
func (m *MockAuctionRepository) GetAuctionByID(id string) (*MockAuction, error) {
	m.recordCall("GetAuctionByID", id)
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.auctions[id], nil
}

// SaveAuction saves an auction
func (m *MockAuctionRepository) SaveAuction(auction *MockAuction) error {
	m.recordCall("SaveAuction", auction.ID)
	if m.SaveAuctionError != nil {
		return m.SaveAuctionError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.auctions[auction.ID] = auction
	return nil
}

// SaveBid saves a bid
func (m *MockAuctionRepository) SaveBid(bid *MockBid) error {
	m.recordCall("SaveBid", bid.ID)
	if m.SaveBidError != nil {
		return m.SaveBidError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.bids = append(m.bids, bid)
	return nil
}

// GetBidsByAuction returns bids for an auction
func (m *MockAuctionRepository) GetBidsByAuction(auctionID string, limit int) ([]*MockBid, error) {
	m.recordCall("GetBidsByAuction", auctionID, limit)

	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*MockBid, 0)
	for _, b := range m.bids {
		if b.AuctionID == auctionID {
			result = append(result, b)
		}
	}
	return result, nil
}

// MockMetricsRegistry for metrics testing
type MockMetricsRegistry struct {
	Counters   map[string]int64
	Gauges     map[string]int64
	Histograms map[string][]float64
	mu         sync.Mutex
}

// NewMockMetricsRegistry creates a mock registry
func NewMockMetricsRegistry() *MockMetricsRegistry {
	return &MockMetricsRegistry{
		Counters:   make(map[string]int64),
		Gauges:     make(map[string]int64),
		Histograms: make(map[string][]float64),
	}
}

// IncrementCounter increments a counter
func (m *MockMetricsRegistry) IncrementCounter(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Counters[name]++
}

// SetGauge sets a gauge value
func (m *MockMetricsRegistry) SetGauge(name string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Gauges[name] = value
}

// ObserveHistogram records a histogram value
func (m *MockMetricsRegistry) ObserveHistogram(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Histograms[name] = append(m.Histograms[name], value)
}
