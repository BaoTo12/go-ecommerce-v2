package domain

import (
	"time"

	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/google/uuid"
)

type Review struct {
	ID        string
	UserID    string
	ProductID string
	Rating    int
	Comment   string
	Images    []string
	CreatedAt time.Time
}

type ReviewStats struct {
	ProductID          string
	AverageRating      float64
	TotalReviews       int
	RatingDistribution map[int]int
}

func NewReview(userID, productID string, rating int, comment string, images []string) (*Review, error) {
	if rating < 1 || rating > 5 {
		return nil, errors.New(errors.ErrInvalidInput, "rating must be between 1 and 5")
	}
	if userID == "" || productID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "user ID and product ID are required")
	}

	return &Review{
		ID:        uuid.New().String(),
		UserID:    userID,
		ProductID: productID,
		Rating:    rating,
		Comment:   comment,
		Images:    images,
		CreatedAt: time.Now(),
	}, nil
}
