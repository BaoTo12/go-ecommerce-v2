package domain_test

import (
	"testing"
	"time"

	"github.com/titan-commerce/backend/referral-service/internal/domain"
)

func TestNewReferralCode(t *testing.T) {
	code := domain.NewReferralCode("user123", "NguyenA")

	if code.UserID != "user123" {
		t.Errorf("Expected UserID user123, got %s", code.UserID)
	}
	if code.Uses != 0 {
		t.Errorf("Expected Uses 0, got %d", code.Uses)
	}
	if code.MaxUses != 50 {
		t.Errorf("Expected MaxUses 50, got %d", code.MaxUses)
	}
	if code.RewardPerUse != 500 {
		t.Errorf("Expected RewardPerUse 500, got %d", code.RewardPerUse)
	}
	if code.Code == "" {
		t.Error("Expected non-empty code")
	}
	if code.Link == "" {
		t.Error("Expected non-empty link")
	}
	// Check expiration is ~1 month from now
	expectedExpiry := time.Now().AddDate(0, 1, 0)
	if code.ExpiresAt.Before(expectedExpiry.Add(-24*time.Hour)) || code.ExpiresAt.After(expectedExpiry.Add(24*time.Hour)) {
		t.Errorf("Expected expiry around 1 month from now, got %v", code.ExpiresAt)
	}
}

func TestNewReferral(t *testing.T) {
	referral := domain.NewReferral("referrer1", "referred1", "Nguyen Van B", "CODE123")

	if referral.ReferrerID != "referrer1" {
		t.Errorf("Expected ReferrerID referrer1, got %s", referral.ReferrerID)
	}
	if referral.ReferredID != "referred1" {
		t.Errorf("Expected ReferredID referred1, got %s", referral.ReferredID)
	}
	if referral.Status != domain.ReferralPending {
		t.Errorf("Expected Status pending, got %s", referral.Status)
	}
	if referral.Reward != 0 {
		t.Errorf("Expected Reward 0, got %d", referral.Reward)
	}
	if referral.ID == "" {
		t.Error("Expected non-empty ID")
	}
}

func TestReferral_Complete(t *testing.T) {
	referral := domain.NewReferral("r1", "r2", "User", "CODE")

	if referral.Status != domain.ReferralPending {
		t.Fatal("Expected initial status to be pending")
	}

	referral.Complete(500)

	if referral.Status != domain.ReferralCompleted {
		t.Errorf("Expected Status completed, got %s", referral.Status)
	}
	if referral.Reward != 500 {
		t.Errorf("Expected Reward 500, got %d", referral.Reward)
	}
	if referral.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}
}

func TestReferralCode_Unique(t *testing.T) {
	codes := make(map[string]bool)

	for i := 0; i < 100; i++ {
		code := domain.NewReferralCode("user", "User")
		if codes[code.Code] {
			t.Errorf("Duplicate code generated: %s", code.Code)
		}
		codes[code.Code] = true
	}
}

func TestReferralStatus_Values(t *testing.T) {
	tests := []struct {
		status domain.ReferralStatus
		value  string
	}{
		{domain.ReferralPending, "pending"},
		{domain.ReferralCompleted, "completed"},
		{domain.ReferralExpired, "expired"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.value {
			t.Errorf("Expected %s, got %s", tt.value, tt.status)
		}
	}
}

func BenchmarkNewReferralCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		domain.NewReferralCode("user", "User")
	}
}

func BenchmarkNewReferral(b *testing.B) {
	for i := 0; i < b.N; i++ {
		domain.NewReferral("r1", "r2", "User", "CODE")
	}
}
