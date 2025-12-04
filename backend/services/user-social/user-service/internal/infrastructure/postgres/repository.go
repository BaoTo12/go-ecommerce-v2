package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/titan-commerce/backend/user-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewUserRepository(databaseURL string, logger *logger.Logger) (*UserRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("User PostgreSQL repository initialized")
	return &UserRepository{db: db, logger: logger}, nil
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	addressesJSON, err := json.Marshal(user.Addresses)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to marshal addresses", err)
	}

	preferencesJSON, err := json.Marshal(user.Preferences)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to marshal preferences", err)
	}

	query := `
		INSERT INTO users (
			id, email, full_name, phone_number, avatar_url, 
			addresses, preferences, created_at, updated_at, version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			full_name = EXCLUDED.full_name,
			phone_number = EXCLUDED.phone_number,
			avatar_url = EXCLUDED.avatar_url,
			addresses = EXCLUDED.addresses,
			preferences = EXCLUDED.preferences,
			updated_at = EXCLUDED.updated_at,
			version = EXCLUDED.version
	`

	_, err = r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.FullName, user.PhoneNumber, user.AvatarURL,
		addressesJSON, preferencesJSON, user.CreatedAt, user.UpdatedAt, user.Version,
	)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save user", err)
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	return r.Save(ctx, user)
}

func (r *UserRepository) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `
		SELECT id, email, full_name, phone_number, avatar_url, 
			   addresses, preferences, created_at, updated_at, version
		FROM users WHERE id = $1
	`

	var user domain.User
	var addressesJSON, preferencesJSON []byte

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.FullName, &user.PhoneNumber, &user.AvatarURL,
		&addressesJSON, &preferencesJSON, &user.CreatedAt, &user.UpdatedAt, &user.Version,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "user not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find user", err)
	}

	if err := json.Unmarshal(addressesJSON, &user.Addresses); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to unmarshal addresses", err)
	}
	if err := json.Unmarshal(preferencesJSON, &user.Preferences); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to unmarshal preferences", err)
	}

	return &user, nil
}
