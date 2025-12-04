package application

import (
	"context"

	"github.com/titan-commerce/backend/gamification-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type GamificationService struct {
	repo   domain.GamificationRepository
	logger *logger.Logger
}

func NewGamificationService(repo domain.GamificationRepository, logger *logger.Logger) *GamificationService {
	return &GamificationService{repo: repo, logger: logger}
}

// EarnPoints awards points to a user
func (s *GamificationService) EarnPoints(
	ctx context.Context,
	userID string,
	points int,
	transType, reference, description string,
) error {
	userPoints, err := s.repo.GetUserPoints(ctx, userID)
	if err != nil {
		// Create new account
		userPoints = domain.NewUserPoints(userID)
	}

	transaction := userPoints.Earn(points, transType, reference, description)

	if err := s.repo.SaveUserPoints(ctx, userPoints); err != nil {
		s.logger.Error(err, "failed to save user points")
		return err
	}

	if err := s.repo.SaveTransaction(ctx, transaction); err != nil {
		s.logger.Error(err, "failed to save transaction")
	}

	s.logger.Infof("Points earned: user=%s, points=%d, type=%s, level=%d",
		userID, points, transType, userPoints.Level)

	return nil
}

// RedeemPoints redeems points for a reward
func (s *GamificationService) RedeemPoints(
	ctx context.Context,
	userID, rewardID string,
) error {
	userPoints, err := s.repo.GetUserPoints(ctx, userID)
	if err != nil {
		return err
	}

	reward, err := s.repo.GetReward(ctx, rewardID)
	if err != nil {
		return err
	}

	transaction, err := userPoints.Spend(reward.PointsCost, "REDEEM", rewardID, reward.Name)
	if err != nil || transaction == nil {
		s.logger.Warnf("Insufficient points: user=%s, available=%d, required=%d",
			userID, userPoints.AvailablePoints, reward.PointsCost)
		return nil
	}

	if err := s.repo.SaveUserPoints(ctx, userPoints); err != nil {
		s.logger.Error(err, "failed to save user points")
		return err
	}

	if err := s.repo.SaveTransaction(ctx, transaction); err != nil {
		s.logger.Error(err, "failed to save transaction")
	}

	s.logger.Infof("Points redeemed: user=%s, reward=%s, points=%d",
		userID, reward.Name, reward.PointsCost)

	return nil
}

// GetUserPoints retrieves user points balance
func (s *GamificationService) GetUserPoints(ctx context.Context, userID string) (*domain.UserPoints, error) {
	return s.repo.GetUserPoints(ctx, userID)
}

