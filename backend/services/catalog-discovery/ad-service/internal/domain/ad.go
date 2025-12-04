package domain

import (
	"time"

	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/google/uuid"
)

type Campaign struct {
	ID              string
	SellerID        string
	ProductID       string
	Budget          float64
	RemainingBudget float64
	BidAmount       float64
	Status          string
	StartTime       time.Time
	EndTime         time.Time
	CreatedAt       time.Time
}

type AdEvent struct {
	ID        string
	AdID      string
	UserID    string
	EventType string
	Timestamp time.Time
}

func NewCampaign(sellerID, productID string, budget, bidAmount float64, start, end time.Time) (*Campaign, error) {
	if budget <= 0 || bidAmount <= 0 {
		return nil, errors.New(errors.ErrInvalidInput, "budget and bid amount must be positive")
	}

	return &Campaign{
		ID:              uuid.New().String(),
		SellerID:        sellerID,
		ProductID:       productID,
		Budget:          budget,
		RemainingBudget: budget,
		BidAmount:       bidAmount,
		Status:          "active",
		StartTime:       start,
		EndTime:         end,
		CreatedAt:       time.Now(),
	}, nil
}
