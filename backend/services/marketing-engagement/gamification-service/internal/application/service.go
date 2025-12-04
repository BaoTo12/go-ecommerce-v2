package application

import (
	"context"
	"math/rand"
	"time"

	"github.com/titan-commerce/backend/gamification-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	"github.com/google/uuid"
)

type GamificationRepository interface {
	GetWallet(ctx context.Context, userID string) (*domain.CoinWallet, error)
	SaveWallet(ctx context.Context, wallet *domain.CoinWallet) error
	SaveTransaction(ctx context.Context, txn *domain.CoinTransaction) error
	GetTransactions(ctx context.Context, userID string, limit int) ([]*domain.CoinTransaction, error)
	GetCheckIn(ctx context.Context, userID string) (*domain.DailyCheckIn, error)
	SaveCheckIn(ctx context.Context, checkIn *domain.DailyCheckIn) error
	GetActiveMissions(ctx context.Context) ([]*domain.Mission, error)
	GetUserMission(ctx context.Context, userID, missionID string) (*domain.UserMission, error)
	SaveUserMission(ctx context.Context, um *domain.UserMission) error
	GetUserMissions(ctx context.Context, userID string) ([]*domain.UserMission, error)
	GetPrizes(ctx context.Context) ([]*domain.LuckyDrawPrize, error)
	SaveDrawResult(ctx context.Context, result *domain.LuckyDrawResult) error
}

type GamificationService struct {
	repo   GamificationRepository
	logger *logger.Logger
}

func NewGamificationService(repo GamificationRepository, logger *logger.Logger) *GamificationService {
	rand.Seed(time.Now().UnixNano())
	return &GamificationService{
		repo:   repo,
		logger: logger,
	}
}

// GetBalance returns user's coin balance
func (s *GamificationService) GetBalance(ctx context.Context, userID string) (*domain.CoinWallet, error) {
	wallet, err := s.repo.GetWallet(ctx, userID)
	if err != nil {
		// Create new wallet
		wallet = domain.NewCoinWallet(userID)
		if err := s.repo.SaveWallet(ctx, wallet); err != nil {
			return nil, err
		}
	}
	return wallet, nil
}

// EarnCoins adds coins to user's wallet
func (s *GamificationService) EarnCoins(ctx context.Context, userID string, amount int, source, description string) (*domain.CoinWallet, error) {
	wallet, err := s.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	wallet.Earn(amount)
	if err := s.repo.SaveWallet(ctx, wallet); err != nil {
		return nil, err
	}

	txn := domain.NewCoinTransaction(userID, amount, domain.CoinTxnEarn, source, description)
	s.repo.SaveTransaction(ctx, txn)

	s.logger.Infof("User %s earned %d coins from %s", userID, amount, source)
	return wallet, nil
}

// SpendCoins deducts coins from user's wallet
func (s *GamificationService) SpendCoins(ctx context.Context, userID string, amount int, description string) (*domain.CoinWallet, error) {
	wallet, err := s.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !wallet.Spend(amount) {
		return nil, errors.New(errors.ErrInsufficientBalance, "not enough coins")
	}

	if err := s.repo.SaveWallet(ctx, wallet); err != nil {
		return nil, err
	}

	txn := domain.NewCoinTransaction(userID, amount, domain.CoinTxnSpend, "redeem", description)
	s.repo.SaveTransaction(ctx, txn)

	s.logger.Infof("User %s spent %d coins", userID, amount)
	return wallet, nil
}

// DailyCheckIn performs daily check-in and rewards coins
func (s *GamificationService) DailyCheckIn(ctx context.Context, userID string) (int, int, error) {
	checkIn, err := s.repo.GetCheckIn(ctx, userID)
	if err != nil {
		checkIn = &domain.DailyCheckIn{
			UserID:        userID,
			CurrentStreak: 0,
		}
	}

	reward, success := checkIn.CheckIn()
	if !success {
		return 0, checkIn.CurrentStreak, errors.New(errors.ErrInvalidInput, "already checked in today")
	}

	if err := s.repo.SaveCheckIn(ctx, checkIn); err != nil {
		return 0, 0, err
	}

	// Award coins
	if reward > 0 {
		s.EarnCoins(ctx, userID, reward, "check_in", "Daily check-in reward")
	}

	s.logger.Infof("User %s checked in, streak: %d, reward: %d coins", 
		userID, checkIn.CurrentStreak, reward)
	return reward, checkIn.CurrentStreak, nil
}

