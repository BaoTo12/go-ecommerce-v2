package application

import (
	"context"

	"github.com/titan-commerce/backend/wallet-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type WalletRepository interface {
	Save(ctx context.Context, wallet *domain.Wallet) error
	FindByUserID(ctx context.Context, userID string) (*domain.Wallet, error)
	Update(ctx context.Context, wallet *domain.Wallet) error
}

type TransactionRepository interface {
	Save(ctx context.Context, txn *domain.Transaction) error
	FindByWalletID(ctx context.Context, walletID string, page, pageSize int) ([]*domain.Transaction, int, error)
}

type WalletService struct {
	walletRepo WalletRepository
	txnRepo    TransactionRepository
	logger     *logger.Logger
}

func NewWalletService(walletRepo WalletRepository, txnRepo TransactionRepository, logger *logger.Logger) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
		txnRepo:    txnRepo,
		logger:     logger,
	}
}

// GetBalance retrieves user's wallet (Query)
func (s *WalletService) GetBalance(ctx context.Context, userID string) (*domain.Wallet, error) {
	wallet, err := s.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		wallet = domain.NewWallet(userID, "USD")
		if err := s.walletRepo.Save(ctx, wallet); err != nil {
			return nil, err
		}
	}
	return wallet, nil
}

// Deposit adds funds to wallet (Command)
func (s *WalletService) Deposit(ctx context.Context, userID string, amount float64) (*domain.Wallet, error) {
	wallet, err := s.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := wallet.Deposit(amount); err != nil {
		return nil, err
	}

	if err := s.walletRepo.Update(ctx, wallet); err != nil {
		return nil, err
	}

	// Record transaction
	txn := domain.NewTransaction(wallet.WalletID, "DEPOSIT", amount, "Wallet top-up")
	if err := s.txnRepo.Save(ctx, txn); err != nil {
		s.logger.Error(err, "failed to save transaction")
	}

	s.logger.Infof("Deposit: user=%s, amount=%.2f", userID, amount)
	return wallet, nil
}

// Withdraw removes funds from wallet (Command)
func (s *WalletService) Withdraw(ctx context.Context, userID string, amount float64) (*domain.Wallet, error) {
	wallet, err := s.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := wallet.Withdraw(amount); err != nil {
		return nil, err
	}

	if err := s.walletRepo.Update(ctx, wallet); err != nil {
		return nil, err
	}

	txn := domain.NewTransaction(wallet.WalletID, "WITHDRAWAL", amount, "Withdrawal to bank")
	if err := s.txnRepo.Save(ctx, txn); err != nil {
		s.logger.Error(err, "failed to save transaction")
	}

	s.logger.Infof("Withdraw: user=%s, amount=%.2f", userID, amount)
	return wallet, nil
}

// HoldFunds holds funds in escrow (Command)
func (s *WalletService) HoldFunds(ctx context.Context, userID, orderID string, amount float64) (string, error) {
	wallet, err := s.GetBalance(ctx, userID)
	if err != nil {
		return "", err
	}

	if err := wallet.HoldFunds(amount); err != nil {
		return "", err
	}

	if err := s.walletRepo.Update(ctx, wallet); err != nil {
		return "", err
	}

	holdID := orderID // Use order ID as hold ID for simplicity
	txn := domain.NewTransaction(wallet.WalletID, "HOLD", amount, "Escrow for order: "+orderID)
	if err := s.txnRepo.Save(ctx, txn); err != nil {
		s.logger.Error(err, "failed to save transaction")
	}

	s.logger.Infof("Hold funds: user=%s, amount=%.2f, order=%s", userID, amount, orderID)
	return holdID, nil
}

// ReleaseFunds releases held funds (Command)
func (s *WalletService) ReleaseFunds(ctx context.Context, userID, holdID string, amount float64, refund bool) error {
	wallet, err := s.GetBalance(ctx, userID)
	if err != nil {
		return err
	}

	if err := wallet.ReleaseFunds(amount, refund); err != nil {
		return err
	}

	if err := s.walletRepo.Update(ctx, wallet); err != nil {
		return err
	}

	txnType := "RELEASE"
	if refund {
		txnType = "REFUND"
	}
	txn := domain.NewTransaction(wallet.WalletID, txnType, amount, "Release for hold: "+holdID)
	if err := s.txnRepo.Save(ctx, txn); err != nil {
		s.logger.Error(err, "failed to save transaction")
	}

	s.logger.Infof("Release funds: user=%s, amount=%.2f, refund=%v", userID, amount, refund)
	return nil
}

// GetTransactions retrieves transaction history (Query)
func (s *WalletService) GetTransactions(ctx context.Context, userID string, page, pageSize int) ([]*domain.Transaction, int, error) {
	wallet, err := s.GetBalance(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return s.txnRepo.FindByWalletID(ctx, wallet.WalletID, page, pageSize)
}
