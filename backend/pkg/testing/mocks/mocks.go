package mocks

import (
	"context"
	"sync"
	"time"

	"github.com/titan-commerce/backend/auction-service/internal/domain"
)

// MockAuctionRepository is a mock implementation for testing
type MockAuctionRepository struct {
	auctions   map[string]*domain.Auction
	bids       []*domain.Bid
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

// MethodCall tracks method invocations
type MethodCall struct {
	Method string
	Args   []interface{}
	Time   time.Time
}

// NewMockAuctionRepository creates a mock repository
func NewMockAuctionRepository() *MockAuctionRepository {
	return &MockAuctionRepository{
		auctions: make(map[string]*domain.Auction),
		bids:     make([]*domain.Bid, 0),
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
func (m *MockAuctionRepository) SetupAuction(auction *domain.Auction) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.auctions[auction.ID] = auction
}

// Interface implementations

func (m *MockAuctionRepository) GetActiveAuctions(ctx interface{}) ([]*domain.Auction, error) {
	m.recordCall("GetActiveAuctions")
	if m.GetActiveError != nil {
		return nil, m.GetActiveError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*domain.Auction, 0, len(m.auctions))
	for _, a := range m.auctions {
		result = append(result, a)
	}
	return result, nil
}

func (m *MockAuctionRepository) GetAuctionByID(ctx interface{}, id string) (*domain.Auction, error) {
	m.recordCall("GetAuctionByID", id)
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.auctions[id], nil
}

func (m *MockAuctionRepository) SaveAuction(ctx interface{}, auction *domain.Auction) error {
	m.recordCall("SaveAuction", auction.ID)
	if m.SaveAuctionError != nil {
		return m.SaveAuctionError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.auctions[auction.ID] = auction
	return nil
}

func (m *MockAuctionRepository) SaveBid(ctx interface{}, bid *domain.Bid) error {
	m.recordCall("SaveBid", bid.ID)
	if m.SaveBidError != nil {
		return m.SaveBidError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.bids = append(m.bids, bid)
	return nil
}

func (m *MockAuctionRepository) GetBidsByAuction(ctx interface{}, auctionID string, limit int) ([]*domain.Bid, error) {
	m.recordCall("GetBidsByAuction", auctionID, limit)

	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*domain.Bid, 0)
	for _, b := range m.bids {
		if b.AuctionID == auctionID {
			result = append(result, b)
		}
	}
	return result, nil
}

func (m *MockAuctionRepository) GetUserBids(ctx interface{}, userID string) ([]*domain.Bid, error) {
	m.recordCall("GetUserBids", userID)

	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*domain.Bid, 0)
	for _, b := range m.bids {
		if b.UserID == userID {
			result = append(result, b)
		}
	}
	return result, nil
}

// MockEventStore for CQRS testing
type MockEventStore struct {
	events     map[string][]interface{}
	mu         sync.RWMutex
	SaveError  error
	LoadError  error
}

func NewMockEventStore() *MockEventStore {
	return &MockEventStore{
		events: make(map[string][]interface{}),
	}
}

func (m *MockEventStore) Save(ctx context.Context, events ...interface{}) error {
	if m.SaveError != nil {
		return m.SaveError
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	// Simple append for testing
	m.events["all"] = append(m.events["all"], events...)
	return nil
}

func (m *MockEventStore) Load(ctx context.Context, aggregateID string) ([]interface{}, error) {
	if m.LoadError != nil {
		return nil, m.LoadError
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.events[aggregateID], nil
}

// MockMetricsRegistry for metrics testing
type MockMetricsRegistry struct {
	Counters   map[string]int64
	Gauges     map[string]int64
	Histograms map[string][]float64
	mu         sync.Mutex
}

func NewMockMetricsRegistry() *MockMetricsRegistry {
	return &MockMetricsRegistry{
		Counters:   make(map[string]int64),
		Gauges:     make(map[string]int64),
		Histograms: make(map[string][]float64),
	}
}

func (m *MockMetricsRegistry) IncrementCounter(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Counters[name]++
}

func (m *MockMetricsRegistry) SetGauge(name string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Gauges[name] = value
}

func (m *MockMetricsRegistry) ObserveHistogram(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Histograms[name] = append(m.Histograms[name], value)
}
