package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/titan-commerce/backend/category-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type CategoryRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewCategoryRepository(databaseURL string, logger *logger.Logger) (*CategoryRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Category PostgreSQL repository initialized")
	return &CategoryRepository{db: db, logger: logger}, nil
}

func (r *CategoryRepository) Save(ctx context.Context, category *domain.Category) error {
	query := `
		INSERT INTO categories (id, name, description, parent_id, image_url)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		category.ID, category.Name, category.Description,
		category.ParentID, category.ImageURL,
	)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save category", err)
	}
	return nil
}

func (r *CategoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	query := `SELECT id, name, description, parent_id, image_url FROM categories WHERE id = $1`
	var c domain.Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Description, &c.ParentID, &c.ImageURL,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "category not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find category", err)
	}
	return &c, nil
}

func (r *CategoryRepository) List(ctx context.Context, page, pageSize int) ([]*domain.Category, int, error) {
	offset := (page - 1) * pageSize
	query := `
		SELECT id, name, description, parent_id, image_url
		FROM categories
		ORDER BY name
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to list categories", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.ParentID, &c.ImageURL); err != nil {
			return nil, 0, errors.Wrap(errors.ErrInternal, "failed to scan category", err)
		}
		categories = append(categories, &c)
	}

	var total int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM categories").Scan(&total)

	return categories, total, nil
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]*domain.Category, error) {
	query := `SELECT id, name, description, parent_id, image_url FROM categories ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to get all categories", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.ParentID, &c.ImageURL); err != nil {
			return nil, errors.Wrap(errors.ErrInternal, "failed to scan category", err)
		}
		categories = append(categories, &c)
	}
	return categories, nil
}
