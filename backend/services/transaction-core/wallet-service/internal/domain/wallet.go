package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type Wallet struct {
	WalletID         string
	UserID           string
	AvailableBalance float64
	HeldBalance      float64  // Escrow
	Currency         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Version          int  // Optimistic locking
}

func NewWallet(userID, currency string) *Wallet {
	now := time.Now()
	return &Wallet{
		WalletID:         uuid.New().String(),
		UserID:           userID,
		AvailableBalance: 0.0,
		HeldBalance:      0.0,
		Currency:         currency,
		CreatedAt:        now,
		UpdatedAt:        now,
		Version:          1,
	}
}

func (w *Wallet) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New(errors.ErrInvalidInput, "amount must be positive")
	}
	w.AvailableBalance += amount
	w.UpdatedAt = time.Now()
	w.Version++
	return nil
}

func (w *Wallet) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New(errors.ErrInvalidInput, "amount must be positive")
	}
	if w.AvailableBalance < amount {
		return errors.New(errors.ErrInsufficientBalance, "insufficient balance")
	}
	w.AvailableBalance -= amount
	w.UpdatedAt = time.Now()
	w.Version++
	return nil
}

func (w *Wallet) HoldFunds(amount float64) error {
	if amount <= 0 {
		return errors.New(errors.ErrInvalidInput, "amount must be positive")
	}
	if w.AvailableBalance < amount {
		return errors.New(errors.ErrInsufficientBalance, "insufficient balance")
	}
	w.AvailableBalance -= amount
	w.HeldBalance += amount
	w.UpdatedAt = time.Now()
	w.Version++
	return nil
}

func (w *Wallet) ReleaseFunds(amount float64, refund bool) error {
	if amount <= 0 {
		return errors.New(errors.ErrInvalidInput, "amount must be positive")
	}
	if w.HeldBalance < amount {
		return errors.New(errors.ErrInvalidInput, "insufficient held balance")
	}
	
	w.HeldBalance -= amount
	if refund {
		w.AvailableBalance += amount  // Return to user
	}
	// If not refund, funds are released to seller (escrow complete)
	
	w.UpdatedAt = time.Now()
	w.Version++
	return nil
}

func (w *Wallet) GetTotalBalance() float64 {
	return w.AvailableBalance + w.HeldBalance
}

type Transaction struct {
	ID          string
	WalletID    string
	Type        string
	Amount      float64
	Description string
	CreatedAt   time.Time
}

func NewTransaction(walletID, txnType string, amount float64, description string) *Transaction {
	return &Transaction{
		ID:          uuid.New().String(),
		WalletID:    walletID,
		Type:        txnType,
		Amount:      amount,
		Description: description,
		CreatedAt:   time.Now(),
	}
}
