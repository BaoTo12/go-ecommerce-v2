package application

import (
	"context"
	"time"

	"github.com/titan-commerce/backend/fraud-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type FraudService struct {
	repo   domain.FraudRepository
	logger *logger.Logger
}

func NewFraudService(repo domain.FraudRepository, logger *logger.Logger) *FraudService {
	return &FraudService{repo: repo, logger: logger}
}

// AnalyzeTransaction performs fraud detection on a transaction
func (s *FraudService) AnalyzeTransaction(
	ctx context.Context,
	orderID, userID, ipAddress string,
	amount int64,
	deviceFingerprint string,
) (*domain.FraudCheck, error) {

	check := domain.NewFraudCheck(orderID, userID, ipAddress)
	check.DeviceFingerprint = deviceFingerprint

	// Get user risk profile
	profile, _ := s.repo.GetUserProfile(ctx, userID)

	// Check 1: Velocity - too many orders in short time
	recentChecks, _ := s.repo.GetRecentChecks(ctx, userID, 10)
	if len(recentChecks) > 5 {
		check.Flags = append(check.Flags, domain.FraudFlag{
			Type:     "velocity",
			Severity: "high",
			Score:    0.7,
			Message:  "Unusual order velocity detected",
		})
	}

	// Check 2: New account with high-value order
	if profile != nil && profile.AccountAge < 7 && amount > 50000 {
		check.Flags = append(check.Flags, domain.FraudFlag{
			Type:     "new_account_high_value",
			Severity: "medium",
			Score:    0.5,
			Message:  "New account with high-value order",
		})
	}

	// Check 3: User trust score
	if profile != nil && profile.TrustScore < 0.3 {
		check.Flags = append(check.Flags, domain.FraudFlag{
			Type:     "low_trust",
			Severity: "high",
			Score:    0.8,
			Message:  "Low user trust score",
		})
	}

	// Check 4: Previous chargebacks
	if profile != nil && profile.ChargebackCount > 2 {
		check.Flags = append(check.Flags, domain.FraudFlag{
			Type:     "chargeback_history",
			Severity: "critical",
			Score:    0.9,
			Message:  "Multiple previous chargebacks",
		})
	}

	// If no flags, add a clean flag
	if len(check.Flags) == 0 {
		check.Flags = append(check.Flags, domain.FraudFlag{
			Type:     "clean",
			Severity: "low",
			Score:    0.1,
			Message:  "No fraud indicators found",
		})
	}

	check.CalculateRisk()

	// Save check
	s.repo.SaveFraudCheck(ctx, check)

	// Update user profile
	if profile != nil {
		now := time.Now()
		profile.LastFraudCheck = &now
		profile.TotalOrders++
		s.repo.SaveUserProfile(ctx, profile)
	}

	s.logger.Infof("Fraud check: order=%s, risk=%s (%.2f), recommendation=%s",
		orderID, check.RiskLevel, check.RiskScore, check.Recommendation)

	return check, nil
}

// GetFraudCheck retrieves a fraud check
func (s *FraudService) GetFraudCheck(ctx context.Context, orderID string) (*domain.FraudCheck, error) {
	return s.repo.GetFraudCheck(ctx, orderID)
}

// GetUserRiskProfile retrieves user risk profile
func (s *FraudService) GetUserRiskProfile(ctx context.Context, userID string) (*domain.UserRiskProfile, error) {
	return s.repo.GetUserProfile(ctx, userID)
}

