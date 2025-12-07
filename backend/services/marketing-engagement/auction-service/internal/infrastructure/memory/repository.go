package memory

import (
	"sync"
	"time"

	"github.com/titan-commerce/backend/auction-service/internal/domain"
)

// AuctionRepository is an in-memory implementation for development
type AuctionRepository struct {
	auctions map[string]*domain.Auction
	bids     []*domain.Bid
	mu       sync.RWMutex
}

// NewAuctionRepository creates a new in-memory repository with sample data
func NewAuctionRepository() *AuctionRepository {
	repo := &AuctionRepository{
		auctions: make(map[string]*domain.Auction),
		bids:     make([]*domain.Bid, 0),
	}
	repo.seedAuctions()
	return repo
}

func (r *AuctionRepository) seedAuctions() {
	now := time.Now()

	auctions := []*domain.Auction{
		{
			ID:            "a1",
			ProductID:     "p1",
			ProductName:   "iPhone 15 Pro Max 512GB Limited Edition",
			ProductImage:  "https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=400",
			OriginalPrice: 35990000,
			StartingBid:   20000000,
			CurrentBid:    25500000,
			MinIncrement:  100000,
			BidCount:      47,
			StartTime:     now.Add(-1 * time.Hour),
			EndTime:       now.Add(2 * time.Hour),
			Status:        domain.AuctionActive,
			HighestBidder: &domain.Bidder{UserID: "u1", Name: "Nguyễn V***", AvatarURL: "https://ui-avatars.com/api/?name=NV"},
			CreatedAt:     now.Add(-2 * time.Hour),
			UpdatedAt:     now,
		},
		{
			ID:            "a2",
			ProductID:     "p5",
			ProductName:   "Nike Air Jordan 1 Retro High OG Limited",
			ProductImage:  "https://images.unsplash.com/photo-1600269452121-4f2416e55c28?w=400",
			OriginalPrice: 8990000,
			StartingBid:   3000000,
			CurrentBid:    5200000,
			MinIncrement:  100000,
			BidCount:      32,
			StartTime:     now.Add(-30 * time.Minute),
			EndTime:       now.Add(45 * time.Minute),
			Status:        domain.AuctionEnding,
			HighestBidder: &domain.Bidder{UserID: "u2", Name: "Trần T***", AvatarURL: "https://ui-avatars.com/api/?name=TT"},
			CreatedAt:     now.Add(-1 * time.Hour),
			UpdatedAt:     now,
		},
		{
			ID:            "a3",
			ProductID:     "p3",
			ProductName:   "MacBook Pro 14\" M3 Pro Chip",
			ProductImage:  "https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=400",
			OriginalPrice: 52990000,
			StartingBid:   35000000,
			CurrentBid:    42000000,
			MinIncrement:  500000,
			BidCount:      28,
			StartTime:     now.Add(-2 * time.Hour),
			EndTime:       now.Add(5 * time.Hour),
			Status:        domain.AuctionActive,
			HighestBidder: &domain.Bidder{UserID: "u3", Name: "Lê H***", AvatarURL: "https://ui-avatars.com/api/?name=LH"},
			CreatedAt:     now.Add(-3 * time.Hour),
			UpdatedAt:     now,
		},
	}

	for _, a := range auctions {
		r.auctions[a.ID] = a
	}
}

func (r *AuctionRepository) GetActiveAuctions(ctx interface{}) ([]*domain.Auction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Auction, 0)
	now := time.Now()

	for _, a := range r.auctions {
		if a.EndTime.After(now) {
			a.UpdateStatus()
			result = append(result, a)
		}
	}

	return result, nil
}

func (r *AuctionRepository) GetAuctionByID(ctx interface{}, id string) (*domain.Auction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if a, ok := r.auctions[id]; ok {
		return a, nil
	}
	return nil, nil
}

func (r *AuctionRepository) SaveAuction(ctx interface{}, auction *domain.Auction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.auctions[auction.ID] = auction
	return nil
}

func (r *AuctionRepository) SaveBid(ctx interface{}, bid *domain.Bid) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.bids = append(r.bids, bid)
	return nil
}

func (r *AuctionRepository) GetBidsByAuction(ctx interface{}, auctionID string, limit int) ([]*domain.Bid, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Bid, 0)
	for i := len(r.bids) - 1; i >= 0; i-- {
		if r.bids[i].AuctionID == auctionID {
			result = append(result, r.bids[i])
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (r *AuctionRepository) GetUserBids(ctx interface{}, userID string) ([]*domain.Bid, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Bid, 0)
	for _, b := range r.bids {
		if b.UserID == userID {
			result = append(result, b)
		}
	}
	return result, nil
}
