package application

import (
	"context"

	"github.com/titan-commerce/backend/search-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type SearchRepository interface {
	IndexProduct(ctx context.Context, product *domain.ProductDocument) error
	Search(ctx context.Context, query string, page, pageSize int) ([]*domain.ProductDocument, int, error)
}

type SearchService struct {
	repo   SearchRepository
	logger *logger.Logger
}

func NewSearchService(repo SearchRepository, logger *logger.Logger) *SearchService {
	return &SearchService{
		repo:   repo,
		logger: logger,
	}
}

func (s *SearchService) IndexProduct(ctx context.Context, product *domain.ProductDocument) error {
	if err := s.repo.IndexProduct(ctx, product); err != nil {
		s.logger.Error(err, "failed to index product")
		return err
	}
	s.logger.Infof("Indexed product %s", product.ID)
	return nil
}

func (s *SearchService) SearchProducts(ctx context.Context, query string, page, pageSize int) ([]*domain.ProductDocument, int, error) {
	return s.repo.Search(ctx, query, page, pageSize)
}
