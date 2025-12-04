package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

// FlashSaleStatus represents the status of a flash sale
type FlashSaleStatus string

const (
	FlashSaleScheduled FlashSaleStatus = "SCHEDULED"
	FlashSaleLive      FlashSaleStatus = "LIVE"
	FlashSaleEnded     FlashSaleStatus = "ENDED"
	FlashSaleSoldOut   FlashSaleStatus = "SOLD_OUT"
)

// FlashSale represents a limited-time high-concurrency sale event
type FlashSale struct {
	FlashSaleID      string
	Name             string
	Description      string
	Status           FlashSaleStatus
	StartTime        time.Time
	EndTime          time.Time
	TotalStock       int
	SoldCount        int
	ViewCount        int
	AddToCartCount   int
	ConversionRate   float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// FlashSaleProduct represents a product in a flash sale
type FlashSaleProduct struct {
	ProductID     string
	FlashSaleID   string
	OriginalPrice int64
	FlashPrice    int64
	Stock         int
	MaxPerUser    int
	SoldCount     int
}

// NewFlashSale creates a new flash sale
func NewFlashSale(name, description string, startTime, endTime time.Time, totalStock int) (*FlashSale, error) {
	if name == "" {
		return nil, errors.New(errors.ErrInvalidInput, "flash sale name is required")
	}
	if startTime.After(endTime) {
		return nil, errors.New(errors.ErrInvalidInput, "start time must be before end time")
	}
	if totalStock <= 0 {
		return nil, errors.New(errors.ErrInvalidInput, "total stock must be positive")
	}

	now := time.Now()
	return &FlashSale{
		FlashSaleID:    uuid.New().String(),
		Name:           name,
		Description:    description,
		Status:         FlashSaleScheduled,
		StartTime:      startTime,
		EndTime:        endTime,
		TotalStock:     totalStock,
		SoldCount:      0,
		ViewCount:      0,
		AddToCartCount: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// Start starts the flash sale
func (f *FlashSale) Start() error {
	if time.Now().Before(f.StartTime) {
		return errors.New(errors.ErrInvalidInput, "flash sale start time not reached")
	}

	f.Status = FlashSaleLive
	f.UpdatedAt = time.Now()
	return nil
}

// RecordSale records a sale
func (f *FlashSale) RecordSale(quantity int) error {
	if f.Status != FlashSaleLive {
		return errors.New(errors.ErrInvalidInput, "flash sale is not live")
	}

	if f.SoldCount+quantity > f.TotalStock {
		return errors.New(errors.ErrInvalidInput, "insufficient stock")
	}

	f.SoldCount += quantity
	f.UpdatedAt = time.Now()

	if f.SoldCount >= f.TotalStock {
		f.Status = FlashSaleSoldOut
	}

	f.calculateConversion()
	return nil
}

// RecordView increments view count
func (f *FlashSale) RecordView() {
	f.ViewCount++
}

// RecordAddToCart increments add to cart count
func (f *FlashSale) RecordAddToCart() {
	f.AddToCartCount++
}

func (f *FlashSale) calculateConversion() {
	if f.ViewCount > 0 {
		f.ConversionRate = float64(f.SoldCount) / float64(f.ViewCount) * 100
	}
}

// IsActive checks if flash sale is currently active
func (f *FlashSale) IsActive() bool {
	now := time.Now()
	return f.Status == FlashSaleLive &&
		now.After(f.StartTime) &&
		now.Before(f.EndTime) &&
		f.SoldCount < f.TotalStock
}

// GetRemainingStock returns available stock
func (f *FlashSale) GetRemainingStock() int {
	return f.TotalStock - f.SoldCount
}

