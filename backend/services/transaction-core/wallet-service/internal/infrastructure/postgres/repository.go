package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/titan-commerce/backend/wallet-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type WalletRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewWalletRepository(databaseURL string, logger *logger.Logger) (*WalletRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Wallet PostgreSQL repository initialized")
	return &WalletRepository{db: db, logger: logger}, nil
}

func (r *WalletRepository) Save(ctx context.Context, wallet *domain.Wallet) error {
	query := `
		INSERT INTO wallets (
			wallet_id, user_id, available_balance, held_balance, currency,
			created_at, updated_at, version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		wallet.WalletID, wallet.UserID, wallet.AvailableBalance, wallet.HeldBalance,
		wallet.Currency, wallet.CreatedAt, wallet.UpdatedAt, wallet.Version,
	)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save wallet", err)
	}

	return nil
}

func (r *WalletRepository) FindByUserID(ctx context.Context, userID string) (*domain.Wallet, error) {
	query := `
		SELECT wallet_id, user_id, available_balance, held_balance, currency,
			   created_at, updated_at, version
		FROM wallets WHERE user_id = $1
	`

	var wallet domain.Wallet
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&wallet.WalletID, &wallet.UserID, &wallet.AvailableBalance, &wallet.HeldBalance,
		&wallet.Currency, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.Version,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "wallet not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find wallet", err)
	}

	return &wallet, nil
}

func (r *WalletRepository) Update(ctx context.Context, wallet *domain.Wallet) error {
	query := `
		UPDATE wallets
		SET available_balance = $1, held_balance = $2, updated_at = $3, version = $4
		WHERE wallet_id = $5 AND version = $6
	`

	result, err := r.db.ExecContext(ctx, query,
		wallet.AvailableBalance, wallet.HeldBalance, wallet.UpdatedAt,
		wallet.Version, wallet.WalletID, wallet.Version-1,
	)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to update wallet", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to get rows affected", err)
	}

	if rows == 0 {
		return errors.New(errors.ErrConflict, "wallet was modified by another transaction (optimistic lock)")
	}

	return nil
}

type TransactionRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewTransactionRepository(databaseURL string, logger *logger.Logger) (*TransactionRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	logger.Info("Transaction PostgreSQL repository initialized")
	return &TransactionRepository{db: db, logger: logger}, nil
}

func (r *TransactionRepository) Save(ctx context.Context, txn *domain.Transaction) error {
	query := `
		INSERT INTO wallet_transactions (id, wallet_id, type, amount, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		txn.ID, txn.WalletID, txn.Type, txn.Amount, txn.Description, txn.CreatedAt,
	)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save transaction", err)
	}

	return nil
}

func (r *TransactionRepository) FindByWalletID(ctx context.Context, walletID string, page, pageSize int) ([]*domain.Transaction, int, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT id, wallet_id, type, amount, description, created_at
		FROM wallet_transactions
		WHERE wallet_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, walletID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to find transactions", err)
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		var txn domain.Transaction
		err := rows.Scan(&txn.ID, &txn.WalletID, &txn.Type, &txn.Amount, &txn.Description, &txn.CreatedAt)
		if err != nil {
			return nil, 0, errors.Wrap(errors.ErrInternal, "failed to scan transaction", err)
		}
		transactions = append(transactions, &txn)
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM wallet_transactions WHERE wallet_id = $1`
	err = r.db.QueryRowContext(ctx, countQuery, walletID).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to count transactions", err)
	}

	return transactions, total, nil
}
