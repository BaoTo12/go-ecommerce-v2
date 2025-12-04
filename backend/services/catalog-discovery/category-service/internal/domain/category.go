package domain

import (
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/google/uuid"
)

type Category struct {
	ID          string
	Name        string
	Description string
	ParentID    string
	ImageURL    string
}

type CategoryNode struct {
	Category *Category
	Children []*CategoryNode
}

func NewCategory(name, description, parentID, imageURL string) (*Category, error) {
	if name == "" {
		return nil, errors.New(errors.ErrInvalidInput, "category name is required")
	}

	return &Category{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		ParentID:    parentID,
		ImageURL:    imageURL,
	}, nil
}
