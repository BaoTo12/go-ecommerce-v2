package application

import (
	"context"
	"math"
	"time"

	"github.com/titan-commerce/backend/fraud-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
	"github.com/google/uuid"
)

type FraudRepository interface {
	Save(ctx context.Context, check *domain.FraudCheck) error
	FindByID(ctx context.Context, checkID string) (*domain.FraudCheck, error)
	FindByTransaction(ctx context.Context, txnID string) (*domain.FraudCheck, error)
	FindByUser(ctx context.Context, userID string, limit int) ([]*domain.FraudCheck, error)
	GetActiveRules(ctx context.Context) ([]*domain.FraudRule, error)
	SaveAlert(ctx context.Context, alert *domain.FraudAlert) error
	GetUnresolvedAlerts(ctx context.Context) ([]*domain.FraudAlert, error)
	GetUserStats(ctx context.Context, userID string) (*domain.FraudFeatures, error)
	RecordDevice(ctx context.Context, userID, deviceID string) error
	RecordIP(ctx context.Context, userID, ip string) error
	IsNewDevice(ctx context.Context, userID, deviceID string) (bool, error)
	IsNewIP(ctx context.Context, userID, ip string) (bool, error)
}

// FraudScorer simulates an ML model (in production, use ONNX runtime)
type FraudScorer struct {
	weights map[string]float64
}

func NewFraudScorer() *FraudScorer {
	return &FraudScorer{
		weights: map[string]float64{
			"new_device":           0.15,
			"new_ip":               0.10,
			"high_risk_country":    0.20,
			"country_mismatch":     0.15,
			"velocity_score":       0.20,
			"amount_deviation":     0.10,
			"failed_payments":      0.05,
			"account_age":          0.05,
		},
	}
}

func (s *FraudScorer) Score(features *domain.FraudFeatures) float64 {
	score := 0.0

	// New device risk
	if features.NewDevice {
		score += s.weights["new_device"]
	}

	// New IP risk
	if features.NewIP {
		score += s.weights["new_ip"]
	}

	// High risk country
	if features.HighRiskCountry {
		score += s.weights["high_risk_country"]
	}

	// Country mismatch
	if features.CountryMismatch {
		score += s.weights["country_mismatch"]
	}

	// Velocity score (normalized 0-1)
	score += features.VelocityScore * s.weights["velocity_score"]

	// Amount deviation (sigmoid to normalize)
	amountRisk := 1.0 / (1.0 + math.Exp(-features.AmountDeviation+2))
	score += amountRisk * s.weights["amount_deviation"]

	// Failed payments
	if features.FailedPaymentsLast7d > 2 {
		score += s.weights["failed_payments"]
	}

	// Account age (new accounts are riskier)
	if features.AccountAgeDays < 7 {
		score += s.weights["account_age"]
	}

	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}

	return score
}

type FraudService struct {
	repo   FraudRepository
	scorer *FraudScorer
	logger *logger.Logger
}

func NewFraudService(repo FraudRepository, logger *logger.Logger) *FraudService {
	return &FraudService{
		repo:   repo,
		scorer: NewFraudScorer(),
		logger: logger,
	}
}

// CheckTransaction performs real-time fraud check on a transaction
func (s *FraudService) CheckTransaction(ctx context.Context, txnID, userID string, amount float64, currency, ip, deviceID, userAgent string) (*domain.FraudCheck, error) {
	startTime := time.Now()

	check := domain.NewFraudCheck(txnID, userID, amount, currency, ip, deviceID, userAgent)

	// Extract features
	features, err := s.extractFeatures(ctx, check)
	if err != nil {
		s.logger.Error(err, "failed to extract features")
		// Default to review on error
		check.SetScore(0.5)
		check.AddReason("Error extracting features")
	} else {
		check.Features = *features

		// Score using ML model
		score := s.scorer.Score(features)
		check.SetScore(score)

		// Add explainability
		s.addReasons(check, features)
	}

	check.ProcessingTime = time.Since(startTime).Milliseconds()

	// Save check
	if err := s.repo.Save(ctx, check); err != nil {
		return nil, err
	}

	// Record device and IP
	s.repo.RecordDevice(ctx, userID, deviceID)
	s.repo.RecordIP(ctx, userID, ip)

	// Create alert for high risk
	if check.RiskLevel == domain.RiskLevelHigh || check.RiskLevel == domain.RiskLevelCritical {
		s.createAlert(ctx, check)
	}

	s.logger.Infof("Fraud check: txn=%s, score=%.2f, decision=%s, time=%dms",
		txnID, check.Score, check.Decision, check.ProcessingTime)

	return check, nil
}

