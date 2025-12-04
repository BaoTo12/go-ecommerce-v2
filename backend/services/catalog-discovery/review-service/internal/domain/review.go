package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type Review struct {
	ID               string
	ProductID        string
	UserID           string
	UserName         string
	OrderID          string
	Rating           int // 1-5 stars
	Title            string
	Content          string
	ImageURLs        []string
	HelpfulCount     int
	VerifiedPurchase bool
	CreatedAt        time.Time
}

func NewReview(productID, userID, userName, orderID string, rating int, title, content string, imageURLs []string) (*Review, error) {
	if productID == "" || userID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "product ID and user ID are required")
	}
	if rating < 1 || rating > 5 {
		return nil, errors.New(errors.ErrInvalidInput, "rating must be between 1 and 5")
	}
	if title == "" || content == "" {
		return nil, errors.New(errors.ErrInvalidInput, "title and content are required")
	}

	return &Review{
		ID:               uuid.New().String(),
		ProductID:        productID,
		UserID:           userID,
		UserName:         userName,
		OrderID:          orderID,
		Rating:           rating,
		Title:            title,
		Content:          content,
		ImageURLs:        imageURLs,
		HelpfulCount:     0,
		VerifiedPurchase: orderID != "",
		CreatedAt:        time.Now(),
	}, nil
}

func (r *Review) IncrementHelpfulCount() {
	r.HelpfulCount++
}

func (r *Review) DecrementHelpfulCount() {
	if r.HelpfulCount > 0 {
		r.HelpfulCount--
	}
}

// Simple spam detection - check for repeated characters, excessive caps
func (r *Review) DetectSpam() bool {
	// Simplified spam detection
	if len(r.Content) < 10 {
		return true // Too short
	}
	
	// Check for excessive repeated characters  
	repeatedChars := 0
	for i := 1; i < len(r.Content); i++ {
		if r.Content[i] == r.Content[i-1] {
			repeatedChars++
			if repeatedChars > 5 {
				return true
			}
		} else {
			repeatedChars = 0
		}
	}
	
	return false
}
