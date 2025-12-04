package domain

import (
	"time"

	"github.com/google/uuid"
)

type CouponType string
type CouponStatus string

const (
	CouponTypePercentage CouponType = "PERCENTAGE"
	CouponTypeFixed      CouponType = "FIXED"
	CouponTypeFreeShip   CouponType = "FREE_SHIPPING"

	CouponStatusActive   CouponStatus = "ACTIVE"
	CouponStatusExpired  CouponStatus = "EXPIRED"
	CouponStatusDepleted CouponStatus = "DEPLETED"
)

type Coupon struct {
	ID              string
	Code            string
	Name            string
	Description     string
	Type            CouponType
	Value           float64     // discount % or fixed amount
	MinOrderValue   float64     // minimum order to apply
	MaxDiscount     float64     // cap for percentage discounts
	TotalQuantity   int
	UsedQuantity    int
	MaxPerUser      int
	ValidFrom       time.Time
	ValidUntil      time.Time
	ApplicableCategories []string // empty = all
	ApplicableProducts   []string // empty = all
	Status          CouponStatus
	CreatedAt       time.Time
}

type CouponUsage struct {
	ID        string
	CouponID  string
	UserID    string
	OrderID   string
	Discount  float64
	UsedAt    time.Time
}

func NewCoupon(code, name, description string, couponType CouponType, value, minOrder, maxDiscount float64, totalQty, maxPerUser int, validFrom, validUntil time.Time) *Coupon {
	return &Coupon{
		ID:              uuid.New().String(),
		Code:            code,
		Name:            name,
		Description:     description,
		Type:            couponType,
		Value:           value,
		MinOrderValue:   minOrder,
		MaxDiscount:     maxDiscount,
		TotalQuantity:   totalQty,
		UsedQuantity:    0,
		MaxPerUser:      maxPerUser,
		ValidFrom:       validFrom,
		ValidUntil:      validUntil,
		Status:          CouponStatusActive,
		CreatedAt:       time.Now(),
	}
}

func (c *Coupon) IsValid() bool {
	now := time.Now()
	return c.Status == CouponStatusActive && 
		now.After(c.ValidFrom) && 
		now.Before(c.ValidUntil) &&
		c.UsedQuantity < c.TotalQuantity
}

func (c *Coupon) CalculateDiscount(orderValue float64) float64 {
	if orderValue < c.MinOrderValue {
		return 0
	}

	var discount float64
	switch c.Type {
	case CouponTypePercentage:
		discount = orderValue * (c.Value / 100)
		if c.MaxDiscount > 0 && discount > c.MaxDiscount {
			discount = c.MaxDiscount
		}
	case CouponTypeFixed:
		discount = c.Value
	case CouponTypeFreeShip:
		discount = c.Value // shipping cost
	}

	return discount
}

func (c *Coupon) Use() {
	c.UsedQuantity++
	if c.UsedQuantity >= c.TotalQuantity {
		c.Status = CouponStatusDepleted
	}
}

func (c *Coupon) CheckExpiry() {
	if time.Now().After(c.ValidUntil) {
		c.Status = CouponStatusExpired
	}
}

type Repository interface {
	Save(ctx interface{}, coupon *Coupon) error
	FindByID(ctx interface{}, couponID string) (*Coupon, error)
	FindByCode(ctx interface{}, code string) (*Coupon, error)
	Update(ctx interface{}, coupon *Coupon) error
	FindActive(ctx interface{}) ([]*Coupon, error)
	SaveUsage(ctx interface{}, usage *CouponUsage) error
	GetUserUsage(ctx interface{}, couponID, userID string) (int, error)
}
