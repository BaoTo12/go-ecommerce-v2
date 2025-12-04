package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserPoints represents a user's points balance
type UserPoints struct {
	UserID         string
	TotalPoints    int
	AvailablePoints int
	LifetimeEarned int
	LifetimeSpent  int
	Level          int
	UpdatedAt      time.Time
}

// PointsTransaction represents a points transaction
type PointsTransaction struct {
	TransactionID string
	UserID        string
	Points        int // positive for earn, negative for spend
	Type          string // PURCHASE, REVIEW, REFERRAL, REDEEM, etc.
	Reference     string // order_id, review_id, etc.
	Description   string
	CreatedAt     time.Time
}

// Badge represents an achievement badge
type Badge struct {
	BadgeID     string
	Name        string
	Description string
	IconURL     string
	Criteria    BadgeCriteria
	Rarity      string // COMMON, RARE, EPIC, LEGENDARY
}

// BadgeCriteria defines how to earn a badge
type BadgeCriteria struct {
	Type      string // orders, reviews, points, streak
	Threshold int
}

// UserBadge represents a user's earned badge
type UserBadge struct {
	UserID    string
	BadgeID   string
	EarnedAt  time.Time
}

// Reward represents a redeemable reward
type Reward struct {
	RewardID      string
	Name          string
	Description   string
	PointsCost    int
	RewardType    string // DISCOUNT, PRODUCT, FREE_SHIPPING
	Value         int64
	Stock         int
	IsActive      bool
}

// NewUserPoints creates a new user points account
func NewUserPoints(userID string) *UserPoints {
	return &UserPoints{
		UserID:          userID,
		TotalPoints:     0,
		AvailablePoints: 0,
		LifetimeEarned:  0,
		LifetimeSpent:   0,
		Level:           1,
		UpdatedAt:       time.Now(),
	}
}

// Earn adds points to user account
func (u *UserPoints) Earn(points int, transType, reference, description string) *PointsTransaction {
	u.TotalPoints += points
	u.AvailablePoints += points
	u.LifetimeEarned += points
	u.UpdatedAt = time.Now()
	u.calculateLevel()

	return &PointsTransaction{
		TransactionID: uuid.New().String(),
		UserID:        u.UserID,
		Points:        points,
		Type:          transType,
		Reference:     reference,
		Description:   description,
		CreatedAt:     time.Now(),
	}
}

// Spend deducts points from user account
func (u *UserPoints) Spend(points int, transType, reference, description string) (*PointsTransaction, error) {
	if u.AvailablePoints < points {
		return nil, nil // Insufficient points
	}

	u.AvailablePoints -= points
	u.LifetimeSpent += points
	u.UpdatedAt = time.Now()

	return &PointsTransaction{
		TransactionID: uuid.New().String(),
		UserID:        u.UserID,
		Points:        -points,
		Type:          transType,
		Reference:     reference,
		Description:   description,
		CreatedAt:     time.Now(),
	}, nil
}

func (u *UserPoints) calculateLevel() {
	// Simple level calculation: 1 level per 1000 points
	u.Level = (u.LifetimeEarned / 1000) + 1
}

