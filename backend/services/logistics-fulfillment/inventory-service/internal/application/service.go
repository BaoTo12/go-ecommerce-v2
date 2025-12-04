package application

import (
	"context"

	"github.com/titan-commerce/backend/inventory-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type InventoryService struct {
	repo   domain.StockRepository
	logger *logger.Logger
}

func NewInventoryService(repo domain.StockRepository, logger *logger.Logger) *InventoryService {
	return &InventoryService{
		repo:   repo,
		logger: logger,
	}
}

// ReserveStock reserves stock atomically (Command)
// This uses Redis Lua script for atomic reserve operation
func (s *InventoryService) ReserveStock(ctx context.Context, productID string, quantity int, ttlMinutes int) (*domain.Reservation, error) {
	reservation, err := domain.NewReservation(productID, quantity, ttlMinutes)
	if err != nil {
		return nil, err
	}

	success, err := s.repo.ReserveStock(ctx, productID, quantity, reservation.ReservationID, ttlMinutes)
	if err != nil {
		s.logger.Error(err, "failed to reserve stock")
		return nil, err
	}

	if !success {
		s.logger.Warnf("Insufficient stock: product=%s, requested=%d", productID, quantity)
		return nil, nil
	}

	s.logger.Infof("Stock reserved: product=%s, quantity=%d, reservation=%s",
		productID, quantity, reservation.ReservationID)

	return reservation, nil
}

// CommitReservation commits a reservation (Command)
// Converts reserved stock to actual deduction
func (s *InventoryService) CommitReservation(ctx context.Context, reservationID string) error {
	reservation, err := s.repo.GetReservation(ctx, reservationID)
	if err != nil {
		return err
	}

	if err := reservation.Commit(); err != nil {
		return err
	}

	if err := s.repo.CommitReservation(ctx, reservationID); err != nil {
		s.logger.Error(err, "failed to commit reservation")
		return err
	}

	s.logger.Infof("Reservation committed: reservation=%s, product=%s",
		reservationID, reservation.ProductID)

	return nil
}

// RollbackReservation rolls back a reservation (Command)
// Returns reserved stock to available pool
func (s *InventoryService) RollbackReservation(ctx context.Context, reservationID string) error {
	reservation, err := s.repo.GetReservation(ctx, reservationID)
	if err != nil {
		return err
	}

	if err := reservation.Rollback(); err != nil {
		return err
	}

	if err := s.repo.RollbackReservation(ctx, reservationID); err != nil {
		s.logger.Error(err, "failed to rollback reservation")
		return err
	}

	s.logger.Infof("Reservation rolled back: reservation=%s, product=%s",
		reservationID, reservation.ProductID)

	return nil
}

// CheckStockAvailability checks if stock is available (Query)
func (s *InventoryService) CheckStockAvailability(ctx context.Context, productID string, quantity int) (bool, error) {
	return s.repo.CheckAvailability(ctx, productID, quantity)
}

// GetAvailableStock retrieves available stock (Query)
func (s *InventoryService) GetAvailableStock(ctx context.Context, productID string) (int, error) {
	return s.repo.GetAvailableStock(ctx, productID)
}

// GetStockInfo retrieves complete stock information (Query)
func (s *InventoryService) GetStockInfo(ctx context.Context, productID string) (*domain.Stock, error) {
	available, err := s.repo.GetAvailableStock(ctx, productID)
	if err != nil {
		return nil, err
	}

	reserved, err := s.repo.GetReservedStock(ctx, productID)
	if err != nil {
		return nil, err
	}

	return &domain.Stock{
		ProductID:         productID,
		AvailableQuantity: available,
		ReservedQuantity:  reserved,
		TotalQuantity:     available + reserved,
	}, nil
}

// AddStock adds stock (Command)
func (s *InventoryService) AddStock(ctx context.Context, productID string, quantity int) error {
	if err := s.repo.AddStock(ctx, productID, quantity); err != nil {
		s.logger.Error(err, "failed to add stock")
		return err
	}

	s.logger.Infof("Stock added: product=%s, quantity=%d", productID, quantity)
	return nil
}

// CleanExpiredReservations cleans up expired reservations (Command)
func (s *InventoryService) CleanExpiredReservations(ctx context.Context) error {
	if err := s.repo.CleanExpiredReservations(ctx); err != nil {
		s.logger.Error(err, "failed to clean expired reservations")
		return err
	}

	s.logger.Info("Expired reservations cleaned")
	return nil
}

// CheckAndAlertLowStock checks for low stock and creates alerts (Command)
func (s *InventoryService) CheckAndAlertLowStock(ctx context.Context, productID string, threshold int) error {
	available, err := s.repo.GetAvailableStock(ctx, productID)
	if err != nil {
		return err
	}

	if available <= threshold {
		alert := &domain.StockAlert{
			ProductID:     productID,
			CurrentStock:  available,
			ThresholdType: "LOW_STOCK",
		}

		if available == 0 {
			alert.ThresholdType = "OUT_OF_STOCK"
		}

		if err := s.repo.SaveAlert(ctx, alert); err != nil {
			return err
		}

		s.logger.Warnf("Stock alert: product=%s, stock=%d, type=%s",
			productID, available, alert.ThresholdType)
	}

	return nil
}

