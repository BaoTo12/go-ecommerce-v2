package application

import (
	"context"
	"sync"

	"github.com/titan-commerce/backend/auction-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

// AuctionService handles auction operations
type AuctionService struct {
	repo        domain.Repository
	logger      *logger.Logger
	subscribers map[string][]chan *BidEvent
	mu          sync.RWMutex
}

// BidEvent is sent to WebSocket subscribers
type BidEvent struct {
	AuctionID    string  `json:"auction_id"`
	CurrentBid   float64 `json:"current_bid"`
	BidCount     int     `json:"bid_count"`
	BidderName   string  `json:"bidder_name"`
	TimeRemaining string `json:"time_remaining"`
}

// NewAuctionService creates a new auction service
func NewAuctionService(repo domain.Repository, log *logger.Logger) *AuctionService {
	return &AuctionService{
		repo:        repo,
		logger:      log,
		subscribers: make(map[string][]chan *BidEvent),
	}
}

// GetActiveAuctions returns all active auctions
func (s *AuctionService) GetActiveAuctions(ctx context.Context) ([]*domain.Auction, error) {
	auctions, err := s.repo.GetActiveAuctions(ctx)
	if err != nil {
		return nil, err
	}

	// Update status for each auction
	for _, a := range auctions {
		a.UpdateStatus()
	}

	return auctions, nil
}

// GetAuctionByID returns a specific auction
func (s *AuctionService) GetAuctionByID(ctx context.Context, id string) (*domain.Auction, error) {
	auction, err := s.repo.GetAuctionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if auction != nil {
		auction.UpdateStatus()
	}
	return auction, nil
}

// PlaceBid places a bid on an auction
func (s *AuctionService) PlaceBid(ctx context.Context, auctionID, userID, userName, avatar string, amount float64) (*domain.Bid, error) {
	auction, err := s.repo.GetAuctionByID(ctx, auctionID)
	if err != nil {
		return nil, err
	}
	if auction == nil {
		return nil, domain.ErrAuctionEnded
	}

	bid, err := auction.PlaceBid(userID, userName, avatar, amount)
	if err != nil {
		return nil, err
	}

	// Save updated auction
	if err := s.repo.SaveAuction(ctx, auction); err != nil {
		return nil, err
	}

	// Save bid
	if err := s.repo.SaveBid(ctx, bid); err != nil {
		return nil, err
	}

	// Notify subscribers
	s.notifySubscribers(auctionID, &BidEvent{
		AuctionID:     auctionID,
		CurrentBid:    auction.CurrentBid,
		BidCount:      auction.BidCount,
		BidderName:    userName,
		TimeRemaining: formatDuration(auction.TimeRemaining()),
	})

	s.logger.Infof("Bid placed: auction=%s user=%s amount=%.0f", auctionID, userID, amount)
	return bid, nil
}

// GetBidHistory returns bid history for an auction
func (s *AuctionService) GetBidHistory(ctx context.Context, auctionID string, limit int) ([]*domain.Bid, error) {
	return s.repo.GetBidsByAuction(ctx, auctionID, limit)
}

// Subscribe adds a subscriber for auction updates
func (s *AuctionService) Subscribe(auctionID string) chan *BidEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan *BidEvent, 10)
	s.subscribers[auctionID] = append(s.subscribers[auctionID], ch)
	return ch
}

// Unsubscribe removes a subscriber
func (s *AuctionService) Unsubscribe(auctionID string, ch chan *BidEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	subs := s.subscribers[auctionID]
	for i, sub := range subs {
		if sub == ch {
			s.subscribers[auctionID] = append(subs[:i], subs[i+1:]...)
			close(ch)
			break
		}
	}
}

func (s *AuctionService) notifySubscribers(auctionID string, event *BidEvent) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, ch := range s.subscribers[auctionID] {
		select {
		case ch <- event:
		default:
			// Channel full, skip
		}
	}
}

func formatDuration(d interface{}) string {
	// Type assertion for time.Duration
	if dur, ok := d.(interface{ Hours() float64 }); ok {
		h := int(dur.Hours())
		m := int(dur.(interface{ Minutes() float64 }).Minutes()) % 60
		s := int(dur.(interface{ Seconds() float64 }).Seconds()) % 60
		return formatTime(h, m, s)
	}
	return "00:00:00"
}

func formatTime(h, m, s int) string {
	return pad(h) + ":" + pad(m) + ":" + pad(s)
}

func pad(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}
