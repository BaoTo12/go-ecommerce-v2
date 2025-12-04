package domain

import (
	"time"

	"github.com/google/uuid"
)

// RiskLevel represents the fraud risk level
type RiskLevel string

const (
	RiskLow      RiskLevel = "LOW"
	RiskMedium   RiskLevel = "MEDIUM"
	RiskHigh     RiskLevel = "HIGH"
	RiskCritical RiskLevel = "CRITICAL"
)

// FraudCheck represents a fraud detection analysis
type FraudCheck struct {
	CheckID          string
	OrderID          string
	UserID           string
	RiskScore        float64 // 0.0 - 1.0
	RiskLevel        RiskLevel
	Flags            []FraudFlag
	Recommendation   string // APPROVE, REVIEW, REJECT
	Reasons          []string
	DeviceFingerprint string
	IPAddress        string
	Location         string
	CreatedAt        time.Time
}

// FraudFlag represents a specific fraud indicator
type FraudFlag struct {
	Type     string  // velocity, location, card, behavior
	Severity string  // low, medium, high
	Score    float64
	Message  string
}

// NewFraudCheck creates a new fraud check
func NewFraudCheck(orderID, userID, ipAddress string) *FraudCheck {
	return &FraudCheck{
		CheckID:   uuid.New().String(),
		OrderID:   orderID,
		UserID:    userID,
		IPAddress: ipAddress,
		Flags:     []FraudFlag{},
		Reasons:   []string{},
		CreatedAt: time.Now(),
	}
}

// CalculateRisk calculates overall risk score and level
func (f *FraudCheck) CalculateRisk() {
	var totalScore float64
	for _, flag := range f.Flags {
		totalScore += flag.Score
	}

	f.RiskScore = totalScore / float64(len(f.Flags))

	if f.RiskScore < 0.3 {
		f.RiskLevel = RiskLow
		f.Recommendation = "APPROVE"
	} else if f.RiskScore < 0.6 {
		f.RiskLevel = RiskMedium
		f.Recommendation = "REVIEW"
	} else if f.RiskScore < 0.8 {
		f.RiskLevel = RiskHigh
		f.Recommendation = "REVIEW"
	} else {
		f.RiskLevel = RiskCritical
		f.Recommendation = "REJECT"
	}
}

// UserRiskProfile represents a user's fraud risk profile
type UserRiskProfile struct {
	UserID              string
	TrustScore          float64
	TotalOrders         int
	FailedOrders        int
	ChargebackCount     int
	AccountAge          int // days
	VerificationStatus  string
	LastFraudCheck      *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

