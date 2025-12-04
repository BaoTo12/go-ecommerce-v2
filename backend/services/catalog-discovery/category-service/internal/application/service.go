package application

import (
	"context"

	"github.com/titan-commerce/backend/category-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type CategoryRepository interface {
	Save(ctx context.Context, category *domain.Category) error
	FindByID(ctx context.Context, categoryID string) (*domain.Category, error)
	FindByParentID(ctx context.Context, parentID string, pageSize int) ([]*domain.Category, error)
	FindAll(ctx context.Context) ([]*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
}

type CategoryService struct {
	repo   CategoryRepository
	logger *logger.Logger
}

func NewCategoryService(repo CategoryRepository, logger *logger.Logger) *CategoryService {
	return &CategoryService{
		repo:   repo,
		logger: logger,
	}
}

// CreateCategory creates a new category (Command)
func (s *CategoryService) CreateCategory(ctx context.Context, name, parentID, iconURL string) (*domain.Category, error) {
	level := 0
	if parentID != "" {
		parent, err := s.repo.FindByID(ctx, parentID)
		if err != nil {
			return nil, err
		}
		level = parent.Level + 1
	}

	category, err := domain.NewCategory(name, parentID, iconURL, level)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, category); err != nil {
		s.logger.Error(err, "failed to save category")
		return nil, err
	}

	s.logger.Infof("Category created: %s (%s)", category.Name, category.ID)
	return category, nil
}

// GetCategory retrieves a category (Query)
func (s *CategoryService) GetCategory(ctx context.Context, categoryID string) (*domain.Category, error) {
	return s.repo.FindByID(ctx, categoryID)
}

// ListCategories lists categories by parent (Query)
func (s *CategoryService) ListCategories(ctx context.Context, parentID string, pageSize int) ([]*domain.Category, error) {
	return s.repo.FindByParentID(ctx, parentID, pageSize)
}

// GetCategoryTree builds hierarchical tree (Query)
func (s *CategoryService) GetCategoryTree(ctx context.Context, rootID string, maxDepth int) ([]*domain.Category, error) {
	// Simplified - returns all categories
	// In production, build actual tree structure with children
	return s.repo.FindAll(ctx)
}

// UpdateCategory updates category details (Command)
func (s *CategoryService) UpdateCategory(ctx context.Context, categoryID, name, iconURL string) (*domain.Category, error) {
	category, err := s.repo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	if name != "" {
		if err := category.UpdateName(name); err != nil {
			return nil, err
		}
	}

	if iconURL != "" {
		category.UpdateIcon(iconURL)
	}

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	s.logger.Infof("Category updated: %s", categoryID)
	return category, nil
}
