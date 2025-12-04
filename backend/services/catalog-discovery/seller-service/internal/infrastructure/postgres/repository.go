package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/lib/pq"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	"github.com/titan-commerce/backend/seller-service/internal/domain"
)

type SellerRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewSellerRepository(databaseURL string, logger *logger.Logger) (*SellerRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Seller PostgreSQL repository initialized")
	return &SellerRepository{db: db, logger: logger}, nil
}

func (r *SellerRepository) Save(ctx context.Context, seller *domain.Seller) error {
	documentsJSON, _ := json.Marshal(seller.Documents)
	addressJSON, _ := json.Marshal(seller.BusinessAddress)

	query := `
		INSERT INTO sellers (seller_id, user_id, business_name, business_type, registration_number, tax_id, 
			business_address, documents, status, rating, total_sales, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.ExecContext(ctx, query,
		seller.SellerID, seller.UserID, seller.BusinessName, seller.BusinessType, seller.RegistrationNumber,
		seller.TaxID, addressJSON, documentsJSON, seller.Status, seller.Rating, seller.TotalSales,
		seller.CreatedAt, seller.UpdatedAt,
	)
	return err
}

func (r *SellerRepository) FindByID(ctx context.Context, sellerID string) (*domain.Seller, error) {
	query := `SELECT seller_id, user_id, business_name, business_type, registration_number, tax_id, 
		business_address, documents, status, rating, total_sales, created_at, updated_at 
		FROM sellers WHERE seller_id = $1`

	var seller domain.Seller
	var documentsJSON, addressJSON []byte

	err := r.db.QueryRowContext(ctx, query, sellerID).Scan(
		&seller.SellerID, &seller.UserID, &seller.BusinessName, &seller.BusinessType,
		&seller.RegistrationNumber, &seller.TaxID, &addressJSON, &documentsJSON,
		&seller.Status, &seller.Rating, &seller.TotalSales, &seller.CreatedAt, &seller.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "seller not found")
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(documentsJSON, &seller.Documents)
	json.Unmarshal(addressJSON, &seller.BusinessAddress)
	return &seller, nil
}

func (r *SellerRepository) FindByUserID(ctx context.Context, userID string) (*domain.Seller, error) {
	query := `SELECT seller_id, user_id, business_name, business_type, registration_number, tax_id, 
		business_address, documents, status, rating, total_sales, created_at, updated_at 
		FROM sellers WHERE user_id = $1`

	var seller domain.Seller
	var documentsJSON, addressJSON []byte

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&seller.SellerID, &seller.UserID, &seller.BusinessName, &seller.BusinessType,
		&seller.RegistrationNumber, &seller.TaxID, &addressJSON, &documentsJSON,
		&seller.Status, &seller.Rating, &seller.TotalSales, &seller.CreatedAt, &seller.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "seller not found")
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(documentsJSON, &seller.Documents)
	json.Unmarshal(addressJSON, &seller.BusinessAddress)
	return &seller, nil
}

func (r *SellerRepository) Update(ctx context.Context, seller *domain.Seller) error {
	documentsJSON, _ := json.Marshal(seller.Documents)
	addressJSON, _ := json.Marshal(seller.BusinessAddress)

	query := `
		UPDATE sellers 
		SET business_name = $2, business_type = $3, registration_number = $4, tax_id = $5,
			business_address = $6, documents = $7, status = $8, rating = $9, total_sales = $10, updated_at = $11
		WHERE seller_id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		seller.SellerID, seller.BusinessName, seller.BusinessType, seller.RegistrationNumber,
		seller.TaxID, addressJSON, documentsJSON, seller.Status, seller.Rating,
		seller.TotalSales, seller.UpdatedAt,
	)
	return err
}

func (r *SellerRepository) List(ctx context.Context, status domain.SellerStatus, page, pageSize int) ([]*domain.Seller, int, error) {
	offset := (page - 1) * pageSize

	query := `SELECT seller_id, user_id, business_name, business_type, registration_number, tax_id, 
		business_address, documents, status, rating, total_sales, created_at, updated_at 
		FROM sellers WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, status, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sellers []*domain.Seller
	for rows.Next() {
		var seller domain.Seller
		var documentsJSON, addressJSON []byte

		if err := rows.Scan(&seller.SellerID, &seller.UserID, &seller.BusinessName, &seller.BusinessType,
			&seller.RegistrationNumber, &seller.TaxID, &addressJSON, &documentsJSON,
			&seller.Status, &seller.Rating, &seller.TotalSales, &seller.CreatedAt, &seller.UpdatedAt); err != nil {
			return nil, 0, err
		}

		json.Unmarshal(documentsJSON, &seller.Documents)
		json.Unmarshal(addressJSON, &seller.BusinessAddress)
		sellers = append(sellers, &seller)
	}

	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sellers WHERE status = $1`, status).Scan(&total)

	return sellers, total, nil
}

func (r *SellerRepository) Close() error {
	return r.db.Close()
}
