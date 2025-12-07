package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReferralStatus represents the status of a referral
type ReferralStatus string

const (
	ReferralPending   ReferralStatus = "pending"
	ReferralCompleted ReferralStatus = "completed"
	ReferralExpired   ReferralStatus = "expired"
)

// Referral represents a referral relationship
type Referral struct {
	ID           string         `json:"id"`
	ReferrerID   string         `json:"referrer_id"`
	ReferredID   string         `json:"referred_id"`
	ReferredName string         `json:"referred_name"`
	ReferredAvatar string       `json:"referred_avatar"`
	Code         string         `json:"code"`
	Status       ReferralStatus `json:"status"`
	Reward       int            `json:"reward"`
	JoinedAt     time.Time      `json:"joined_at"`
	CompletedAt  *time.Time     `json:"completed_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

// ReferralStats represents user's referral statistics
type ReferralStats struct {
	UserID          string `json:"user_id"`
	TotalReferrals  int    `json:"total_referrals"`
	CompletedCount  int    `json:"completed_count"`
	PendingCount    int    `json:"pending_count"`
	TotalEarned     int    `json:"total_earned"`
	PendingEarnings int    `json:"pending_earnings"`
}

// ReferralCode represents a unique referral code
type ReferralCode struct {
	Code      string    `json:"code"`
	UserID    string    `json:"user_id"`
	Link      string    `json:"link"`
	Uses      int       `json:"uses"`
	MaxUses   int       `json:"max_uses"`
	RewardPerUse int    `json:"reward_per_use"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewReferralCode creates a new referral code for a user
func NewReferralCode(userID, userName string) *ReferralCode {
	code := "SHOPEE-" + userName + "-" + uuid.New().String()[:4]
	return &ReferralCode{
		Code:         code,
		UserID:       userID,
		Link:         "https://shopee.vn/register?ref=" + code,
		Uses:         0,
		MaxUses:      50, // Max 50 referrals per month
		RewardPerUse: 500, // 500 coins per referral
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().AddDate(0, 1, 0), // 1 month validity
	}
}

// NewReferral creates a new referral record
func NewReferral(referrerID, referredID, referredName, code string) *Referral {
	return &Referral{
		ID:           uuid.New().String(),
		ReferrerID:   referrerID,
		ReferredID:   referredID,
		ReferredName: referredName,
		ReferredAvatar: "https://ui-avatars.com/api/?name=" + referredName + "&background=random",
		Code:         code,
		Status:       ReferralPending,
		Reward:       0,
		JoinedAt:     time.Now(),
		CreatedAt:    time.Now(),
	}
}

// Complete marks the referral as completed and assigns reward
func (r *Referral) Complete(reward int) {
	now := time.Now()
	r.Status = ReferralCompleted
	r.Reward = reward
	r.CompletedAt = &now
}

// Repository interface for referral data
type Repository interface {
	GetReferralCode(ctx interface{}, userID string) (*ReferralCode, error)
	SaveReferralCode(ctx interface{}, code *ReferralCode) error
	GetReferrals(ctx interface{}, userID string) ([]*Referral, error)
	SaveReferral(ctx interface{}, referral *Referral) error
	GetReferralStats(ctx interface{}, userID string) (*ReferralStats, error)
	GetReferralByCode(ctx interface{}, code string) (*ReferralCode, error)
}
