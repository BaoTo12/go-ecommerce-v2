package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type CouponType string

const (
	CouponTypePercentage   CouponType = "PERCENTAGE"
	CouponTypeFixedAmount  CouponType = "FIXED_AMOUNT"
	CouponTypeFreeShipping CouponType = "FREE_SHIPPING"
)

type Coupon struct {
	ID            string
	Code          string
	Type          CouponType
	DiscountValue float64
	MinPurchase   float64
	UsageLimit    int
	UsedCount     int
	ValidFrom     time.Time
	ValidUntil    time.Time
	Active        bool
	CreatedAt     time.Time
}

func NewCoupon(code string, couponType CouponType, discountValue, minPurchase float64, usageLimit int, validFrom, validUntil time.Time) (*Coupon, error) {
	if code == "" {
		return nil, errors.New(errors.ErrInvalidInput, "coupon code is required")
	}
	if discountValue <= 0 {
		return nil, errors.New(errors.ErrInvalidInput, "discount value must be positive")
	}
	if validUntil.Before(validFrom) {
		return nil, errors.New(errors.ErrInvalidInput, "valid until must be after valid from")
	}

	return &Coupon{
		ID:            uuid.New().String(),
		Code:          code,
		Type:          couponType,
		DiscountValue: discountValue,
		MinPurchase:   minPurchase,
		UsageLimit:    usageLimit,
		UsedCount:     0,
		ValidFrom:     validFrom,
		ValidUntil:    validUntil,
		Active:        true,
		CreatedAt:     time.Now(),
	}, nil
}

func (c *Coupon) IsValid(orderTotal float64) (bool, string) {
	now := time.Now()

	if !c.Active {
		return false, "coupon is inactive"
	}
	if now.Before(c.ValidFrom) {
		return false, "coupon not yet valid"
	}
	if now.After(c.ValidUntil) {
		return false, "coupon has expired"
	}
	if c.UsageLimit > 0 && c.UsedCount >= c.UsageLimit {
		return false, "coupon usage limit reached"
	}
	if orderTotal < c.MinPurchase {
		return false, "order total does not meet minimum purchase requirement"
	}

	return true, ""
}

func (c *Coupon) CalculateDiscount(orderTotal float64) float64 {
	switch c.Type {
	case CouponTypePercentage:
		return orderTotal * (c.DiscountValue / 100.0)
	case CouponTypeFixedAmount:
		return c.DiscountValue
	case CouponTypeFreeShipping:
		return 0.0 // Shipping discount handled separately
	default:
		return 0.0
	}
}

func (c *Coupon) IncrementUsage() {
	c.UsedCount++
}
