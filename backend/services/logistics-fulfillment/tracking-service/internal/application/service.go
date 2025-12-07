package application

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/titan-commerce/backend/tracking-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

// TrackingService handles order tracking operations
type TrackingService struct {
	repo        domain.Repository
	logger      *logger.Logger
	subscribers map[string][]chan *LocationUpdate
	mu          sync.RWMutex
}

// LocationUpdate is sent to WebSocket subscribers
type LocationUpdate struct {
	OrderID  string  `json:"order_id"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	ETA      int     `json:"eta_minutes"`
	Distance float64 `json:"distance_km"`
}

// NewTrackingService creates a new tracking service
func NewTrackingService(repo domain.Repository, log *logger.Logger) *TrackingService {
	svc := &TrackingService{
		repo:        repo,
		logger:      log,
		subscribers: make(map[string][]chan *LocationUpdate),
	}
	// Start location simulation
	go svc.simulateDriverMovement()
	return svc
}

// GetTracking returns tracking information for an order
func (s *TrackingService) GetTracking(ctx context.Context, orderID string) (*domain.OrderTracking, error) {
	return s.repo.GetTrackingByOrderID(ctx, orderID)
}

// UpdateLocation updates the driver's location for an order
func (s *TrackingService) UpdateLocation(ctx context.Context, orderID string, lat, lng float64) error {
	if err := s.repo.UpdateDriverLocation(ctx, orderID, lat, lng); err != nil {
		return err
	}

	// Notify subscribers
	s.notifySubscribers(orderID, &LocationUpdate{
		OrderID:  orderID,
		Lat:      lat,
		Lng:      lng,
		ETA:      rand.Intn(30) + 5,
		Distance: rand.Float64()*5 + 0.5,
	})

	return nil
}

// Subscribe adds a subscriber for location updates
func (s *TrackingService) Subscribe(orderID string) chan *LocationUpdate {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan *LocationUpdate, 10)
	s.subscribers[orderID] = append(s.subscribers[orderID], ch)
	return ch
}

// Unsubscribe removes a subscriber
func (s *TrackingService) Unsubscribe(orderID string, ch chan *LocationUpdate) {
	s.mu.Lock()
	defer s.mu.Unlock()

	subs := s.subscribers[orderID]
	for i, sub := range subs {
		if sub == ch {
			s.subscribers[orderID] = append(subs[:i], subs[i+1:]...)
			close(ch)
			break
		}
	}
}

func (s *TrackingService) notifySubscribers(orderID string, update *LocationUpdate) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, ch := range s.subscribers[orderID] {
		select {
		case ch <- update:
		default:
			// Channel full, skip
		}
	}
}

// simulateDriverMovement simulates driver movement for demo purposes
func (s *TrackingService) simulateDriverMovement() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	baseLat := 10.775
	baseLng := 106.700

	for range ticker.C {
		s.mu.RLock()
		for orderID := range s.subscribers {
			// Simulate movement
			lat := baseLat + (rand.Float64()-0.5)*0.01
			lng := baseLng + (rand.Float64()-0.5)*0.01

			s.notifySubscribers(orderID, &LocationUpdate{
				OrderID:  orderID,
				Lat:      lat,
				Lng:      lng,
				ETA:      rand.Intn(20) + 5,
				Distance: rand.Float64()*3 + 0.2,
			})
		}
		s.mu.RUnlock()
	}
}
