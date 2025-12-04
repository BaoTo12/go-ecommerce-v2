package application

import (
	"context"
	"strings"
	"time"

	"github.com/titan-commerce/backend/coupon-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	"github.com/google/uuid"
)

type CouponRepository interface {
	Save(ctx context.Context, coupon *domain.Coupon) error
	FindByID(ctx context.Context, couponID string) (*domain.Coupon, error)
	FindByCode(ctx context.Context, code string) (*domain.Coupon, error)
	Update(ctx context.Context, coupon *domain.Coupon) error
	FindActive(ctx context.Context) ([]*domain.Coupon, error)
	SaveUsage(ctx context.Context, usage *domain.CouponUsage) error
	GetUserUsage(ctx context.Context, couponID, userID string) (int, error)
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

// CreateCoupon creates a new coupon
func (s *CouponService) CreateCoupon(ctx context.Context, code, name, description string, couponType domain.CouponType, value, minOrder, maxDiscount float64, totalQty, maxPerUser int, validFrom, validUntil time.Time) (*domain.Coupon, error) {
	// Normalize code
	code = strings.ToUpper(strings.TrimSpace(code))

	// Check if code exists
	if existing, _ := s.repo.FindByCode(ctx, code); existing != nil {
		return nil, errors.New(errors.ErrConflict, "coupon code already exists")
	}

	coupon := domain.NewCoupon(code, name, description, couponType, value, minOrder, maxDiscount, totalQty, maxPerUser, validFrom, validUntil)
	
	if err := s.repo.Save(ctx, coupon); err != nil {
		return nil, err
	}

	s.logger.Infof("Coupon created: %s (%s)", code, couponType)
	return coupon, nil
}

// GetCoupon retrieves a coupon by ID
func (s *CouponService) GetCoupon(ctx context.Context, couponID string) (*domain.Coupon, error) {
	return s.repo.FindByID(ctx, couponID)
}

// GetCouponByCode retrieves a coupon by code
func (s *CouponService) GetCouponByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	return s.repo.FindByCode(ctx, code)
}

// ValidateCoupon validates a coupon for a user and order
func (s *CouponService) ValidateCoupon(ctx context.Context, code, userID string, orderValue float64, categoryIDs, productIDs []string) (*domain.Coupon, float64, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	
	coupon, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, 0, errors.New(errors.ErrNotFound, "coupon not found")
	}

	// Check validity
	if !coupon.IsValid() {
		return nil, 0, errors.New(errors.ErrInvalidInput, "coupon is not valid")
	}

	// Check minimum order
	if orderValue < coupon.MinOrderValue {
		return nil, 0, errors.New(errors.ErrInvalidInput, "order value below minimum")
	}

	// Check user usage limit
	usageCount, _ := s.repo.GetUserUsage(ctx, coupon.ID, userID)
	if usageCount >= coupon.MaxPerUser {
		return nil, 0, errors.New(errors.ErrInvalidInput, "coupon usage limit reached")
	}

	// Check applicable categories
	if len(coupon.ApplicableCategories) > 0 {
		found := false
		for _, cat := range categoryIDs {
			for _, applicable := range coupon.ApplicableCategories {
				if cat == applicable {
					found = true
					break
				}
			}
		}
		if !found {
			return nil, 0, errors.New(errors.ErrInvalidInput, "coupon not applicable to these categories")
		}
	}

	// Check applicable products
	if len(coupon.ApplicableProducts) > 0 {
		found := false
		for _, prod := range productIDs {
			for _, applicable := range coupon.ApplicableProducts {
				if prod == applicable {
					found = true
					break
				}
			}
		}
		if !found {
			return nil, 0, errors.New(errors.ErrInvalidInput, "coupon not applicable to these products")
		}
	}

	// Calculate discount
	discount := coupon.CalculateDiscount(orderValue)

	return coupon, discount, nil
}

// ApplyCoupon applies a coupon to an order
func (s *CouponService) ApplyCoupon(ctx context.Context, code, userID, orderID string, orderValue float64) (float64, error) {
	coupon, discount, err := s.ValidateCoupon(ctx, code, userID, orderValue, nil, nil)
	if err != nil {
		return 0, err
	}

	// Mark coupon as used
	coupon.Use()
	if err := s.repo.Update(ctx, coupon); err != nil {
		return 0, err
	}

	// Record usage
	usage := &domain.CouponUsage{
		ID:       uuid.New().String(),
		CouponID: coupon.ID,
		UserID:   userID,
		OrderID:  orderID,
		Discount: discount,
		UsedAt:   time.Now(),
	}
	if err := s.repo.SaveUsage(ctx, usage); err != nil {
		s.logger.Error(err, "failed to save coupon usage")
	}

	s.logger.Infof("Coupon %s applied: user=%s, order=%s, discount=%.2f", 
		code, userID, orderID, discount)
	return discount, nil
}

// GetActiveCoupons returns all active coupons
func (s *CouponService) GetActiveCoupons(ctx context.Context) ([]*domain.Coupon, error) {
	return s.repo.FindActive(ctx)
}

// DeactivateCoupon deactivates a coupon
func (s *CouponService) DeactivateCoupon(ctx context.Context, couponID string) error {
	coupon, err := s.repo.FindByID(ctx, couponID)
	if err != nil {
		return err
	}

	coupon.Status = domain.CouponStatusExpired
	return s.repo.Update(ctx, coupon)
}
