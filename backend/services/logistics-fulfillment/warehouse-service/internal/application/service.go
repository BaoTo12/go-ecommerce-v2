package application

import (
	"context"
	"math"
	"sort"

	"github.com/titan-commerce/backend/warehouse-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type WarehouseService struct {
	warehouseRepo domain.WarehouseRepository
	stockRepo     domain.StockRepository
	logger        *logger.Logger
}

func NewWarehouseService(warehouseRepo domain.WarehouseRepository, stockRepo domain.StockRepository, logger *logger.Logger) *WarehouseService {
	return &WarehouseService{
		warehouseRepo: warehouseRepo,
		stockRepo:     stockRepo,
		logger:        logger,
	}
}

// CreateWarehouse creates a new warehouse (Command)
func (s *WarehouseService) CreateWarehouse(ctx context.Context, name, code string, address domain.Address, capacity, priority int) (*domain.Warehouse, error) {
	warehouse, err := domain.NewWarehouse(name, code, address, capacity, priority)
	if err != nil {
		return nil, err
	}

	if err := s.warehouseRepo.Save(ctx, warehouse); err != nil {
		s.logger.Error(err, "failed to save warehouse")
		return nil, err
	}

	s.logger.Infof("Warehouse created: id=%s, name=%s, code=%s", warehouse.ID, warehouse.Name, warehouse.Code)
	return warehouse, nil
}

// GetWarehouse retrieves warehouse details (Query)
func (s *WarehouseService) GetWarehouse(ctx context.Context, warehouseID string) (*domain.Warehouse, error) {
	return s.warehouseRepo.FindByID(ctx, warehouseID)
}

// ListWarehouses lists warehouses with pagination (Query)
func (s *WarehouseService) ListWarehouses(ctx context.Context, status domain.WarehouseStatus, page, pageSize int) ([]*domain.Warehouse, int, error) {
	return s.warehouseRepo.List(ctx, status, page, pageSize)
}

// UpdateWarehouse updates warehouse details (Command)
func (s *WarehouseService) UpdateWarehouse(ctx context.Context, warehouseID, name string, status domain.WarehouseStatus, capacity, priority int) (*domain.Warehouse, error) {
	warehouse, err := s.warehouseRepo.FindByID(ctx, warehouseID)
	if err != nil {
		return nil, err
	}

	warehouse.Update(name, capacity, priority)
	warehouse.UpdateStatus(status)

	if err := s.warehouseRepo.Update(ctx, warehouse); err != nil {
		s.logger.Error(err, "failed to update warehouse")
		return nil, err
	}

	s.logger.Infof("Warehouse updated: id=%s, status=%s", warehouseID, status)
	return warehouse, nil
}

// AllocateStock allocates stock from warehouses based on priority and proximity (Command)
func (s *WarehouseService) AllocateStock(ctx context.Context, productID string, requiredQuantity int, customerAddress string, preferredWarehouseIDs []string) ([]domain.StockAllocation, error) {
	// Get all active warehouses
	warehouses, _, err := s.warehouseRepo.List(ctx, domain.StatusActive, 1, 100)
	if err != nil {
		return nil, err
	}

	// Sort by priority (lower number = higher priority)
	sort.Slice(warehouses, func(i, j int) bool {
		return warehouses[i].Priority < warehouses[j].Priority
	})

	var allocations []domain.StockAllocation
	remainingQuantity := requiredQuantity

	// Try to allocate from warehouses
	for _, warehouse := range warehouses {
		if remainingQuantity <= 0 {
			break
		}

		// Check stock availability
		stock, err := s.stockRepo.GetStock(ctx, warehouse.ID, productID)
		if err != nil || stock == nil || stock.AvailableQuantity <= 0 {
			continue
		}

		// Allocate available quantity
		allocQuantity := stock.AvailableQuantity
		if allocQuantity > remainingQuantity {
			allocQuantity = remainingQuantity
		}

		allocations = append(allocations, domain.StockAllocation{
			WarehouseID: warehouse.ID,
			Quantity:    allocQuantity,
		})

		remainingQuantity -= allocQuantity
	}

	if remainingQuantity > 0 {
		s.logger.Warnf("Insufficient stock: product=%s, required=%d, allocated=%d",
			productID, requiredQuantity, requiredQuantity-remainingQuantity)
	}

	return allocations, nil
}

