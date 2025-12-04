package domain

import (
	"time"

	"github.com/google/uuid"
)

type FraudDecision string
type RiskLevel string

const (
	FraudDecisionAllow   FraudDecision = "ALLOW"
	FraudDecisionReview  FraudDecision = "REVIEW"
	FraudDecisionBlock   FraudDecision = "BLOCK"

	RiskLevelLow      RiskLevel = "LOW"
	RiskLevelMedium   RiskLevel = "MEDIUM"
	RiskLevelHigh     RiskLevel = "HIGH"
	RiskLevelCritical RiskLevel = "CRITICAL"
)

type FraudCheck struct {
	ID              string
	TransactionID   string
	UserID          string
	Amount          float64
	Currency        string
	IP              string
	DeviceID        string
	UserAgent       string
	Features        FraudFeatures
	Score           float64 // 0-1, higher = more fraudulent
	RiskLevel       RiskLevel
	Decision        FraudDecision
	Reasons         []string
	ProcessingTime  int64 // milliseconds
	CreatedAt       time.Time
}

type FraudFeatures struct {
	AccountAgeDays       int
	TotalOrders          int
	AvgOrderValue        float64
	OrdersLast24h        int
	OrdersLast7d         int
	FailedPaymentsLast7d int
	UniqueDevices        int
	UniqueIPs            int
	CountryMismatch      bool
	HighRiskCountry      bool
	NewDevice            bool
	NewIP                bool
	VelocityScore        float64
	AmountDeviation      float64 // std deviations from user's avg
}

type FraudRule struct {
	ID          string
	Name        string
	Description string
	Condition   string // e.g., "amount > 10000 && new_device"
	Action      FraudDecision
	ScoreWeight float64
	IsActive    bool
	Priority    int
	CreatedAt   time.Time
}

type FraudAlert struct {
	ID           string
	FraudCheckID string
	AlertType    string
	Severity     RiskLevel
	Message      string
	Acknowledged bool
	ResolvedAt   *time.Time
	CreatedAt    time.Time
}

func NewFraudCheck(txnID, userID string, amount float64, currency, ip, deviceID, userAgent string) *FraudCheck {
	return &FraudCheck{
		ID:            uuid.New().String(),
		TransactionID: txnID,
		UserID:        userID,
		Amount:        amount,
		Currency:      currency,
		IP:            ip,
		DeviceID:      deviceID,
		UserAgent:     userAgent,
		Reasons:       []string{},
		CreatedAt:     time.Now(),
	}
}

func (fc *FraudCheck) SetScore(score float64) {
	fc.Score = score
	
	switch {
	case score < 0.3:
		fc.RiskLevel = RiskLevelLow
		fc.Decision = FraudDecisionAllow
	case score < 0.7:
		fc.RiskLevel = RiskLevelMedium
		fc.Decision = FraudDecisionReview
	case score < 0.9:
		fc.RiskLevel = RiskLevelHigh
		fc.Decision = FraudDecisionReview
	default:
		fc.RiskLevel = RiskLevelCritical
		fc.Decision = FraudDecisionBlock
	}
}

func (fc *FraudCheck) AddReason(reason string) {
	fc.Reasons = append(fc.Reasons, reason)
}

type Repository interface {
	Save(ctx interface{}, check *FraudCheck) error
	FindByID(ctx interface{}, checkID string) (*FraudCheck, error)
	FindByTransaction(ctx interface{}, txnID string) (*FraudCheck, error)
	FindByUser(ctx interface{}, userID string, limit int) ([]*FraudCheck, error)
	GetActiveRules(ctx interface{}) ([]*FraudRule, error)
	SaveAlert(ctx interface{}, alert *FraudAlert) error
	GetUnresolvedAlerts(ctx interface{}) ([]*FraudAlert, error)
	GetUserStats(ctx interface{}, userID string) (*FraudFeatures, error)
}
