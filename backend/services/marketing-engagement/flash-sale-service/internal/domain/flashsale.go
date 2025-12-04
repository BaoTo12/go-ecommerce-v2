package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type FlashSale struct {
	ID             string
	ProductID      string
	ProductName    string
	OriginalPrice  float64
	FlashPrice     float64
	TotalStock     int
	RemainingStock int
	StartTime      time.Time
	EndTime        time.Time
	IsActive       bool
	CreatedAt      time.Time
}

func NewFlashSale(productID, productName string, originalPrice, flashPrice float64, stock int, startTime, endTime time.Time) (*FlashSale, error) {
	if productID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "product ID is required")
	}
	if flashPrice >= originalPrice {
		return nil, errors.New(errors.ErrInvalidInput, "flash price must be less than original price")
	}
	if stock <= 0 {
		return nil, errors.New(errors.ErrInvalidInput, "stock must be positive")
	}
	if endTime.Before(startTime) {
		return nil, errors.New(errors.ErrInvalidInput, "end time must be after start time")
	}

	return &FlashSale{
		ID:             uuid.New().String(),
		ProductID:      productID,
		ProductName:    productName,
		OriginalPrice:  originalPrice,
		FlashPrice:     flashPrice,
		TotalStock:     stock,
		RemainingStock: stock,
		StartTime:      startTime,
		EndTime:        endTime,
		IsActive:       time.Now().After(startTime) && time.Now().Before(endTime),
		CreatedAt:      time.Now(),
	}, nil
}

func (f *FlashSale) IsLive() bool {
	now := time.Now()
	return now.After(f.StartTime) && now.Before(f.EndTime) && f.RemainingStock > 0
}

func (f *FlashSale) GetDiscountPercentage() int {
	return int(((f.OriginalPrice - f.FlashPrice) / f.OriginalPrice) * 100)
}
