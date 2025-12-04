package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Currency    string
	CategoryID  string
	Images      []string
	Attributes  map[string]string
	Stock       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProduct(name, description, categoryID, currency string, price float64, stock int, images []string, attributes map[string]string) (*Product, error) {
	if name == "" {
		return nil, errors.New(errors.ErrInvalidInput, "product name is required")
	}
	if price <= 0 {
		return nil, errors.New(errors.ErrInvalidInput, "price must be positive")
	}

	now := time.Now()
	if currency == "" {
		currency = "USD"
	}
	return &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		Currency:    currency,
		CategoryID:  categoryID,
		Images:      images,
		Attributes:  attributes,
		Stock:       stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
