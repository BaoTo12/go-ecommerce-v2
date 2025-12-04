package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/lib/pq"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	"github.com/titan-commerce/backend/product-service/internal/domain"
)

type ProductRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewProductRepository(databaseURL string, logger *logger.Logger) (*ProductRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Product PostgreSQL repository initialized")
	return &ProductRepository{db: db, logger: logger}, nil
}

func (r *ProductRepository) Save(ctx context.Context, product *domain.Product) error {
	imagesJSON, err := json.Marshal(product.Images)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to marshal images", err)
	}

	attributesJSON, err := json.Marshal(product.Attributes)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to marshal attributes", err)
	}

	query := `
		INSERT INTO products (
			id, name, description, price, currency, category_id, 
			images, attributes, stock, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price = EXCLUDED.price,
			currency = EXCLUDED.currency,
			category_id = EXCLUDED.category_id,
			images = EXCLUDED.images,
			attributes = EXCLUDED.attributes,
			stock = EXCLUDED.stock,
			updated_at = EXCLUDED.updated_at
	`

	_, err = r.db.ExecContext(ctx, query,
		product.ID, product.Name, product.Description, product.Price, product.Currency,
		product.CategoryID, imagesJSON, attributesJSON, product.Stock,
		product.CreatedAt, product.UpdatedAt,
	)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save product", err)
	}

	return nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `
		SELECT id, name, description, price, currency, category_id,
			   images, attributes, stock, created_at, updated_at
		FROM products WHERE id = $1
	`

	var p domain.Product
	var imagesJSON, attributesJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Currency, &p.CategoryID,
		&imagesJSON, &attributesJSON, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "product not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find product", err)
	}

	json.Unmarshal(imagesJSON, &p.Images)
	json.Unmarshal(attributesJSON, &p.Attributes)

	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	return r.Save(ctx, product)
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to delete product", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New(errors.ErrNotFound, "product not found")
	}

	return nil
}

func (r *ProductRepository) List(ctx context.Context, page, pageSize int) ([]*domain.Product, int, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT id, name, description, price, currency, category_id,
			   images, attributes, stock, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to list products", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var p domain.Product
		var imagesJSON, attributesJSON []byte
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Price, &p.Currency, &p.CategoryID,
			&imagesJSON, &attributesJSON, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.Wrap(errors.ErrInternal, "failed to scan product", err)
		}
		json.Unmarshal(imagesJSON, &p.Images)
		json.Unmarshal(attributesJSON, &p.Attributes)
		products = append(products, &p)
	}

	var total int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM products").Scan(&total)

	return products, total, nil
}
