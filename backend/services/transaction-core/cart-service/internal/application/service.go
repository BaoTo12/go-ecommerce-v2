package application

import (
	"context"

	"github.com/titan-commerce/backend/cart-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type CartRepository interface {
	Save(ctx context.Context, cart *domain.Cart) error
	FindByUserID(ctx context.Context, userID string) (*domain.Cart, error)
	Delete(ctx context.Context, userID string) error
}

type CartService struct {
	repo   CartRepository
	logger *logger.Logger
}

func NewCartService(repo CartRepository, logger *logger.Logger) *CartService {
	return &CartService{
		repo:   repo,
		logger: logger,
	}
}

// AddToCart adds an item to user's cart (Command)
func (s *CartService) AddToCart(ctx context.Context, userID, productID, productName string, quantity int, unitPrice float64) (*domain.Cart, error) {
	cart, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		s.logger.Error(err, "failed to get cart")
		return nil, err
	}

	cart.AddItem(productID, productName, quantity, unitPrice)

	if err := s.repo.Save(ctx, cart); err != nil {
		s.logger.Error(err, "failed to save cart")
		return nil, err
	}

	s.logger.Infof("Added to cart: user=%s, product=%s, qty=%d", userID, productID, quantity)
	return cart, nil
}

// RemoveFromCart removes an item from cart (Command)
func (s *CartService) RemoveFromCart(ctx context.Context, userID, productID string) (*domain.Cart, error) {
	cart, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	cart.RemoveItem(productID)

	if err := s.repo.Save(ctx, cart); err != nil {
		return nil, err
	}

	s.logger.Infof("Removed from cart: user=%s, product=%s", userID, productID)
	return cart, nil
}

// UpdateQuantity updates item quantity in cart (Command)
func (s *CartService) UpdateQuantity(ctx context.Context, userID, productID string, quantity int) (*domain.Cart, error) {
	cart, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	cart.UpdateQuantity(productID, quantity)

	if err := s.repo.Save(ctx, cart); err != nil {
		return nil, err
	}

	s.logger.Infof("Updated cart quantity: user=%s, product=%s, qty=%d", userID, productID, quantity)
	return cart, nil
}

// GetCart retrieves user's cart (Query)
func (s *CartService) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	cart, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		s.logger.Error(err, "failed to get cart")
		return nil, err
	}
	return cart, nil
}

// ClearCart empties user's cart (Command)
func (s *CartService) ClearCart(ctx context.Context, userID string) error {
	if err := s.repo.Delete(ctx, userID); err != nil {
		s.logger.Error(err, "failed to clear cart")
		return err
	}

	s.logger.Infof("Cart cleared: user=%s", userID)
	return nil
}
