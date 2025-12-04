package domain

import "context"

type WarehouseRepository interface {
	Save(ctx context.Context, warehouse *Warehouse) error
	FindByID(ctx context.Context, warehouseID string) (*Warehouse, error)
	FindByCode(ctx context.Context, code string) (*Warehouse, error)
	List(ctx context.Context, status WarehouseStatus, page, pageSize int) ([]*Warehouse, int, error)
	Update(ctx context.Context, warehouse *Warehouse) error
	Delete(ctx context.Context, warehouseID string) error
}

type StockRepository interface {
	SaveStock(ctx context.Context, stock *WarehouseStock) error
	GetStock(ctx context.Context, warehouseID, productID string) (*WarehouseStock, error)
	GetStocksByWarehouse(ctx context.Context, warehouseID string) ([]*WarehouseStock, error)
	UpdateStock(ctx context.Context, stock *WarehouseStock) error
	SaveMovement(ctx context.Context, movement *StockMovement) error
	GetMovementHistory(ctx context.Context, warehouseID, productID string) ([]*StockMovement, error)
}

