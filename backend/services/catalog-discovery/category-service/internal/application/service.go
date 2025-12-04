package application

import (
	"context"

	"github.com/titan-commerce/backend/category-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type CategoryRepository interface {
	Save(ctx context.Context, category *domain.Category) error
	FindByID(ctx context.Context, id string) (*domain.Category, error)
	List(ctx context.Context, page, pageSize int) ([]*domain.Category, int, error)
	GetAll(ctx context.Context) ([]*domain.Category, error)
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

func (s *CategoryService) CreateCategory(ctx context.Context, name, description, parentID, imageURL string) (*domain.Category, error) {
	category, err := domain.NewCategory(name, description, parentID, imageURL)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, category); err != nil {
		s.logger.Error(err, "failed to save category")
		return nil, err
	}

	s.logger.Infof("Category created: %s", category.Name)
	return category, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, id string) (*domain.Category, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CategoryService) ListCategories(ctx context.Context, page, pageSize int) ([]*domain.Category, int, error) {
	return s.repo.List(ctx, page, pageSize)
}

func (s *CategoryService) GetCategoryTree(ctx context.Context) ([]*domain.CategoryNode, error) {
	allCategories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Build tree
	categoryMap := make(map[string]*domain.CategoryNode)
	var roots []*domain.CategoryNode

	// First pass: create nodes
	for _, c := range allCategories {
		categoryMap[c.ID] = &domain.CategoryNode{
			Category: c,
			Children: []*domain.CategoryNode{},
		}
	}

	// Second pass: link children
	for _, c := range allCategories {
		node := categoryMap[c.ID]
		if c.ParentID == "" {
			roots = append(roots, node)
		} else {
			if parent, ok := categoryMap[c.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			} else {
				// Parent not found, treat as root (or handle error)
				roots = append(roots, node)
			}
		}
	}

	return roots, nil
}
