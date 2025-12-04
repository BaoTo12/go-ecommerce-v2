package application

import (
	"context"

	"github.com/titan-commerce/backend/review-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type ReviewRepository interface {
	Save(ctx context.Context, review *domain.Review) error
	FindByID(ctx context.Context, reviewID string) (*domain.Review, error)
	FindByProductID(ctx context.Context, productID string, page, pageSize int, minRating int) ([]*domain.Review, int, error)
	Update(ctx context.Context, review *domain.Review) error
	CalculateAverageRating(ctx context.Context, productID string) (float64, error)
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

// CreateReview creates a product review (Command)
func (s *ReviewService) CreateReview(ctx context.Context, productID, userID, userName, orderID string, rating int, title, content string, imageURLs []string) (*domain.Review, error) {
	review, err := domain.NewReview(productID, userID, userName, orderID, rating, title, content, imageURLs)
	if err != nil {
		return nil, err
	}

	// Spam detection
	if review.DetectSpam() {
		return nil, errors.New(errors.ErrInvalidInput, "review appears to be spam")
	}

	if err := s.repo.Save(ctx, review); err != nil {
		s.logger.Error(err, "failed to save review")
		return nil, err
	}

	s.logger.Infof("Review created: product=%s, user=%s, rating=%d", productID, userID, rating)
	return review, nil
}

// GetReviews retrieves reviews for a product (Query)
func (s *ReviewService) GetReviews(ctx context.Context, productID string, page, pageSize, minRating int) ([]*domain.Review, int, float64, error) {
	reviews, total, err := s.repo.FindByProductID(ctx, productID, page, pageSize, minRating)
	if err != nil {
		return nil, 0, 0, err
	}

	// Calculate average rating
	avgRating, err := s.repo.CalculateAverageRating(ctx, productID)
	if err != nil {
		avgRating = 0
	}

	return reviews, total, avgRating, nil
}

// VoteReview marks review as helpful (Command)
func (s *ReviewService) VoteReview(ctx context.Context, reviewID, userID string, helpful bool) error {
	review, err := s.repo.FindByID(ctx, reviewID)
	if err != nil {
		return err
	}

	if helpful {
		review.IncrementHelpfulCount()
	} else {
		review.DecrementHelpfulCount()
	}

	if err := s.repo.Update(ctx, review); err != nil {
		return err
	}

	s.logger.Infof("Review voted: review=%s, helpful=%v", reviewID, helpful)
	return nil
}