// GetMissions returns active missions and user progress
func (s *GamificationService) GetMissions(ctx context.Context, userID string) ([]*domain.Mission, []*domain.UserMission, error) {
	missions, err := s.repo.GetActiveMissions(ctx)
	if err != nil {
		return nil, nil, err
	}

	userMissions, err := s.repo.GetUserMissions(ctx, userID)
	if err != nil {
		userMissions = []*domain.UserMission{}
	}

	return missions, userMissions, nil
}

// UpdateMissionProgress updates progress on a mission
func (s *GamificationService) UpdateMissionProgress(ctx context.Context, userID, missionID string, progress int) (*domain.UserMission, error) {
	um, err := s.repo.GetUserMission(ctx, userID, missionID)
	if err != nil {
		um = &domain.UserMission{
			ID:        uuid.New().String(),
			UserID:    userID,
			MissionID: missionID,
			Progress:  0,
		}
	}

	um.Progress += progress
	um.UpdatedAt = time.Now()

	if err := s.repo.SaveUserMission(ctx, um); err != nil {
		return nil, err
	}

	return um, nil
}

// ClaimMissionReward claims reward for completed mission
func (s *GamificationService) ClaimMissionReward(ctx context.Context, userID, missionID string) (int, error) {
	// Get mission details
	missions, err := s.repo.GetActiveMissions(ctx)
	if err != nil {
		return 0, err
	}

	var mission *domain.Mission
	for _, m := range missions {
		if m.ID == missionID {
			mission = m
			break
		}
	}
	if mission == nil {
		return 0, errors.New(errors.ErrNotFound, "mission not found")
	}

	um, err := s.repo.GetUserMission(ctx, userID, missionID)
	if err != nil {
		return 0, err
	}

	if um.Progress < mission.Target {
		return 0, errors.New(errors.ErrInvalidInput, "mission not completed")
	}

	if um.ClaimedAt != nil {
		return 0, errors.New(errors.ErrInvalidInput, "already claimed")
	}

	// Mark claimed
	now := time.Now()
	um.Completed = true
	um.ClaimedAt = &now
	s.repo.SaveUserMission(ctx, um)

	// Award coins
	s.EarnCoins(ctx, userID, mission.Reward, "mission", "Mission reward: "+mission.Name)

	s.logger.Infof("User %s claimed mission %s reward: %d coins", userID, missionID, mission.Reward)
	return mission.Reward, nil
}

// SpinLuckyDraw spins the lucky draw wheel
func (s *GamificationService) SpinLuckyDraw(ctx context.Context, userID string, spinCost int) (*domain.LuckyDrawResult, error) {
	// Deduct spin cost
	if spinCost > 0 {
		_, err := s.SpendCoins(ctx, userID, spinCost, "Lucky draw spin")
		if err != nil {
			return nil, err
		}
	}

	// Get prizes
	prizes, err := s.repo.GetPrizes(ctx)
	if err != nil {
		return nil, err
	}

	// Weighted random selection
	prize := s.selectPrize(prizes)

	result := &domain.LuckyDrawResult{
		ID:      uuid.New().String(),
		UserID:  userID,
		PrizeID: prize.ID,
		Prize:   prize,
		SpunAt:  time.Now(),
	}

	if err := s.repo.SaveDrawResult(ctx, result); err != nil {
		return nil, err
	}

	// Award prize
	if prize.Type == "coins" && prize.Value > 0 {
		s.EarnCoins(ctx, userID, prize.Value, "lucky_draw", "Lucky draw prize")
	}

	s.logger.Infof("User %s won %s from lucky draw", userID, prize.Name)
	return result, nil
}

func (s *GamificationService) selectPrize(prizes []*domain.LuckyDrawPrize) *domain.LuckyDrawPrize {
	if len(prizes) == 0 {
		return &domain.LuckyDrawPrize{ID: "none", Name: "No Prize", Type: "nothing"}
	}

	// Calculate total probability
	var total float64
	for _, p := range prizes {
		total += p.Probability
	}

	// Random selection
	r := rand.Float64() * total
	var cumulative float64
	for _, p := range prizes {
		cumulative += p.Probability
		if r <= cumulative {
			return p
		}
	}

	return prizes[len(prizes)-1]
}

// GetTransactionHistory returns coin transaction history
func (s *GamificationService) GetTransactionHistory(ctx context.Context, userID string, limit int) ([]*domain.CoinTransaction, error) {
	return s.repo.GetTransactions(ctx, userID, limit)
}
