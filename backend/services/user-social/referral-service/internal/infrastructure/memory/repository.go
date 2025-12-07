package memory

import (
	"sync"
	"time"

	"github.com/titan-commerce/backend/referral-service/internal/domain"
)

// ReferralRepository is an in-memory implementation
type ReferralRepository struct {
	codes     map[string]*domain.ReferralCode
	referrals []*domain.Referral
	mu        sync.RWMutex
}

// NewReferralRepository creates a new in-memory repository with sample data
func NewReferralRepository() *ReferralRepository {
	repo := &ReferralRepository{
		codes:     make(map[string]*domain.ReferralCode),
		referrals: make([]*domain.Referral, 0),
	}
	repo.seedData()
	return repo
}

func (r *ReferralRepository) seedData() {
	// Sample referral code
	code := &domain.ReferralCode{
		Code:         "SHOPEE-NGUYENVANA-2024",
		UserID:       "u1",
		Link:         "https://shopee.vn/register?ref=SHOPEE-NGUYENVANA-2024",
		Uses:         3,
		MaxUses:      50,
		RewardPerUse: 500,
		CreatedAt:    time.Now().Add(-7 * 24 * time.Hour),
		ExpiresAt:    time.Now().Add(23 * 24 * time.Hour),
	}
	r.codes["u1"] = code
	r.codes["SHOPEE-NGUYENVANA-2024"] = code

	// Sample referrals
	completedAt := time.Now().Add(-24 * time.Hour)
	r.referrals = []*domain.Referral{
		{ID: "r1", ReferrerID: "u1", ReferredID: "u2", ReferredName: "Nguyễn Văn B", ReferredAvatar: "https://ui-avatars.com/api/?name=NB&background=random", Code: code.Code, Status: domain.ReferralCompleted, Reward: 500, JoinedAt: time.Now().Add(-5 * 24 * time.Hour), CompletedAt: &completedAt},
		{ID: "r2", ReferrerID: "u1", ReferredID: "u3", ReferredName: "Trần Thị C", ReferredAvatar: "https://ui-avatars.com/api/?name=TC&background=random", Code: code.Code, Status: domain.ReferralCompleted, Reward: 500, JoinedAt: time.Now().Add(-4 * 24 * time.Hour), CompletedAt: &completedAt},
		{ID: "r3", ReferrerID: "u1", ReferredID: "u4", ReferredName: "Lê Văn D", ReferredAvatar: "https://ui-avatars.com/api/?name=LD&background=random", Code: code.Code, Status: domain.ReferralPending, Reward: 0, JoinedAt: time.Now().Add(-1 * 24 * time.Hour)},
	}
}

func (r *ReferralRepository) GetReferralCode(ctx interface{}, userID string) (*domain.ReferralCode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if code, ok := r.codes[userID]; ok {
		return code, nil
	}
	return nil, nil
}

func (r *ReferralRepository) SaveReferralCode(ctx interface{}, code *domain.ReferralCode) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.codes[code.UserID] = code
	r.codes[code.Code] = code
	return nil
}

func (r *ReferralRepository) GetReferrals(ctx interface{}, userID string) ([]*domain.Referral, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Referral, 0)
	for _, ref := range r.referrals {
		if ref.ReferrerID == userID {
			result = append(result, ref)
		}
	}
	return result, nil
}

func (r *ReferralRepository) SaveReferral(ctx interface{}, referral *domain.Referral) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Update if exists
	for i, ref := range r.referrals {
		if ref.ID == referral.ID {
			r.referrals[i] = referral
			return nil
		}
	}
	// Add new
	r.referrals = append(r.referrals, referral)
	return nil
}

func (r *ReferralRepository) GetReferralStats(ctx interface{}, userID string) (*domain.ReferralStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := &domain.ReferralStats{UserID: userID}
	for _, ref := range r.referrals {
		if ref.ReferrerID == userID {
			stats.TotalReferrals++
			if ref.Status == domain.ReferralCompleted {
				stats.CompletedCount++
				stats.TotalEarned += ref.Reward
			} else if ref.Status == domain.ReferralPending {
				stats.PendingCount++
				stats.PendingEarnings += 500 // Potential reward
			}
		}
	}
	return stats, nil
}

func (r *ReferralRepository) GetReferralByCode(ctx interface{}, code string) (*domain.ReferralCode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if refCode, ok := r.codes[code]; ok {
		return refCode, nil
	}
	return nil, nil
}
