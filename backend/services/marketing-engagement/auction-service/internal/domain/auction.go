package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAuctionEnded        = errors.New("auction has ended")
	ErrAuctionNotStarted   = errors.New("auction has not started")
	ErrBidTooLow           = errors.New("bid must be higher than current bid + minimum increment")
	ErrInsufficientBalance = errors.New("insufficient balance to place bid")
)

type AuctionStatus string

const (
	AuctionPending  AuctionStatus = "pending"
	AuctionActive   AuctionStatus = "active"
	AuctionEnding   AuctionStatus = "ending"   // Last 10 minutes
	AuctionEnded    AuctionStatus = "ended"
	AuctionCanceled AuctionStatus = "canceled"
)

// Auction represents an auction item
type Auction struct {
	ID             string        `json:"id"`
	ProductID      string        `json:"product_id"`
	ProductName    string        `json:"product_name"`
	ProductImage   string        `json:"product_image"`
	OriginalPrice  float64       `json:"original_price"`
	StartingBid    float64       `json:"starting_bid"`
	CurrentBid     float64       `json:"current_bid"`
	MinIncrement   float64       `json:"min_increment"`
	BidCount       int           `json:"bid_count"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	Status         AuctionStatus `json:"status"`
	HighestBidder  *Bidder       `json:"highest_bidder,omitempty"`
	WinnerID       string        `json:"winner_id,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// Bidder represents a user who placed a bid
type Bidder struct {
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// Bid represents a single bid on an auction
type Bid struct {
	ID        string    `json:"id"`
	AuctionID string    `json:"auction_id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	IsWinning bool      `json:"is_winning"`
}

// NewAuction creates a new auction
func NewAuction(productID, name, image string, originalPrice, startingBid float64, duration time.Duration) *Auction {
	now := time.Now()
	return &Auction{
		ID:            uuid.New().String(),
		ProductID:     productID,
		ProductName:   name,
		ProductImage:  image,
		OriginalPrice: originalPrice,
		StartingBid:   startingBid,
		CurrentBid:    startingBid,
		MinIncrement:  100000, // 100,000 VND minimum increment
		BidCount:      0,
		StartTime:     now,
		EndTime:       now.Add(duration),
		Status:        AuctionActive,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// UpdateStatus updates the auction status based on current time
func (a *Auction) UpdateStatus() {
	now := time.Now()

	if now.Before(a.StartTime) {
		a.Status = AuctionPending
		return
	}

	if now.After(a.EndTime) {
		a.Status = AuctionEnded
		if a.HighestBidder != nil {
			a.WinnerID = a.HighestBidder.UserID
		}
		return
	}

	// Last 10 minutes = ending soon
	if a.EndTime.Sub(now) < 10*time.Minute {
		a.Status = AuctionEnding
		return
	}

	a.Status = AuctionActive
}

// PlaceBid attempts to place a bid on the auction
func (a *Auction) PlaceBid(userID, userName, avatar string, amount float64) (*Bid, error) {
	now := time.Now()

	// Check if auction is active
	if now.Before(a.StartTime) {
		return nil, ErrAuctionNotStarted
	}
	if now.After(a.EndTime) {
		return nil, ErrAuctionEnded
	}

	// Check bid amount
	minBid := a.CurrentBid + a.MinIncrement
	if amount < minBid {
		return nil, ErrBidTooLow
	}

	// Create bid
	bid := &Bid{
		ID:        uuid.New().String(),
		AuctionID: a.ID,
		UserID:    userID,
		Amount:    amount,
		CreatedAt: now,
		IsWinning: true,
	}

	// Update auction
	a.CurrentBid = amount
	a.BidCount++
	a.HighestBidder = &Bidder{
		UserID:    userID,
		Name:      userName,
		AvatarURL: avatar,
	}
	a.UpdatedAt = now
	a.UpdateStatus()

	// Extend auction by 2 minutes if bid placed in last 2 minutes
	if a.EndTime.Sub(now) < 2*time.Minute {
		a.EndTime = now.Add(2 * time.Minute)
	}

	return bid, nil
}

// TimeRemaining returns the time remaining until auction ends
func (a *Auction) TimeRemaining() time.Duration {
	remaining := a.EndTime.Sub(time.Now())
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Repository interface for auction data
type Repository interface {
	GetActiveAuctions(ctx interface{}) ([]*Auction, error)
	GetAuctionByID(ctx interface{}, id string) (*Auction, error)
	SaveAuction(ctx interface{}, auction *Auction) error
	SaveBid(ctx interface{}, bid *Bid) error
	GetBidsByAuction(ctx interface{}, auctionID string, limit int) ([]*Bid, error)
	GetUserBids(ctx interface{}, userID string) ([]*Bid, error)
}
