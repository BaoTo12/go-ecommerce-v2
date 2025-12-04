package domain

import (
	"time"

	"github.com/google/uuid"
)

type CoinTransactionType string

const (
	CoinTxnEarn   CoinTransactionType = "EARN"
	CoinTxnSpend  CoinTransactionType = "SPEND"
	CoinTxnExpire CoinTransactionType = "EXPIRE"
)

// CoinWallet represents a user's Shopee Coins balance
type CoinWallet struct {
	UserID    string
	Balance   int
	Lifetime  int // Total coins ever earned
	UpdatedAt time.Time
}

type CoinTransaction struct {
	ID          string
	UserID      string
	Amount      int
	Type        CoinTransactionType
	Source      string // "check_in", "purchase", "game", "mission"
	Description string
	CreatedAt   time.Time
}

// DailyCheckIn tracks user check-in streaks
type DailyCheckIn struct {
	UserID        string
	LastCheckIn   time.Time
	CurrentStreak int
	LongestStreak int
	TotalCheckIns int
}

type Mission struct {
	ID          string
	Name        string
	Description string
	Type        string // "purchase", "review", "share", "invite"
	Target      int    // e.g., buy 3 items
	Reward      int    // coins reward
	StartTime   time.Time
	EndTime     time.Time
	IsActive    bool
}

type UserMission struct {
	ID         string
	UserID     string
	MissionID  string
	Progress   int
	Completed  bool
	ClaimedAt  *time.Time
	UpdatedAt  time.Time
}

type LuckyDrawPrize struct {
	ID          string
	Name        string
	Type        string // "coins", "voucher", "product", "nothing"
	Value       int
	Probability float64 // 0-100
	ImageURL    string
}

type LuckyDrawResult struct {
	ID         string
	UserID     string
	PrizeID    string
	Prize      *LuckyDrawPrize
	SpunAt     time.Time
}

func NewCoinWallet(userID string) *CoinWallet {
	return &CoinWallet{
		UserID:    userID,
		Balance:   0,
		Lifetime:  0,
		UpdatedAt: time.Now(),
	}
}

func (w *CoinWallet) Earn(amount int) {
	w.Balance += amount
	w.Lifetime += amount
	w.UpdatedAt = time.Now()
}

func (w *CoinWallet) Spend(amount int) bool {
	if w.Balance < amount {
		return false
	}
	w.Balance -= amount
	w.UpdatedAt = time.Now()
	return true
}

func NewCoinTransaction(userID string, amount int, txnType CoinTransactionType, source, desc string) *CoinTransaction {
	return &CoinTransaction{
		ID:          uuid.New().String(),
		UserID:      userID,
		Amount:      amount,
		Type:        txnType,
		Source:      source,
		Description: desc,
		CreatedAt:   time.Now(),
	}
}

func (d *DailyCheckIn) CheckIn() (int, bool) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastDay := time.Date(d.LastCheckIn.Year(), d.LastCheckIn.Month(), d.LastCheckIn.Day(), 0, 0, 0, 0, d.LastCheckIn.Location())

	// Already checked in today
	if today.Equal(lastDay) {
		return 0, false
	}

	d.LastCheckIn = now
	d.TotalCheckIns++

	// Check if streak continues
	yesterDay := today.AddDate(0, 0, -1)
	if lastDay.Equal(yesterDay) {
		d.CurrentStreak++
	} else {
		d.CurrentStreak = 1
	}

	if d.CurrentStreak > d.LongestStreak {
		d.LongestStreak = d.CurrentStreak
	}

	// Calculate reward based on streak (Day 1: 1 coin, Day 7: 10 coins)
	reward := d.CurrentStreak
	if reward > 10 {
		reward = 10
	}

	return reward, true
}

type Repository interface {
	GetWallet(ctx interface{}, userID string) (*CoinWallet, error)
	SaveWallet(ctx interface{}, wallet *CoinWallet) error
	SaveTransaction(ctx interface{}, txn *CoinTransaction) error
	GetTransactions(ctx interface{}, userID string, limit int) ([]*CoinTransaction, error)
	GetCheckIn(ctx interface{}, userID string) (*DailyCheckIn, error)
	SaveCheckIn(ctx interface{}, checkIn *DailyCheckIn) error
	GetActiveMissions(ctx interface{}) ([]*Mission, error)
	GetUserMission(ctx interface{}, userID, missionID string) (*UserMission, error)
	SaveUserMission(ctx interface{}, um *UserMission) error
	GetPrizes(ctx interface{}) ([]*LuckyDrawPrize, error)
	SaveDrawResult(ctx interface{}, result *LuckyDrawResult) error
}
