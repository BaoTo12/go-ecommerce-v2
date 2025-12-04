package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/titan-commerce/backend/auth-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type AuthRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewAuthRepository(databaseURL string, logger *logger.Logger) (*AuthRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Auth PostgreSQL repository initialized")
	return &AuthRepository{db: db, logger: logger}, nil
}

func (r *AuthRepository) SaveUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO auth_users (
			id, email, password_hash, full_name, 
			mfa_enabled, mfa_secret, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FullName,
		user.MFAEnabled, user.MFASecret, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save user", err)
	}

	return nil
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, 
			   mfa_enabled, mfa_secret, created_at, updated_at
		FROM auth_users WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.MFAEnabled, &user.MFASecret, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "user not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find user", err)
	}

	return &user, nil
}

func (r *AuthRepository) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, 
			   mfa_enabled, mfa_secret, created_at, updated_at
		FROM auth_users WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.MFAEnabled, &user.MFASecret, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "user not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find user", err)
	}

	return &user, nil
}

func (r *AuthRepository) UpdateMFA(ctx context.Context, userID, secret string, enabled bool) error {
	query := `UPDATE auth_users SET mfa_secret = $1, mfa_enabled = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, secret, enabled, userID)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to update MFA", err)
	}
	return nil
}
