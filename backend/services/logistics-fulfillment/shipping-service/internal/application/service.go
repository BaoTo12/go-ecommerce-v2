package application

import (
	"context"

	"github.com/titan-commerce/backend/shipping-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type ShipmentRepository interface {
	Save(ctx context.Context, shipment *domain.Shipment) error
	FindByID(ctx context.Context, shipmentID string) (*domain.Shipment, error)
	FindByOrderID(ctx context.Context, orderID string) (*domain.Shipment, error)
	Update(ctx context.Context, shipment *domain.Shipment) error
}

type ShippingService struct {
	repo   ShipmentRepository
	logger *logger.Logger
}

func NewShippingService(repo ShipmentRepository, logger *logger.Logger) *ShippingService {
	return &ShippingService{
		repo:   repo,
		logger: logger,
	}
}

// CreateShipment creates a new shipment (Command)
func (s *ShippingService) CreateShipment(ctx context.Context, orderID, carrier, originAddr, destAddr string, weight float64) (*domain.Shipment, error) {
	shipment, err := domain.NewShipment(orderID, carrier, originAddr, destAddr, weight)
	if err != nil {
		return nil, err
	}

	// Calculate shipping cost based on carrier and weight
	cost := s.calculateShippingCost(carrier, weight)
	shipment.SetShippingCost(cost)

	if err := s.repo.Save(ctx, shipment); err != nil {
		s.logger.Error(err, "failed to save shipment")
		return nil, err
	}

	s.logger.Infof("Shipment created: order=%s, carrier=%s, tracking=%s", orderID, carrier, shipment.TrackingNumber)
	return shipment, nil
}

// GetShipment retrieves shipment details (Query)
func (s *ShippingService) GetShipment(ctx context.Context, shipmentID string) (*domain.Shipment, error) {
	return s.repo.FindByID(ctx, shipmentID)
}

// UpdateStatus updates shipment status (Command)
func (s *ShippingService) UpdateStatus(ctx context.Context, shipmentID string, status domain.ShipmentStatus) (*domain.Shipment, error) {
	shipment, err := s.repo.FindByID(ctx, shipmentID)
	if err != nil {
		return nil, err
	}

	shipment.UpdateStatus(status)

	if err := s.repo.Update(ctx, shipment); err != nil {
		return nil, err
	}

	s.logger.Infof("Shipment status updated: shipment=%s, status=%s", shipmentID, status)
	return shipment, nil
}

// Simple cost calculation - in production, integrate with carrier APIs
func (s *ShippingService) calculateShippingCost(carrier string, weight float64) float64 {
	baseCost := 5.0
	perKgCost := 2.0

	switch carrier {
	case "DHL":
		return baseCost + (weight * perKgCost * 1.2)
	case "FedEx":
		return baseCost + (weight * perKgCost * 1.1)
	case "UPS":
		return baseCost + (weight * perKgCost)
	default:
		return baseCost + (weight * perKgCost)
	}
}
