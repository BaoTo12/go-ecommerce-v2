package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/titan-commerce/backend/warehouse-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type WarehouseRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

type StockRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewWarehouseRepository(databaseURL string, logger *logger.Logger) (*WarehouseRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Warehouse PostgreSQL repository initialized")
	return &WarehouseRepository{db: db, logger: logger}, nil
}

func NewStockRepository(db *sql.DB, logger *logger.Logger) *StockRepository {
	return &StockRepository{db: db, logger: logger}
}

func (r *WarehouseRepository) Save(ctx context.Context, warehouse *domain.Warehouse) error {
	addressJSON, _ := json.Marshal(warehouse.Address)
	zonesJSON, _ := json.Marshal(warehouse.Zones)

	query := `
		INSERT INTO warehouses (id, name, code, address, zones, capacity, priority, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		warehouse.ID, warehouse.Name, warehouse.Code, addressJSON, zonesJSON,
		warehouse.Capacity, warehouse.Priority, warehouse.Status,
		warehouse.CreatedAt, warehouse.UpdatedAt,
	)
	return err
}

func (r *WarehouseRepository) FindByID(ctx context.Context, warehouseID string) (*domain.Warehouse, error) {
	query := `SELECT id, name, code, address, zones, capacity, priority, status, created_at, updated_at FROM warehouses WHERE id = $1`
	
	var w domain.Warehouse
	var addressJSON, zonesJSON []byte
	
	err := r.db.QueryRowContext(ctx, query, warehouseID).Scan(
		&w.ID, &w.Name, &w.Code, &addressJSON, &zonesJSON, &w.Capacity,
		&w.Priority, &w.Status, &w.CreatedAt, &w.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "warehouse not found")
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(addressJSON, &w.Address)
	json.Unmarshal(zonesJSON, &w.Zones)
	return &w, nil
}

func (r *WarehouseRepository) FindByCode(ctx context.Context, code string) (*domain.Warehouse, error) {
	query := `SELECT id, name, code, address, zones, capacity, priority, status, created_at, updated_at FROM warehouses WHERE code = $1`
	
	var w domain.Warehouse
	var addressJSON, zonesJSON []byte
	
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&w.ID, &w.Name, &w.Code, &addressJSON, &zonesJSON, &w.Capacity,
		&w.Priority, &w.Status, &w.CreatedAt, &w.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "warehouse not found")
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(addressJSON, &w.Address)
	json.Unmarshal(zonesJSON, &w.Zones)
	return &w, nil
}

func (r *WarehouseRepository) List(ctx context.Context, status domain.WarehouseStatus, page, pageSize int) ([]*domain.Warehouse, int, error) {
	offset := (page - 1) * pageSize
	
	query := `SELECT id, name, code, address, zones, capacity, priority, status, created_at, updated_at 
			  FROM warehouses WHERE status = $1 ORDER BY priority ASC LIMIT $2 OFFSET $3`
	
	rows, err := r.db.QueryContext(ctx, query, status, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var warehouses []*domain.Warehouse
	for rows.Next() {
		var w domain.Warehouse
		var addressJSON, zonesJSON []byte
		
		if err := rows.Scan(&w.ID, &w.Name, &w.Code, &addressJSON, &zonesJSON,
			&w.Capacity, &w.Priority, &w.Status, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, 0, err
		}
		
		json.Unmarshal(addressJSON, &w.Address)
		json.Unmarshal(zonesJSON, &w.Zones)
		warehouses = append(warehouses, &w)
	}

	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM warehouses WHERE status = $1`, status).Scan(&total)

	return warehouses, total, nil
}

func (r *WarehouseRepository) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	addressJSON, _ := json.Marshal(warehouse.Address)
	zonesJSON, _ := json.Marshal(warehouse.Zones)

	query := `
		UPDATE warehouses 
		SET name = $2, address = $3, zones = $4, capacity = $5, priority = $6, status = $7, updated_at = $8
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		warehouse.ID, warehouse.Name, addressJSON, zonesJSON, warehouse.Capacity,
		warehouse.Priority, warehouse.Status, warehouse.UpdatedAt,
	)
	return err
}

func (r *WarehouseRepository) Delete(ctx context.Context, warehouseID string) error {
	query := `DELETE FROM warehouses WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, warehouseID)
	return err
}

