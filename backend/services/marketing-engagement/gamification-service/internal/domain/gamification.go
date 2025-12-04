package domain

import (
	"time"

	"github.com/titan-commerce/backend/pkg/errors"
)

type UserCoins struct {
	UserID         string
	Balance        int
	LifetimeEarned int
	LifetimeSpent  int
	UpdatedAt      time.Time
}

func NewUserCoins(userID string) *UserCoins {
	return &UserCoins{
		UserID:         userID,
		Balance:        0,
		LifetimeEarned: 0,
		LifetimeSpent:  0,
		UpdatedAt:      time.Now(),
	}
}

func (c *UserCoins) AddCoins(amount int, reason string) error {
	if amount <= 0 {
		return errors.New(errors.ErrInvalidInput, "amount must be positive")
	}
	c.Balance += amount
	c.LifetimeEarned += amount
	c.UpdatedAt = time.Now()
	return nil
}

func (c *UserCoins) SpendCoins(amount int) error {
	if amount <= 0 {
		return errors.New(errors.ErrInvalidInput, "amount must be positive")
	}
	if c.Balance < amount {
		return errors.New(errors.ErrInsufficientBalance, "insufficient coins")
	}
	c.Balance -= amount
	c.LifetimeSpent += amount
	c.UpdatedAt = time.Now()
	return nil
}

type Mission struct {
	MissionID       string
	UserID          string
	Title           string
	Description     string
	RewardCoins     int
	CurrentProgress int
	TargetProgress  int
	Completed       bool
	CompletedAt     *time.Time
	ExpiresAt       time.Time
}

func (m *Mission) UpdateProgress(progress int) error {
	if m.Completed {
		return errors.New(errors.ErrInvalidInput, "mission already completed")
	}
	
	m.CurrentProgress = progress
	if m.CurrentProgress >= m.TargetProgress {
		m.Completed = true
		now := time.Now()
		m.CompletedAt = &now
	}
	return nil
}

func (m *Mission) IsExpired() bool {
	return time.Now().After(m.ExpiresAt)
}

type CheckInStreak struct {
	UserID          string
	CurrentStreak   int
	LastCheckInDate time.Time
}

func (s *CheckInStreak) CheckIn() (int, error) {
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	lastCheckIn := s.LastCheckInDate.Truncate(24 * time.Hour)

	// Check if already checked in today
	if today.Equal(lastCheckIn) {
		return 0, errors.New(errors.ErrAlreadyExists, "already checked in today")
	}

	// Check if streak continues (checked in yesterday)
	yesterday := today.Add(-24 * time.Hour)
	if lastCheckIn.Equal(yesterday) {
		s.CurrentStreak++
	} else {
		s.CurrentStreak = 1 // Reset streak
	}

	s.LastCheckInDate = now

	// Calculate reward based on streak
	reward := s.calculateReward()
	return reward, nil
}

func (s *CheckInStreak) calculateReward() int {
	// Day 1: 10 coins, Day 7: 100 coins, Day 30: 500 coins
	if s.CurrentStreak >= 30 {
		return 500
	} else if s.CurrentStreak >= 7 {
		return 100
	} else {
		return 10 * s.CurrentStreak
	}
}
