package domain

import "context"

type FlashSaleRepository interface {
	SaveFlashSale(ctx context.Context, sale *FlashSale) error
	GetFlashSale(ctx context.Context, flashSaleID string) (*FlashSale, error)
	UpdateFlashSale(ctx context.Context, sale *FlashSale) error
	GetActiveFlashSales(ctx context.Context) ([]*FlashSale, error)
	SaveProduct(ctx context.Context, product *FlashSaleProduct) error
	GetProducts(ctx context.Context, flashSaleID string) ([]*FlashSaleProduct, error)
}