// StockRepository implementations
func (r *StockRepository) SaveStock(ctx context.Context, stock *domain.WarehouseStock) error {
	query := `
		INSERT INTO warehouse_stocks (warehouse_id, product_id, available, reserved, zone, bin, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (warehouse_id, product_id) DO UPDATE 
		SET available = $3, reserved = $4, zone = $5, bin = $6, updated_at = $7
	`
	_, err := r.db.ExecContext(ctx, query,
		stock.WarehouseID, stock.ProductID, stock.Available, stock.Reserved,
		stock.Zone, stock.Bin, stock.UpdatedAt,
	)
	return err
}

func (r *StockRepository) GetStock(ctx context.Context, warehouseID, productID string) (*domain.WarehouseStock, error) {
	query := `SELECT warehouse_id, product_id, available, reserved, zone, bin, updated_at 
			  FROM warehouse_stocks WHERE warehouse_id = $1 AND product_id = $2`
	
	var stock domain.WarehouseStock
	err := r.db.QueryRowContext(ctx, query, warehouseID, productID).Scan(
		&stock.WarehouseID, &stock.ProductID, &stock.Available, &stock.Reserved,
		&stock.Zone, &stock.Bin, &stock.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "stock not found")
	}
	return &stock, err
}

func (r *StockRepository) GetStocksByWarehouse(ctx context.Context, warehouseID string) ([]*domain.WarehouseStock, error) {
	query := `SELECT warehouse_id, product_id, available, reserved, zone, bin, updated_at 
			  FROM warehouse_stocks WHERE warehouse_id = $1`
	
	rows, err := r.db.QueryContext(ctx, query, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []*domain.WarehouseStock
	for rows.Next() {
		var stock domain.WarehouseStock
		if err := rows.Scan(&stock.WarehouseID, &stock.ProductID, &stock.Available,
			&stock.Reserved, &stock.Zone, &stock.Bin, &stock.UpdatedAt); err != nil {
			return nil, err
		}
		stocks = append(stocks, &stock)
	}
	return stocks, nil
}

func (r *StockRepository) UpdateStock(ctx context.Context, stock *domain.WarehouseStock) error {
	query := `
		UPDATE warehouse_stocks 
		SET available = $3, reserved = $4, zone = $5, bin = $6, updated_at = $7
		WHERE warehouse_id = $1 AND product_id = $2
	`
	_, err := r.db.ExecContext(ctx, query,
		stock.WarehouseID, stock.ProductID, stock.Available, stock.Reserved,
		stock.Zone, stock.Bin, stock.UpdatedAt,
	)
	return err
}

func (r *StockRepository) SaveMovement(ctx context.Context, movement *domain.StockMovement) error {
	query := `
		INSERT INTO stock_movements (movement_id, warehouse_id, product_id, movement_type, quantity, from_zone, to_zone, reference_id, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		movement.MovementID, movement.WarehouseID, movement.ProductID, movement.MovementType,
		movement.Quantity, movement.FromZone, movement.ToZone, movement.ReferenceID,
		movement.Notes, movement.CreatedAt,
	)
	return err
}

func (r *StockRepository) GetMovementHistory(ctx context.Context, warehouseID, productID string) ([]*domain.StockMovement, error) {
	query := `SELECT movement_id, warehouse_id, product_id, movement_type, quantity, from_zone, to_zone, reference_id, notes, created_at
			  FROM stock_movements WHERE warehouse_id = $1 AND product_id = $2 ORDER BY created_at DESC LIMIT 100`
	
	rows, err := r.db.QueryContext(ctx, query, warehouseID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []*domain.StockMovement
	for rows.Next() {
		var m domain.StockMovement
		if err := rows.Scan(&m.MovementID, &m.WarehouseID, &m.ProductID, &m.MovementType,
			&m.Quantity, &m.FromZone, &m.ToZone, &m.ReferenceID, &m.Notes, &m.CreatedAt); err != nil {
			return nil, err
		}
		movements = append(movements, &m)
	}
	return movements, nil
}

func (r *WarehouseRepository) Close() error {
	return r.db.Close()
}
