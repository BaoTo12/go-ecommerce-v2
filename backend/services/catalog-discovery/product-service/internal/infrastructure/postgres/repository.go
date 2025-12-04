package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/titan-commerce/backend/product-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
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
	variantsJSON, err := json.Marshal(product.Variants)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to marshal variants", err)
	}

	imageURLsJSON, err := json.Marshal(product.ImageURLs)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to marshal image URLs", err)
	}

	query := `
		INSERT INTO products (
			id, seller_id, name, description, category_id, 
			variants, image_urls, status, rating, review_count, sold_count,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			category_id = EXCLUDED.category_id,
			variants = EXCLUDED.variants,
			image_urls = EXCLUDED.image_urls,
			status = EXCLUDED.status,
			rating = EXCLUDED.rating,
			review_count = EXCLUDED.review_count,
			sold_count = EXCLUDED.sold_count,
			updated_at = EXCLUDED.updated_at
	`

	_, err = r.db.ExecContext(ctx, query,
		product.ID, product.SellerID, product.Name, product.Description, product.CategoryID,
		variantsJSON, imageURLsJSON, product.Status, product.Rating, product.ReviewCount,
		product.SoldCount, product.CreatedAt, product.UpdatedAt,
	)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save product", err)
	}

	return nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `
		SELECT id, seller_id, name, description, category_id,
			   variants, image_urls, status, rating, review_count, sold_count,
			   created_at, updated_at
		FROM products WHERE id = $1
	`

	var p domain.Product
	var variantsJSON, imageURLsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SellerID, &p.Name, &p.Description, &p.CategoryID,
		&variantsJSON, &imageURLsJSON, &p.Status, &p.Rating,
		&p.ReviewCount, &p.SoldCount, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "product not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find product", err)
	}

	if err := json.Unmarshal(variantsJSON, &p.Variants); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to unmarshal variants", err)
	}
	if err := json.Unmarshal(imageURLsJSON, &p.ImageURLs); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to unmarshal image URLs", err)
	}

	return &p, nil
}

func (r *ProductRepository) List(ctx context.Context, page, pageSize int) ([]*domain.Product, int, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT id, seller_id, name, description, category_id,
			   variants, image_urls, status, rating, review_count, sold_count,
			   created_at, updated_at
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
		var variantsJSON, imageURLsJSON []byte
		err := rows.Scan(
			&p.ID, &p.SellerID, &p.Name, &p.Description, &p.CategoryID,
			&variantsJSON, &imageURLsJSON, &p.Status, &p.Rating,
			&p.ReviewCount, &p.SoldCount, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.Wrap(errors.ErrInternal, "failed to scan product", err)
		}
		json.Unmarshal(variantsJSON, &p.Variants)
		json.Unmarshal(imageURLsJSON, &p.ImageURLs)
		products = append(products, &p)
	}

	var total int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM products").Scan(&total)

	return products, total, nil
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
