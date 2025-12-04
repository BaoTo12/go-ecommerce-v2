package application

import (
	"context"

	"github.com/titan-commerce/backend/feed-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type FeedRepository interface {
	SavePost(ctx context.Context, post *domain.Post) error
	DeletePost(ctx context.Context, postID string) error
	GetGlobalFeed(ctx context.Context, page, pageSize int) ([]*domain.Post, error)
	GetUserFeed(ctx context.Context, userID string, page, pageSize int) ([]*domain.Post, error)
}

type FeedService struct {
	repo   FeedRepository
	logger *logger.Logger
}

func NewFeedService(repo FeedRepository, logger *logger.Logger) *FeedService {
	return &FeedService{
		repo:   repo,
		logger: logger,
	}
}

func (s *FeedService) PublishPost(ctx context.Context, userID, content, mediaURL string, tags []string) (*domain.Post, error) {
	post, err := domain.NewPost(userID, content, mediaURL, tags)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SavePost(ctx, post); err != nil {
		s.logger.Error(err, "failed to save post")
		return nil, err
	}

	s.logger.Infof("User %s published post %s", userID, post.ID)
	return post, nil
}

func (s *FeedService) DeletePost(ctx context.Context, postID, userID string) error {
	// In a real app, we'd verify ownership here or in the repo
	if err := s.repo.DeletePost(ctx, postID); err != nil {
		s.logger.Error(err, "failed to delete post")
		return err
	}
	return nil
}

func (s *FeedService) GetFeed(ctx context.Context, userID string, page, pageSize int) ([]*domain.Post, error) {
	// Simple algorithm:
	// 1. If userID is present, try to get personalized feed (e.g. from followed users)
	// 2. Fallback to global feed (recent posts)
	
	// For MVP, we'll just return the global feed
	return s.repo.GetGlobalFeed(ctx, page, pageSize)
}
