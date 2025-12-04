package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

// CouponType represents the type of coupon
type CouponType string

const (
	CouponPercentage CouponType = "PERCENTAGE"
	CouponFixed      CouponType = "FIXED"
	CouponFreeShip   CouponType = "FREE_SHIPPING"
	CouponBOGO       CouponType = "BOGO"
)

// Coupon represents a discount coupon
type Coupon struct {
	CouponID       string
	Code           string
	Type           CouponType
	DiscountValue  int64 // percentage * 100 or fixed amount in cents
	MinOrderValue  int64
	MaxDiscount    int64
	UsageLimit     int
	UsageCount     int
	PerUserLimit   int
	ValidFrom      time.Time
	ValidUntil     time.Time
	IsActive       bool
	ApplicableProducts []string
	ApplicableCategories []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewCoupon creates a new coupon
func NewCoupon(code string, couponType CouponType, discountValue, minOrderValue int64) (*Coupon, error) {
	if code == "" {
		return nil, errors.New(errors.ErrInvalidInput, "coupon code is required")
	}
	if discountValue <= 0 {
		return nil, errors.New(errors.ErrInvalidInput, "discount value must be positive")
	}

	now := time.Now()
	return &Coupon{
		CouponID:      uuid.New().String(),
		Code:          code,
		Type:          couponType,
		DiscountValue: discountValue,
		MinOrderValue: minOrderValue,
		UsageLimit:    -1, // unlimited
		UsageCount:    0,
		PerUserLimit:  -1, // unlimited
		ValidFrom:     now,
		ValidUntil:    now.AddDate(0, 1, 0), // 1 month default
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// IsValid checks if coupon is currently valid
func (c *Coupon) IsValid() bool {
	now := time.Now()
	return c.IsActive &&
		now.After(c.ValidFrom) &&
		now.Before(c.ValidUntil) &&
		(c.UsageLimit == -1 || c.UsageCount < c.UsageLimit)
}

// CanApplyToOrder checks if coupon can be applied to an order
func (c *Coupon) CanApplyToOrder(orderValue int64) bool {
	return c.IsValid() && orderValue >= c.MinOrderValue
}

// CalculateDiscount calculates the discount amount
func (c *Coupon) CalculateDiscount(orderValue int64) int64 {
	if !c.CanApplyToOrder(orderValue) {
		return 0
	}

	var discount int64

	switch c.Type {
	case CouponPercentage:
		discount = (orderValue * c.DiscountValue) / 10000 // discount value is percentage * 100
		if c.MaxDiscount > 0 && discount > c.MaxDiscount {
			discount = c.MaxDiscount
		}
	case CouponFixed:
		discount = c.DiscountValue
		if discount > orderValue {
			discount = orderValue
		}
	case CouponFreeShip:
		// Handled by order service
		discount = 0
	}

	return discount
}

// Use increments usage count
func (c *Coupon) Use() error {
	if !c.IsValid() {
		return errors.New(errors.ErrInvalidInput, "coupon is not valid")
	}

	c.UsageCount++
	c.UpdatedAt = time.Now()
	return nil
}

// CouponUsage tracks individual coupon usage
type CouponUsage struct {
	UsageID    string
	CouponID   string
	UserID     string
	OrderID    string
	Discount   int64
	UsedAt     time.Time
}

// NewCouponUsage creates a new usage record
func NewCouponUsage(couponID, userID, orderID string, discount int64) *CouponUsage {
	return &CouponUsage{
		UsageID:  uuid.New().String(),
		CouponID: couponID,
		UserID:   userID,
		OrderID:  orderID,
		Discount: discount,
		UsedAt:   time.Now(),
	}
}

