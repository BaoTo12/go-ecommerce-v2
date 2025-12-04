package application

import (
	"context"

	"github.com/titan-commerce/backend/social-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type SocialRepository interface {
	SaveFollow(ctx context.Context, follow *domain.Follow) error
	DeleteFollow(ctx context.Context, followerID, followeeID string) error
	GetFollowers(ctx context.Context, userID string, page, pageSize int) ([]*domain.Follow, int, error)
	GetFollowing(ctx context.Context, userID string, page, pageSize int) ([]*domain.Follow, int, error)
	GetStats(ctx context.Context, userID string) (*domain.SocialStats, error)
}

type SocialService struct {
	repo   SocialRepository
	logger *logger.Logger
}

func NewSocialService(repo SocialRepository, logger *logger.Logger) *SocialService {
	return &SocialService{
		repo:   repo,
		logger: logger,
	}
}

func (s *SocialService) FollowUser(ctx context.Context, followerID, followeeID string) error {
	follow, err := domain.NewFollow(followerID, followeeID)
	if err != nil {
		return err
	}

	if err := s.repo.SaveFollow(ctx, follow); err != nil {
		s.logger.Error(err, "failed to follow user")
		return err
	}

	s.logger.Infof("User %s followed %s", followerID, followeeID)
	return nil
}

func (s *SocialService) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	if err := s.repo.DeleteFollow(ctx, followerID, followeeID); err != nil {
		s.logger.Error(err, "failed to unfollow user")
		return err
	}

	s.logger.Infof("User %s unfollowed %s", followerID, followeeID)
	return nil
}

func (s *SocialService) GetFollowers(ctx context.Context, userID string, page, pageSize int) ([]*domain.Follow, int, error) {
	return s.repo.GetFollowers(ctx, userID, page, pageSize)
}

func (s *SocialService) GetFollowing(ctx context.Context, userID string, page, pageSize int) ([]*domain.Follow, int, error) {
	return s.repo.GetFollowing(ctx, userID, page, pageSize)
}

func (s *SocialService) GetStats(ctx context.Context, userID string) (*domain.SocialStats, error) {
	return s.repo.GetStats(ctx, userID)
}
