package domain

import (
	"time"

	"github.com/google/uuid"
)

type FlashSaleStatus string

const (
	FlashSaleStatusScheduled FlashSaleStatus = "SCHEDULED"
	FlashSaleStatusActive    FlashSaleStatus = "ACTIVE"
	FlashSaleStatusEnded     FlashSaleStatus = "ENDED"
	FlashSaleStatusSoldOut   FlashSaleStatus = "SOLD_OUT"
)

type FlashSale struct {
	ID              string
	ProductID       string
	OriginalPrice   float64
	SalePrice       float64
	DiscountPercent int
	TotalQuantity   int
	SoldQuantity    int
	MaxPerUser      int
	StartTime       time.Time
	EndTime         time.Time
	Status          FlashSaleStatus
	CreatedAt       time.Time
}

type FlashSaleReservation struct {
	ID          string
	FlashSaleID string
	UserID      string
	Quantity    int
	Status      string // "reserved", "confirmed", "expired"
	ExpiresAt   time.Time
	CreatedAt   time.Time
}

type FlashSalePurchase struct {
	ID          string
	FlashSaleID string
	UserID      string
	Quantity    int
	PurchasedAt time.Time
}

func NewFlashSale(productID string, originalPrice, salePrice float64, totalQty, maxPerUser int, start, end time.Time) *FlashSale {
	discountPercent := int((1 - salePrice/originalPrice) * 100)
	return &FlashSale{
		ID:              uuid.New().String(),
		ProductID:       productID,
		OriginalPrice:   originalPrice,
		SalePrice:       salePrice,
		DiscountPercent: discountPercent,
		TotalQuantity:   totalQty,
		SoldQuantity:    0,
		MaxPerUser:      maxPerUser,
		StartTime:       start,
		EndTime:         end,
		Status:          FlashSaleStatusScheduled,
		CreatedAt:       time.Now(),
	}
}

func (fs *FlashSale) IsActive() bool {
	now := time.Now()
	return fs.Status == FlashSaleStatusActive && now.After(fs.StartTime) && now.Before(fs.EndTime)
}

func (fs *FlashSale) CanPurchase(quantity int) bool {
	return fs.IsActive() && fs.SoldQuantity+quantity <= fs.TotalQuantity
}

func (fs *FlashSale) RemainingQuantity() int {
	return fs.TotalQuantity - fs.SoldQuantity
}

func (fs *FlashSale) Activate() {
	fs.Status = FlashSaleStatusActive
}

func (fs *FlashSale) End() {
	if fs.SoldQuantity >= fs.TotalQuantity {
		fs.Status = FlashSaleStatusSoldOut
	} else {
		fs.Status = FlashSaleStatusEnded
	}
}

func (fs *FlashSale) Purchase(quantity int) {
	fs.SoldQuantity += quantity
	if fs.SoldQuantity >= fs.TotalQuantity {
		fs.Status = FlashSaleStatusSoldOut
	}
}

func NewReservation(saleID, userID string, quantity int, ttl time.Duration) *FlashSaleReservation {
	return &FlashSaleReservation{
		ID:          uuid.New().String(),
		FlashSaleID: saleID,
		UserID:      userID,
		Quantity:    quantity,
		Status:      "reserved",
		ExpiresAt:   time.Now().Add(ttl),
		CreatedAt:   time.Now(),
	}
}

type Repository interface {
	Save(ctx interface{}, sale *FlashSale) error
	FindByID(ctx interface{}, saleID string) (*FlashSale, error)
	FindActive(ctx interface{}) ([]*FlashSale, error)
	FindUpcoming(ctx interface{}) ([]*FlashSale, error)
	Update(ctx interface{}, sale *FlashSale) error
	SaveReservation(ctx interface{}, res *FlashSaleReservation) error
	SavePurchase(ctx interface{}, purchase *FlashSalePurchase) error
	GetUserPurchases(ctx interface{}, saleID, userID string) (int, error)
	DecrementStock(ctx interface{}, saleID string, quantity int) (bool, error)
}
