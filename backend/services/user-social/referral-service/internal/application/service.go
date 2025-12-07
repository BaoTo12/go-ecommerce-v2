package application

import (
	"context"

	"github.com/titan-commerce/backend/referral-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

// ReferralService handles referral operations
type ReferralService struct {
	repo   domain.Repository
	logger *logger.Logger
}

// NewReferralService creates a new referral service
func NewReferralService(repo domain.Repository, log *logger.Logger) *ReferralService {
	return &ReferralService{
		repo:   repo,
		logger: log,
	}
}

// GenerateCode generates a unique referral code for a user
func (s *ReferralService) GenerateCode(ctx context.Context, userID, userName string) (*domain.ReferralCode, error) {
	// Check if user already has a code
	existing, _ := s.repo.GetReferralCode(ctx, userID)
	if existing != nil {
		return existing, nil
	}

	// Create new code
	code := domain.NewReferralCode(userID, userName)
	if err := s.repo.SaveReferralCode(ctx, code); err != nil {
		return nil, err
	}

	s.logger.Infof("Generated referral code for user %s: %s", userID, code.Code)
	return code, nil
}

// GetReferralStats returns referral statistics for a user
func (s *ReferralService) GetReferralStats(ctx context.Context, userID string) (*domain.ReferralStats, error) {
	return s.repo.GetReferralStats(ctx, userID)
}

// GetReferrals returns all referrals for a user
func (s *ReferralService) GetReferrals(ctx context.Context, userID string) ([]*domain.Referral, error) {
	return s.repo.GetReferrals(ctx, userID)
}

// RedeemCode processes a referral code redemption
func (s *ReferralService) RedeemCode(ctx context.Context, code, newUserID, newUserName string) (*domain.Referral, error) {
	// Get referral code
	refCode, err := s.repo.GetReferralByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if refCode == nil {
		return nil, nil // Code not found
	}

	// Check if code is still valid
	if refCode.Uses >= refCode.MaxUses {
		return nil, nil // Max uses reached
	}

	// Create referral record
	referral := domain.NewReferral(refCode.UserID, newUserID, newUserName, code)
	if err := s.repo.SaveReferral(ctx, referral); err != nil {
		return nil, err
	}

	// Update code usage
	refCode.Uses++
	if err := s.repo.SaveReferralCode(ctx, refCode); err != nil {
		return nil, err
	}

	s.logger.Infof("Referral code %s redeemed by %s", code, newUserID)
	return referral, nil
}

// CompleteReferral marks a referral as completed (after first order)
func (s *ReferralService) CompleteReferral(ctx context.Context, referredUserID string) error {
	referrals, err := s.repo.GetReferrals(ctx, referredUserID)
	if err != nil {
		return err
	}

	for _, ref := range referrals {
		if ref.ReferredID == referredUserID && ref.Status == domain.ReferralPending {
			ref.Complete(500) // 500 coins reward
			if err := s.repo.SaveReferral(ctx, ref); err != nil {
				return err
			}
			s.logger.Infof("Referral completed: referrer=%s referred=%s reward=%d", ref.ReferrerID, ref.ReferredID, ref.Reward)
		}
	}

	return nil
}
