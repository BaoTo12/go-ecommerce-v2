package application

import (
	"context"

	"github.com/titan-commerce/backend/order-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

// OrderService is the application service for orders
type OrderService struct {
	repo   domain.Repository
	events domain.EventStore
	logger *logger.Logger
}

// NewOrderService creates a new order service
func NewOrderService(repo domain.Repository, events domain.EventStore, logger *logger.Logger) *OrderService {
	return &OrderService{
		repo:   repo,
		events: events,
		logger: logger,
	}
}

// CreateOrder creates a new order (Command)
func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []domain.OrderItem, shippingAddress string) (*domain.Order, error) {
	// Create order aggregate
	order, err := domain.NewOrder(userID, items, shippingAddress)
	if err != nil {
		s.logger.Error(err, "failed to create order")
		return nil, err
	}

	// Save order
	if err := s.repo.Save(ctx, order); err != nil {
		s.logger.Error(err, "failed to save order")
		return nil, err
	}

	// Publish domain event
	event := domain.NewOrderCreatedEvent(order)
	if err := s.events.SaveEvent(ctx, event); err != nil {
		s.logger.Error(err, "failed to publish order created event")
		// Don't fail the operation, event will be eventually published
	}

	s.logger.Infof("Order created: %s for user: %s", order.ID, userID)
	return order, nil
}

// GetOrder retrieves an order (Query)
func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		s.logger.Error(err, "failed to get order")
		return nil, err
	}
	return order, nil
}

// ListOrders lists orders for a user (Query)
func (s *OrderService) ListOrders(ctx context.Context, userID string, limit, offset int) ([]*domain.Order, error) {
	orders, err := s.repo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Error(err, "failed to list orders")
		return nil, err
	}
	return orders, nil
}

// CancelOrder cancels an order (Command)
func (s *OrderService) CancelOrder(ctx context.Context, orderID string, reason string) (*domain.Order, error) {
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if err := order.Cancel(reason); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, order); err != nil {
		return nil, err
	}

	s.logger.Infof("Order cancelled: %s, reason: %s", orderID, reason)
	return order, nil
}
