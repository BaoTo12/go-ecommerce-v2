package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/titan-commerce/backend/flash-sale-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type FlashSalePostgresRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewFlashSalePostgresRepository(databaseURL string, logger *logger.Logger) (*FlashSalePostgresRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to open database", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	repo := &FlashSalePostgresRepository{db: db, logger: logger}
	if err := repo.createTables(ctx); err != nil {
		return nil, err
	}

	logger.Info("Flash Sale PostgreSQL repository initialized")
	return repo, nil
}

func (r *FlashSalePostgresRepository) createTables(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS flash_sales (
			id VARCHAR(64) PRIMARY KEY,
			product_id VARCHAR(64) NOT NULL,
			original_price DECIMAL(12,2) NOT NULL,
			sale_price DECIMAL(12,2) NOT NULL,
			discount_percent INT NOT NULL,
			total_quantity INT NOT NULL,
			sold_quantity INT NOT NULL DEFAULT 0,
			max_per_user INT NOT NULL DEFAULT 5,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL,
			status VARCHAR(20) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_flash_status ON flash_sales(status, start_time)`,
		`CREATE TABLE IF NOT EXISTS flash_sale_reservations (
			id VARCHAR(64) PRIMARY KEY,
			flash_sale_id VARCHAR(64) NOT NULL REFERENCES flash_sales(id),
			user_id VARCHAR(64) NOT NULL,
			quantity INT NOT NULL,
			status VARCHAR(20) NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_reservation_user ON flash_sale_reservations(flash_sale_id, user_id)`,
		`CREATE TABLE IF NOT EXISTS flash_sale_purchases (
			id VARCHAR(64) PRIMARY KEY,
			flash_sale_id VARCHAR(64) NOT NULL REFERENCES flash_sales(id),
			user_id VARCHAR(64) NOT NULL,
			quantity INT NOT NULL,
			purchased_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_purchase_user ON flash_sale_purchases(flash_sale_id, user_id)`,
	}

	for _, query := range queries {
		if _, err := r.db.ExecContext(ctx, query); err != nil {
			return errors.Wrap(errors.ErrInternal, "failed to create table", err)
		}
	}
	return nil
}

func (r *FlashSalePostgresRepository) Save(ctx context.Context, sale *domain.FlashSale) error {
	query := `INSERT INTO flash_sales 
			  (id, product_id, original_price, sale_price, discount_percent, total_quantity,
			   sold_quantity, max_per_user, start_time, end_time, status, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.ExecContext(ctx, query,
		sale.ID, sale.ProductID, sale.OriginalPrice, sale.SalePrice, sale.DiscountPercent,
		sale.TotalQuantity, sale.SoldQuantity, sale.MaxPerUser, sale.StartTime, sale.EndTime,
		sale.Status, sale.CreatedAt)
	return err
}

func (r *FlashSalePostgresRepository) FindByID(ctx context.Context, saleID string) (*domain.FlashSale, error) {
	query := `SELECT id, product_id, original_price, sale_price, discount_percent,
			  total_quantity, sold_quantity, max_per_user, start_time, end_time, status, created_at
			  FROM flash_sales WHERE id = $1`

	var sale domain.FlashSale
	err := r.db.QueryRowContext(ctx, query, saleID).Scan(
		&sale.ID, &sale.ProductID, &sale.OriginalPrice, &sale.SalePrice, &sale.DiscountPercent,
		&sale.TotalQuantity, &sale.SoldQuantity, &sale.MaxPerUser, &sale.StartTime, &sale.EndTime,
		&sale.Status, &sale.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "flash sale not found")
	}
	return &sale, err
}

func (r *FlashSalePostgresRepository) FindActive(ctx context.Context) ([]*domain.FlashSale, error) {
	query := `SELECT id, product_id, original_price, sale_price, discount_percent,
			  total_quantity, sold_quantity, max_per_user, start_time, end_time, status, created_at
			  FROM flash_sales WHERE status = 'ACTIVE' AND NOW() BETWEEN start_time AND end_time`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []*domain.FlashSale
	for rows.Next() {
		var sale domain.FlashSale
		if err := rows.Scan(&sale.ID, &sale.ProductID, &sale.OriginalPrice, &sale.SalePrice,
			&sale.DiscountPercent, &sale.TotalQuantity, &sale.SoldQuantity, &sale.MaxPerUser,
			&sale.StartTime, &sale.EndTime, &sale.Status, &sale.CreatedAt); err != nil {
			return nil, err
		}
		sales = append(sales, &sale)
	}
	return sales, nil
}

func (r *FlashSalePostgresRepository) FindUpcoming(ctx context.Context) ([]*domain.FlashSale, error) {
	query := `SELECT id, product_id, original_price, sale_price, discount_percent,
			  total_quantity, sold_quantity, max_per_user, start_time, end_time, status, created_at
			  FROM flash_sales WHERE status = 'SCHEDULED' AND start_time > NOW()
			  ORDER BY start_time ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []*domain.FlashSale
	for rows.Next() {
		var sale domain.FlashSale
		if err := rows.Scan(&sale.ID, &sale.ProductID, &sale.OriginalPrice, &sale.SalePrice,
			&sale.DiscountPercent, &sale.TotalQuantity, &sale.SoldQuantity, &sale.MaxPerUser,
			&sale.StartTime, &sale.EndTime, &sale.Status, &sale.CreatedAt); err != nil {
			return nil, err
		}
		sales = append(sales, &sale)
	}
	return sales, nil
}

func (r *FlashSalePostgresRepository) Update(ctx context.Context, sale *domain.FlashSale) error {
	query := `UPDATE flash_sales SET sold_quantity = $2, status = $3 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, sale.ID, sale.SoldQuantity, sale.Status)
	return err
}

func (r *FlashSalePostgresRepository) SaveReservation(ctx context.Context, res *domain.FlashSaleReservation) error {
	query := `INSERT INTO flash_sale_reservations 
			  (id, flash_sale_id, user_id, quantity, status, expires_at, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		res.ID, res.FlashSaleID, res.UserID, res.Quantity, res.Status, res.ExpiresAt, res.CreatedAt)
	return err
}

func (r *FlashSalePostgresRepository) SavePurchase(ctx context.Context, purchase *domain.FlashSalePurchase) error {
	query := `INSERT INTO flash_sale_purchases (id, flash_sale_id, user_id, quantity, purchased_at)
			  VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query,
		purchase.ID, purchase.FlashSaleID, purchase.UserID, purchase.Quantity, purchase.PurchasedAt)
	return err
}

func (r *FlashSalePostgresRepository) GetUserPurchases(ctx context.Context, saleID, userID string) (int, error) {
	query := `SELECT COALESCE(SUM(quantity), 0) FROM flash_sale_purchases 
			  WHERE flash_sale_id = $1 AND user_id = $2`

	var total int
	err := r.db.QueryRowContext(ctx, query, saleID, userID).Scan(&total)
	return total, err
}

func (r *FlashSalePostgresRepository) DecrementStock(ctx context.Context, saleID string, quantity int) (bool, error) {
	query := `UPDATE flash_sales SET sold_quantity = sold_quantity + $2 
			  WHERE id = $1 AND sold_quantity + $2 <= total_quantity
			  RETURNING sold_quantity`

	var newSold int
	err := r.db.QueryRowContext(ctx, query, saleID, quantity).Scan(&newSold)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *FlashSalePostgresRepository) Close() error {
	return r.db.Close()
}
