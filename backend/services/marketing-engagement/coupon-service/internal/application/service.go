package application

import (
	"context"

	"github.com/titan-commerce/backend/coupon-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type CouponRepository interface {
	Save(ctx context.Context, coupon *domain.Coupon) error
	FindByCode(ctx context.Context, code string) (*domain.Coupon, error)
	Update(ctx context.Context, coupon *domain.Coupon) error
}

type CouponService struct {
	repo   CouponRepository
	logger *logger.Logger
}

func NewCouponService(repo CouponRepository, logger *logger.Logger) *CouponService {
	return &CouponService{
		repo:   repo,
		logger: logger,
	}
}

// ValidateCoupon validates a coupon for an order (Query)
func (s *CouponService) ValidateCoupon(ctx context.Context, code string, orderTotal float64) (*domain.Coupon, string, error) {
	coupon, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, "coupon not found", errors.New(errors.ErrNotFound, "invalid coupon code")
	}

	valid, message := coupon.IsValid(orderTotal)
	if !valid {
		return coupon, message, errors.New(errors.ErrInvalidInput, message)
	}

	return coupon, "", nil
}

// ApplyCoupon applies a coupon to an order (Command)
func (s *CouponService) ApplyCoupon(ctx context.Context, code string, userID, orderID string, orderTotal float64) (float64, float64, error) {
	coupon, validationMsg, err := s.ValidateCoupon(ctx, orderTotal, orderTotal)
	if err != nil {
		return 0, orderTotal, err
	}

	// Calculate discount
	discountAmount := coupon.CalculateDiscount(orderTotal)
	finalTotal := orderTotal - discountAmount

	// Increment usage
	coupon.IncrementUsage()
	if err := s.repo.Update(ctx, coupon); err != nil {
		return 0, orderTotal, err
	}

	s.logger.Infof("Coupon applied: code=%s, user=%s, order=%s, discount=%.2f", code, userID, orderID, discountAmount)
	return discountAmount, finalTotal, nil
}
