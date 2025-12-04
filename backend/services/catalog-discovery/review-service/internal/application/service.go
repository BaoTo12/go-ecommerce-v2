package application

import (
	"context"

	"github.com/titan-commerce/backend/review-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type ReviewRepository interface {
	Save(ctx context.Context, review *domain.Review) error
	GetByProduct(ctx context.Context, productID string, page, pageSize int) ([]*domain.Review, int, error)
	GetStats(ctx context.Context, productID string) (*domain.ReviewStats, error)
}

type ReviewService struct {
	repo   ReviewRepository
	logger *logger.Logger
}

func NewReviewService(repo ReviewRepository, logger *logger.Logger) *ReviewService {
	return &ReviewService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, userID, productID string, rating int, comment string, images []string) (*domain.Review, error) {
	review, err := domain.NewReview(userID, productID, rating, comment, images)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, review); err != nil {
		s.logger.Error(err, "failed to save review")
		return nil, err
	}

	s.logger.Infof("Review created for product %s by user %s", productID, userID)
	return review, nil
}

func (s *ReviewService) GetProductReviews(ctx context.Context, productID string, page, pageSize int) ([]*domain.Review, int, error) {
	return s.repo.GetByProduct(ctx, productID, page, pageSize)
}

func (s *ReviewService) GetStats(ctx context.Context, productID string) (*domain.ReviewStats, error) {
	return s.repo.GetStats(ctx, productID)
}
