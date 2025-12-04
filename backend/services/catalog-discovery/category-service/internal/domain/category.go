package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type Category struct {
	ID           string
	Name         string
	Slug         string
	ParentID     string
	Level        int
	IconURL      string
	ProductCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewCategory(name, parentID, iconURL string, level int) (*Category, error) {
	if name == "" {
		return nil, errors.New(errors.ErrInvalidInput, "category name is required")
	}

	slug := generateSlug(name)
	now := time.Now()

	return &Category{
		ID:           uuid.New().String(),
		Name:         name,
		Slug:         slug,
		ParentID:     parentID,
		Level:        level,
		IconURL:      iconURL,
		ProductCount: 0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (c *Category) UpdateName(name string) error {
	if name == "" {
		return errors.New(errors.ErrInvalidInput, "category name cannot be empty")
	}
	c.Name = name
	c.Slug = generateSlug(name)
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Category) UpdateIcon(iconURL string) {
	c.IconURL = iconURL
	c.UpdatedAt = time.Now()
}

func (c *Category) IncrementProductCount() {
	c.ProductCount++
	c.UpdatedAt = time.Now()
}

func (c *Category) DecrementProductCount() {
	if c.ProductCount > 0 {
		c.ProductCount--
		c.UpdatedAt = time.Now()
	}
}

// Simple slug generation (replace spaces with hyphens, lowercase)
func generateSlug(name string) string {
	// Simplified - in production use proper slug library
	slug := name
	// This is a placeholder - use proper slug generation in production
	return slug
}
