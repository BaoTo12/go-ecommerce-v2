package application

import (
	"context"

	"github.com/titan-commerce/backend/product-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type ProductRepository interface {
	Save(ctx context.Context, product *domain.Product) error
	FindByID(ctx context.Context, productID string) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, productID string) error
	List(ctx context.Context, page, pageSize int) ([]*domain.Product, int, error)
}

type ProductService struct {
	repo   ProductRepository
	logger *logger.Logger
}

func NewProductService(repo ProductRepository, logger *logger.Logger) *ProductService {
	return &ProductService{
		repo:   repo,
		logger: logger,
	}
}

// CreateProduct creates a new product (Command)
func (s *ProductService) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	if err := s.repo.Save(ctx, product); err != nil {
		s.logger.Error(err, "failed to save product")
		return nil, err
	}

	s.logger.Infof("Product created: %s", product.ID)
	return product, nil
}

// GetProduct retrieves a product by ID (Query)
func (s *ProductService) GetProduct(ctx context.Context, productID string) (*domain.Product, error) {
	product, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		s.logger.Error(err, "failed to get product")
		return nil, err
	}
	return product, nil
}

// UpdateProduct updates product details (Command)
func (s *ProductService) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	if err := s.repo.Update(ctx, product); err != nil {
		s.logger.Error(err, "failed to update product")
		return nil, err
	}

	s.logger.Infof("Product updated: %s", product.ID)
	return product, nil
}

// ListProducts returns products (Query)
func (s *ProductService) ListProducts(ctx context.Context, page, pageSize int) ([]*domain.Product, int, error) {
	return s.repo.List(ctx, page, pageSize)
}
