package application

import (
	"context"
	"math/rand"
	"time"

	"github.com/titan-commerce/backend/gamification-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type CoinsRepository interface {
	Save(ctx context.Context, coins *domain.UserCoins) error
	FindByUserID(ctx context.Context, userID string) (*domain.UserCoins, error)
	Update(ctx context.Context, coins *domain.UserCoins) error
}

type MissionRepository interface {
	FindActiveByUserID(ctx context.Context, userID string) ([]*domain.Mission, error)
	Update(ctx context.Context, mission *domain.Mission) error
}

type CheckInRepository interface {
	FindByUserID(ctx context.Context, userID string) (*domain.CheckInStreak, error)
	Save(ctx context.Context, streak *domain.CheckInStreak) error
	Update(ctx context.Context, streak *domain.CheckInStreak) error
}

type GamificationService struct {
	coinsRepo   CoinsRepository
	missionRepo MissionRepository
	checkInRepo CheckInRepository
	logger      *logger.Logger
	rand        *rand.Rand
}

func NewGamificationService(coinsRepo CoinsRepository, missionRepo MissionRepository, checkInRepo CheckInRepository, logger *logger.Logger) *GamificationService {
	return &GamificationService{
		coinsRepo:   coinsRepo,
		missionRepo: missionRepo,
		checkInRepo: checkInRepo,
		logger:      logger,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// PlayShakeGame handles the shake-shake game (Command)
func (s *GamificationService) PlayShakeGame(ctx context.Context, userID string, shakeIntensity float32) (int, error) {
	// Validate intensity (prevent cheating)
	if shakeIntensity < 0 || shakeIntensity > 100 {
		return 0, errors.New(errors.ErrInvalidInput, "invalid shake intensity")
	}

	// Calculate coins based on intensity (1-50 coins)
	coinsWon := int(shakeIntensity/2) + s.rand.Intn(10)
	if coinsWon > 50 {
		coinsWon = 50
	}
	if coinsWon < 1 {
		coinsWon = 1
	}

	// Award coins
	coins, err := s.getOrCreateCoins(ctx, userID)
	if err != nil {
		return 0, err
	}

	if err := coins.AddCoins(coinsWon, "shake_game"); err != nil {
		return 0, err
	}

	if err := s.coinsRepo.Update(ctx, coins); err != nil {
		return 0, err
	}

	s.logger.Infof("Shake game played: user=%s, intensity=%.2f, won=%d coins", userID, shakeIntensity, coinsWon)
	return coinsWon, nil
}

// DailyCheckIn handles daily check-in (Command)
func (s *GamificationService) DailyCheckIn(ctx context.Context, userID string) (int, int, error) {
	streak, err := s.checkInRepo.FindByUserID(ctx, userID)
	if err != nil {
		streak = &domain.CheckInStreak{UserID: userID, CurrentStreak: 0}
	}

	coinsWon, err := streak.CheckIn()
	if err != nil {
		return 0, 0, err
	}

	// Update streak
	if err := s.checkInRepo.Update(ctx, streak); err != nil {
		return 0, 0, err
	}

	// Award coins
	coins, err := s.getOrCreateCoins(ctx, userID)
	if err != nil {
		return 0, 0, err
	}

	if err := coins.AddCoins(coinsWon, "daily_checkin"); err != nil {
		return 0, 0, err
	}

	if err := s.coinsRepo.Update(ctx, coins); err != nil {
		return 0, 0, err
	}

	s.logger.Infof("Daily check-in: user=%s, streak=%d, won=%d coins", userID, streak.CurrentStreak, coinsWon)
	return coinsWon, streak.CurrentStreak, nil
}

// PlayLuckyDraw handles lucky draw spin (Command)
func (s *GamificationService) PlayLuckyDraw(ctx context.Context, userID string, costCoins int) (int, string, error) {
	coins, err := s.getOrCreateCoins(ctx, userID)
	if err != nil {
		return 0, "", err
	}

	// Spend coins for lucky draw
	if err := coins.SpendCoins(costCoins); err != nil {
		return 0, "", err
	}

	// Determine prize (weighted random)
	roll := s.rand.Intn(100)
	var coinsWon int
	var prize string

	if roll < 50 { // 50% chance - small win
		coinsWon = 5 + s.rand.Intn(20)
		prize = "Small prize"
	} else if roll < 80 { // 30% chance - medium win
		coinsWon = 50 + s.rand.Intn(100)
		prize = "Medium prize"
	} else if roll < 95 { // 15% chance - big win
		coinsWon = 200 + s.rand.Intn(300)
		prize = "Big prize!"
	} else { // 5% chance - jackpot
		coinsWon = 1000
		prize = "JACKPOT!!!"
	}

	if err := coins.AddCoins(coinsWon, "lucky_draw"); err != nil {
		return 0, "", err
	}

	if err := s.coinsRepo.Update(ctx, coins); err != nil {
		return 0, "", err
	}

	s.logger.Infof("Lucky draw: user=%s, spent=%d, won=%d coins (%s)", userID, costCoins, coinsWon, prize)
	return coinsWon, prize, nil
}

// RedeemCoins redeems coins for discount (Command)
func (s *GamificationService) RedeemCoins(ctx context.Context, userID string, coinsToRedeem int) (string, error) {
	coins, err := s.getOrCreateCoins(ctx, userID)
	if err != nil {
		return "", err
	}

	if err := coins.SpendCoins(coinsToRedeem); err != nil {
		return "", err
	}

	// 100 coins = $1 discount
	discountAmount := float64(coinsToRedeem) / 100.0

	if err := s.coinsRepo.Update(ctx, coins); err != nil {
		return "", err
	}

	// Generate voucher code (simplified)
	voucherCode := "COINS" + time.Now().Format("20060102150405")

	s.logger.Infof("Coins redeemed: user=%s, coins=%d, discount=$%.2f", userID, coinsToRedeem, discountAmount)
	return voucherCode, nil
}

func (s *GamificationService) getOrCreateCoins(ctx context.Context, userID string) (*domain.UserCoins, error) {
	coins, err := s.coinsRepo.FindByUserID(ctx, userID)
	if err != nil {
		coins = domain.NewUserCoins(userID)
		if err :=s.coinsRepo.Save(ctx, coins); err != nil {
			return nil, err
		}
	}
	return coins, nil
}