// GetStockByWarehouse retrieves stock information for a warehouse (Query)
func (s *WarehouseService) GetStockByWarehouse(ctx context.Context, warehouseID, productID string) ([]*domain.WarehouseStock, error) {
	if productID != "" {
		stock, err := s.stockRepo.GetStock(ctx, warehouseID, productID)
		if err != nil {
			return nil, err
		}
		return []*domain.WarehouseStock{stock}, nil
	}
	return s.stockRepo.GetStocksByWarehouse(ctx, warehouseID)
}

// TransferStock transfers stock between warehouses (Command)
func (s *WarehouseService) TransferStock(ctx context.Context, fromWarehouseID, toWarehouseID, productID string, quantity int, notes string) (string, error) {
	// Get source stock
	fromStock, err := s.stockRepo.GetStock(ctx, fromWarehouseID, productID)
	if err != nil {
		return "", err
	}

	// Remove from source
	if !fromStock.RemoveStock(quantity) {
		s.logger.Errorf("Insufficient stock for transfer: warehouse=%s, product=%s, available=%d, required=%d",
			fromWarehouseID, productID, fromStock.AvailableQuantity, quantity)
		return "", nil
	}

	// Add to destination
	toStock, err := s.stockRepo.GetStock(ctx, toWarehouseID, productID)
	if err != nil {
		// Create new stock entry if doesn't exist
		toStock = domain.NewWarehouseStock(toWarehouseID, productID, 0, "", "")
	}
	toStock.AddStock(quantity)

	// Update both stocks
	if err := s.stockRepo.UpdateStock(ctx, fromStock); err != nil {
		return "", err
	}
	if err := s.stockRepo.UpdateStock(ctx, toStock); err != nil {
		return "", err
	}

	// Record movements
	outboundMovement := domain.NewStockMovement(fromWarehouseID, productID, domain.MovementTransfer, quantity, toWarehouseID, notes, "system")
	inboundMovement := domain.NewStockMovement(toWarehouseID, productID, domain.MovementTransfer, quantity, fromWarehouseID, notes, "system")

	if err := s.stockRepo.SaveMovement(ctx, outboundMovement); err != nil {
		return "", err
	}
	if err := s.stockRepo.SaveMovement(ctx, inboundMovement); err != nil {
		return "", err
	}

	s.logger.Infof("Stock transferred: from=%s, to=%s, product=%s, quantity=%d",
		fromWarehouseID, toWarehouseID, productID, quantity)

	return outboundMovement.MovementID, nil
}

// RecordStockMovement records a stock movement (Command)
func (s *WarehouseService) RecordStockMovement(ctx context.Context, warehouseID, productID string, movementType domain.StockMovementType, quantity int, referenceID, notes string) (*domain.StockMovement, error) {
	// Get stock
	stock, err := s.stockRepo.GetStock(ctx, warehouseID, productID)
	if err != nil {
		// Create new stock entry if doesn't exist
		stock = domain.NewWarehouseStock(warehouseID, productID, 0, "", "")
	}

	// Update stock based on movement type
	switch movementType {
	case domain.MovementInbound, domain.MovementReturn:
		stock.AddStock(quantity)
	case domain.MovementOutbound:
		if !stock.RemoveStock(quantity) {
			s.logger.Errorf("Insufficient stock for outbound movement: warehouse=%s, product=%s",
				warehouseID, productID)
			return nil, nil
		}
	case domain.MovementAdjustment:
		// Adjustment can be positive or negative
		stock.AvailableQuantity += quantity
		stock.TotalQuantity += quantity
	}

	// Update stock
	if err := s.stockRepo.UpdateStock(ctx, stock); err != nil {
		return nil, err
	}

	// Create movement record
	movement := domain.NewStockMovement(warehouseID, productID, movementType, quantity, referenceID, notes, "system")
	if err := s.stockRepo.SaveMovement(ctx, movement); err != nil {
		return nil, err
	}

	s.logger.Infof("Stock movement recorded: warehouse=%s, product=%s, type=%s, quantity=%d",
		warehouseID, productID, movementType, quantity)

	return movement, nil
}

// calculateDistance calculates distance between two points (simple haversine formula)
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // km

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

