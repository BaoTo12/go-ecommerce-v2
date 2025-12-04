package application

import (
	"context"

	"github.com/titan-commerce/backend/coupon-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type CouponService struct {
	repo   domain.CouponRepository
	logger *logger.Logger
}

func NewCouponService(repo domain.CouponRepository, logger *logger.Logger) *CouponService {
	return &CouponService{repo: repo, logger: logger}
}

// CreateCoupon creates a new coupon
func (s *CouponService) CreateCoupon(
	ctx context.Context,
	code string,
	couponType domain.CouponType,
	discountValue, minOrderValue int64,
) (*domain.Coupon, error) {
	coupon, err := domain.NewCoupon(code, couponType, discountValue, minOrderValue)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveCoupon(ctx, coupon); err != nil {
		s.logger.Error(err, "failed to create coupon")
		return nil, err
	}

	s.logger.Infof("Coupon created: code=%s, type=%s, discount=%d", code, couponType, discountValue)
	return coupon, nil
}

// ValidateCoupon validates a coupon for an order
func (s *CouponService) ValidateCoupon(
	ctx context.Context,
	code, userID string,
	orderValue int64,
) (*domain.Coupon, int64, error) {
	coupon, err := s.repo.GetCouponByCode(ctx, code)
	if err != nil {
		return nil, 0, errors.New(errors.ErrNotFound, "coupon not found")
	}

	if !coupon.IsValid() {
		return nil, 0, errors.New(errors.ErrInvalidInput, "coupon is expired or invalid")
	}

	// Check per-user limit
	if coupon.PerUserLimit > 0 {
		usageCount, _ := s.repo.GetUserUsageCount(ctx, coupon.CouponID, userID)
		if usageCount >= coupon.PerUserLimit {
			return nil, 0, errors.New(errors.ErrInvalidInput, "coupon usage limit reached for this user")
		}
	}

	discount := coupon.CalculateDiscount(orderValue)
	if discount == 0 {
		return nil, 0, errors.New(errors.ErrInvalidInput, "minimum order value not met")
	}

	return coupon, discount, nil
}

// ApplyCoupon applies a coupon to an order
func (s *CouponService) ApplyCoupon(
	ctx context.Context,
	couponID, userID, orderID string,
	discount int64,
) error {
	coupon, err := s.repo.GetCoupon(ctx, couponID)
	if err != nil {
		return err
	}

	if err := coupon.Use(); err != nil {
		return err
	}

	// Save usage
	usage := domain.NewCouponUsage(couponID, userID, orderID, discount)
	if err := s.repo.SaveUsage(ctx, usage); err != nil {
		s.logger.Error(err, "failed to save coupon usage")
		return err
	}

	// Update coupon
	if err := s.repo.UpdateCoupon(ctx, coupon); err != nil {
		s.logger.Error(err, "failed to update coupon")
		return err
	}

	s.logger.Infof("Coupon applied: id=%s, user=%s, order=%s, discount=%d",
		couponID, userID, orderID, discount)

	return nil
}

