package application

import (
	"context"
	"time"

	"github.com/titan-commerce/backend/flash-sale-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type FlashSaleService struct {
	repo   domain.FlashSaleRepository
	logger *logger.Logger
}

func NewFlashSaleService(repo domain.FlashSaleRepository, logger *logger.Logger) *FlashSaleService {
	return &FlashSaleService{repo: repo, logger: logger}
}

// CreateFlashSale creates a new flash sale event
func (s *FlashSaleService) CreateFlashSale(
	ctx context.Context,
	name, description string,
	startTime, endTime time.Time,
	totalStock int,
) (*domain.FlashSale, error) {
	sale, err := domain.NewFlashSale(name, description, startTime, endTime, totalStock)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveFlashSale(ctx, sale); err != nil {
		s.logger.Error(err, "failed to create flash sale")
		return nil, err
	}

	s.logger.Infof("Flash sale created: %s, stock=%d, start=%v", name, totalStock, startTime)
	return sale, nil
}

// StartFlashSale manually starts a flash sale
func (s *FlashSaleService) StartFlashSale(ctx context.Context, flashSaleID string) error {
	sale, err := s.repo.GetFlashSale(ctx, flashSaleID)
	if err != nil {
		return err
	}

	if err := sale.Start(); err != nil {
		return err
	}

	if err := s.repo.UpdateFlashSale(ctx, sale); err != nil {
		s.logger.Error(err, "failed to start flash sale")
		return err
	}

	s.logger.Infof("Flash sale started: %s", flashSaleID)
	return nil
}

// RecordPurchase records a flash sale purchase
func (s *FlashSaleService) RecordPurchase(ctx context.Context, flashSaleID string, quantity int) error {
	sale, err := s.repo.GetFlashSale(ctx, flashSaleID)
	if err != nil {
		return err
	}

	if err := sale.RecordSale(quantity); err != nil {
		return err
	}

	if err := s.repo.UpdateFlashSale(ctx, sale); err != nil {
		s.logger.Error(err, "failed to record sale")
		return err
	}

	s.logger.Infof("Flash sale purchase: sale=%s, qty=%d, remaining=%d",
		flashSaleID, quantity, sale.GetRemainingStock())

	return nil
}

// GetActiveFlashSales retrieves all active flash sales
func (s *FlashSaleService) GetActiveFlashSales(ctx context.Context) ([]*domain.FlashSale, error) {
	return s.repo.GetActiveFlashSales(ctx)
}

