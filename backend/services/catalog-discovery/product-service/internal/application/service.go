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
	ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Product, int, error)
	Search(ctx context.Context, query, categoryID string, page, pageSize int) ([]*domain.Product, int, error)
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
func (s *ProductService) CreateProduct(ctx context.Context, sellerID, name, description, categoryID string, variants []domain.ProductVariant, imageURLs []string) (*domain.Product, error) {
	product, err := domain.NewProduct(sellerID, name, description, categoryID, variants, imageURLs)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, product); err != nil {
		s.logger.Error(err, "failed to save product")
		return nil, err
	}

	s.logger.Infof("Product created: %s by seller: %s", product.ID, sellerID)
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
func (s *ProductService) UpdateProduct(ctx context.Context, productID, name, description string, variants []domain.ProductVariant) (*domain.Product, error) {
	product, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	product.Name = name
	product.Description = description
	product.Variants = variants

	if err := s.repo.Update(ctx, product); err != nil {
		s.logger.Error(err, "failed to update product")
		return nil, err
	}

	s.logger.Infof("Product updated: %s", productID)
	return product, nil
}

// PublishProduct publishes a product (Command)
func (s *ProductService) PublishProduct(ctx context.Context, productID string) error {
	product, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	// Workflow: Submit → Approve → Publish
	if product.Status == domain.ProductStatusDraft {
		if err := product.SubmitForReview(); err != nil {
			return err
		}
	}
	if product.Status == domain.ProductStatusPendingReview {
		if err := product.Approve(); err != nil {
			return err
		}
	}
	if product.Status == domain.ProductStatusApproved {
		if err := product.Publish(); err != nil {
			return err
		}
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return err
	}

	s.logger.Infof("Product published: %s", productID)
	return nil
}

// ListProducts returns products for a seller (Query)
func (s *ProductService) ListProducts(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Product, int, error) {
	return s.repo.ListBySeller(ctx, sellerID, page, pageSize)
}

// SearchProducts searches products (Query)
func (s *ProductService) SearchProducts(ctx context.Context, query, categoryID string, page, pageSize int) ([]*domain.Product, int, error) {
	return s.repo.Search(ctx, query, categoryID, page, pageSize)
}