func (s *FraudService) extractFeatures(ctx context.Context, check *domain.FraudCheck) (*domain.FraudFeatures, error) {
	userStats, err := s.repo.GetUserStats(ctx, check.UserID)
	if err != nil {
		// New user, use defaults
		userStats = &domain.FraudFeatures{
			AccountAgeDays:  0,
			TotalOrders:     0,
			AvgOrderValue:   0,
		}
	}

	// Check if new device/IP
	newDevice, _ := s.repo.IsNewDevice(ctx, check.UserID, check.DeviceID)
	newIP, _ := s.repo.IsNewIP(ctx, check.UserID, check.IP)

	// Calculate velocity score (orders in last 24h / typical)
	velocityScore := 0.0
	if userStats.TotalOrders > 0 {
		avgDailyOrders := float64(userStats.TotalOrders) / float64(max(userStats.AccountAgeDays, 1))
		if avgDailyOrders > 0 {
			velocityScore = float64(userStats.OrdersLast24h) / (avgDailyOrders * 3) // 3x baseline is suspicious
			if velocityScore > 1 {
				velocityScore = 1
			}
		}
	}

	// Calculate amount deviation
	amountDeviation := 0.0
	if userStats.AvgOrderValue > 0 {
		amountDeviation = (check.Amount - userStats.AvgOrderValue) / userStats.AvgOrderValue
	}

	features := &domain.FraudFeatures{
		AccountAgeDays:       userStats.AccountAgeDays,
		TotalOrders:          userStats.TotalOrders,
		AvgOrderValue:        userStats.AvgOrderValue,
		OrdersLast24h:        userStats.OrdersLast24h,
		OrdersLast7d:         userStats.OrdersLast7d,
		FailedPaymentsLast7d: userStats.FailedPaymentsLast7d,
		UniqueDevices:        userStats.UniqueDevices,
		UniqueIPs:            userStats.UniqueIPs,
		NewDevice:            newDevice,
		NewIP:                newIP,
		VelocityScore:        velocityScore,
		AmountDeviation:      amountDeviation,
	}

	return features, nil
}

func (s *FraudService) addReasons(check *domain.FraudCheck, features *domain.FraudFeatures) {
	if features.NewDevice {
		check.AddReason("Transaction from new device")
	}
	if features.NewIP {
		check.AddReason("Transaction from new IP address")
	}
	if features.VelocityScore > 0.5 {
		check.AddReason("Unusual order velocity")
	}
	if features.AmountDeviation > 3 {
		check.AddReason("Order amount significantly higher than usual")
	}
	if features.FailedPaymentsLast7d > 2 {
		check.AddReason("Multiple failed payments recently")
	}
	if features.AccountAgeDays < 7 {
		check.AddReason("New account")
	}
}

func (s *FraudService) createAlert(ctx context.Context, check *domain.FraudCheck) {
	alert := &domain.FraudAlert{
		ID:           uuid.New().String(),
		FraudCheckID: check.ID,
		AlertType:    "high_risk_transaction",
		Severity:     check.RiskLevel,
		Message:      "High risk transaction detected: " + check.TransactionID,
		CreatedAt:    time.Now(),
	}
	s.repo.SaveAlert(ctx, alert)
	s.logger.Warnf("Fraud alert created: %s", alert.ID)
}

// GetFraudCheck retrieves a fraud check by ID
func (s *FraudService) GetFraudCheck(ctx context.Context, checkID string) (*domain.FraudCheck, error) {
	return s.repo.FindByID(ctx, checkID)
}

// GetUserFraudHistory returns fraud check history for a user
func (s *FraudService) GetUserFraudHistory(ctx context.Context, userID string, limit int) ([]*domain.FraudCheck, error) {
	return s.repo.FindByUser(ctx, userID, limit)
}

// GetPendingAlerts returns unresolved fraud alerts
func (s *FraudService) GetPendingAlerts(ctx context.Context) ([]*domain.FraudAlert, error) {
	return s.repo.GetUnresolvedAlerts(ctx)
}

// OverrideDecision allows manual override of fraud decision
func (s *FraudService) OverrideDecision(ctx context.Context, checkID string, decision domain.FraudDecision, reason string) error {
	check, err := s.repo.FindByID(ctx, checkID)
	if err != nil {
		return err
	}

	check.Decision = decision
	check.AddReason("Manual override: " + reason)

	s.logger.Infof("Fraud decision overridden: %s -> %s", checkID, decision)
	return s.repo.Save(ctx, check)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
