package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type ProductStatus string

const (
	ProductStatusDraft          ProductStatus = "DRAFT"
	ProductStatusPendingReview  ProductStatus = "PENDING_REVIEW"
	ProductStatusApproved       ProductStatus = "APPROVED"
	ProductStatusPublished      ProductStatus = "PUBLISHED"
	ProductStatusSuspended      ProductStatus = "SUSPENDED"
)

type ProductVariant struct {
	VariantID string
	Name      string  // e.g., "Red - Large"
	Price     float64
	Stock     int
	SKU       string
}

type Product struct {
	ID           string
	SellerID     string
	Name         string
	Description  string
	CategoryID   string
	Variants     []ProductVariant
	ImageURLs    []string
	Status       ProductStatus
	Rating       float64
	ReviewCount  int
	SoldCount    int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewProduct(sellerID, name, description, categoryID string, variants []ProductVariant, imageURLs []string) (*Product, error) {
	if sellerID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "seller ID is required")
	}
	if name == "" {
		return nil, errors.New(errors.ErrInvalidInput, "product name is required")
	}
	if len(variants) == 0 {
		return nil, errors.New(errors.ErrInvalidInput, "at least one variant is required")
	}

	// Validate variants
	for _, v := range variants {
		if v.Price <= 0 {
			return nil, errors.New(errors.ErrInvalidInput, "variant price must be positive")
		}
		if v.Stock < 0 {
			return nil, errors.New(errors.ErrInvalidInput, "variant stock cannot be negative")
		}
	}

	now := time.Now()
	return &Product{
		ID:          uuid.New().String(),
		SellerID:    sellerID,
		Name:        name,
		Description: description,
		CategoryID:  categoryID,
		Variants:    variants,
		ImageURLs:   imageURLs,
		Status:      ProductStatusDraft,
		Rating:      0.0,
		ReviewCount: 0,
		SoldCount:   0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (p *Product) SubmitForReview() error {
	if p.Status != ProductStatusDraft {
		return errors.New(errors.ErrInvalidInput, "only draft products can be submitted for review")
	}
	p.Status = ProductStatusPendingReview
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) Approve() error {
	if p.Status != ProductStatusPendingReview {
		return errors.New(errors.ErrInvalidInput, "only pending products can be approved")
	}
	p.Status = ProductStatusApproved
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) Publish() error {
	if p.Status != ProductStatusApproved {
		return errors.New(errors.ErrInvalidInput, "only approved products can be published")
	}
	p.Status = ProductStatusPublished
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) UpdateRating(newRating float64, newReviewCount int) {
	p.Rating = newRating
	p.ReviewCount = newReviewCount
	p.UpdatedAt = time.Now()
}

func (p *Product) IncrementSoldCount(quantity int) {
	p.SoldCount += quantity
	p.UpdatedAt = time.Now()
}
