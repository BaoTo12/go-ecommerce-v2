package domain

import (
	"time"

	"github.com/google/uuid"
)

type StockMovementType string

const (
	MovementInbound    StockMovementType = "INBOUND"
	MovementOutbound   StockMovementType = "OUTBOUND"
	MovementTransfer   StockMovementType = "TRANSFER"
	MovementAdjustment StockMovementType = "ADJUSTMENT"
	MovementReturn     StockMovementType = "RETURN"
)

type WarehouseStock struct {
	WarehouseID       string
	ProductID         string
	AvailableQuantity int
	ReservedQuantity  int
	TotalQuantity     int
	Zone              string
	BinLocation       string
	LastUpdated       time.Time
}

type StockMovement struct {
	MovementID   string
	WarehouseID  string
	ProductID    string
	MovementType StockMovementType
	Quantity     int
	ReferenceID  string // Order ID, Transfer ID, etc.
	Notes        string
	CreatedAt    time.Time
	CreatedBy    string
}

func NewWarehouseStock(warehouseID, productID string, quantity int, zone, binLocation string) *WarehouseStock {
	return &WarehouseStock{
		WarehouseID:       warehouseID,
		ProductID:         productID,
		AvailableQuantity: quantity,
		ReservedQuantity:  0,
		TotalQuantity:     quantity,
		Zone:              zone,
		BinLocation:       binLocation,
		LastUpdated:       time.Now(),
	}
}

func NewStockMovement(warehouseID, productID string, movementType StockMovementType, quantity int, referenceID, notes, createdBy string) *StockMovement {
	return &StockMovement{
		MovementID:   uuid.New().String(),
		WarehouseID:  warehouseID,
		ProductID:    productID,
		MovementType: movementType,
		Quantity:     quantity,
		ReferenceID:  referenceID,
		Notes:        notes,
		CreatedAt:    time.Now(),
		CreatedBy:    createdBy,
	}
}

func (ws *WarehouseStock) AddStock(quantity int) {
	ws.AvailableQuantity += quantity
	ws.TotalQuantity += quantity
	ws.LastUpdated = time.Now()
}

func (ws *WarehouseStock) RemoveStock(quantity int) bool {
	if ws.AvailableQuantity < quantity {
		return false
	}
	ws.AvailableQuantity -= quantity
	ws.TotalQuantity -= quantity
	ws.LastUpdated = time.Now()
	return true
}

func (ws *WarehouseStock) ReserveStock(quantity int) bool {
	if ws.AvailableQuantity < quantity {
		return false
	}
	ws.AvailableQuantity -= quantity
	ws.ReservedQuantity += quantity
	ws.LastUpdated = time.Now()
	return true
}

func (ws *WarehouseStock) ReleaseStock(quantity int) {
	ws.AvailableQuantity += quantity
	ws.ReservedQuantity -= quantity
	ws.LastUpdated = time.Now()
}

type StockAllocation struct {
	WarehouseID string
	Quantity    int
}

