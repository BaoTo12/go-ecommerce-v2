package domain_test

import (
	"testing"
	"time"

	"github.com/titan-commerce/backend/auction-service/internal/domain"
)

func TestNewAuction(t *testing.T) {
	auction := domain.NewAuction("p1", "iPhone 15 Pro", "image.jpg", 30000000, 20000000, 2*time.Hour)

	if auction.ProductID != "p1" {
		t.Errorf("Expected ProductID p1, got %s", auction.ProductID)
	}
	if auction.ProductName != "iPhone 15 Pro" {
		t.Errorf("Expected ProductName iPhone 15 Pro, got %s", auction.ProductName)
	}
	if auction.CurrentBid != 20000000 {
		t.Errorf("Expected CurrentBid 20000000, got %f", auction.CurrentBid)
	}
	if auction.Status != domain.AuctionActive {
		t.Errorf("Expected Status active, got %s", auction.Status)
	}
	if auction.EndTime.Before(time.Now().Add(1*time.Hour)) {
		t.Error("Expected EndTime to be at least 1 hour from now")
	}
}

func TestAuction_PlaceBid(t *testing.T) {
	auction := domain.NewAuction("p1", "Test Product", "img.jpg", 1000000, 500000, 1*time.Hour)

	tests := []struct {
		name      string
		userID    string
		amount    float64
		expectErr error
	}{
		{
			name:      "Valid bid above current + increment",
			userID:    "user1",
			amount:    600100, // 500000 + 100000 (min increment) + 100
			expectErr: nil,
		},
		{
			name:      "Bid too low",
			userID:    "user2",
			amount:    650000, // Less than current (600100) + increment (100000)
			expectErr: domain.ErrBidTooLow,
		},
		{
			name:      "Valid higher bid",
			userID:    "user3",
			amount:    800000,
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := auction.PlaceBid(tt.userID, tt.userID, "avatar.jpg", tt.amount)
			if err != tt.expectErr {
				t.Errorf("PlaceBid() error = %v, want %v", err, tt.expectErr)
			}
		})
	}
}

func TestAuction_PlaceBid_BeforeStart(t *testing.T) {
	auction := &domain.Auction{
		ID:          "a1",
		StartingBid: 100000,
		CurrentBid:  100000,
		MinIncrement: 10000,
		StartTime:   time.Now().Add(1 * time.Hour), // Starts in 1 hour
		EndTime:     time.Now().Add(2 * time.Hour),
		Status:      domain.AuctionPending,
	}

	_, err := auction.PlaceBid("user1", "User", "avatar.jpg", 150000)
	if err != domain.ErrAuctionNotStarted {
		t.Errorf("Expected ErrAuctionNotStarted, got %v", err)
	}
}

func TestAuction_PlaceBid_AfterEnd(t *testing.T) {
	auction := &domain.Auction{
		ID:          "a1",
		StartingBid: 100000,
		CurrentBid:  100000,
		MinIncrement: 10000,
		StartTime:   time.Now().Add(-2 * time.Hour), // Started 2 hours ago
		EndTime:     time.Now().Add(-1 * time.Hour), // Ended 1 hour ago
		Status:      domain.AuctionEnded,
	}

	_, err := auction.PlaceBid("user1", "User", "avatar.jpg", 150000)
	if err != domain.ErrAuctionEnded {
		t.Errorf("Expected ErrAuctionEnded, got %v", err)
	}
}

func TestAuction_UpdateStatus(t *testing.T) {
	tests := []struct {
		name           string
		startTime      time.Time
		endTime        time.Time
		expectedStatus domain.AuctionStatus
	}{
		{
			name:           "Pending - not started",
			startTime:      time.Now().Add(1 * time.Hour),
			endTime:        time.Now().Add(2 * time.Hour),
			expectedStatus: domain.AuctionPending,
		},
		{
			name:           "Active - running",
			startTime:      time.Now().Add(-1 * time.Hour),
			endTime:        time.Now().Add(1 * time.Hour),
			expectedStatus: domain.AuctionActive,
		},
		{
			name:           "Ending - last 10 minutes",
			startTime:      time.Now().Add(-1 * time.Hour),
			endTime:        time.Now().Add(5 * time.Minute),
			expectedStatus: domain.AuctionEnding,
		},
		{
			name:           "Ended - past end time",
			startTime:      time.Now().Add(-2 * time.Hour),
			endTime:        time.Now().Add(-1 * time.Hour),
			expectedStatus: domain.AuctionEnded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auction := &domain.Auction{
				StartTime: tt.startTime,
				EndTime:   tt.endTime,
			}
			auction.UpdateStatus()
			if auction.Status != tt.expectedStatus {
				t.Errorf("UpdateStatus() = %v, want %v", auction.Status, tt.expectedStatus)
			}
		})
	}
}

func TestAuction_TimeRemaining(t *testing.T) {
	auction := &domain.Auction{
		EndTime: time.Now().Add(30 * time.Minute),
	}

	remaining := auction.TimeRemaining()
	if remaining < 29*time.Minute || remaining > 31*time.Minute {
		t.Errorf("TimeRemaining() = %v, expected around 30 minutes", remaining)
	}

	// Test expired auction
	auction.EndTime = time.Now().Add(-10 * time.Minute)
	if auction.TimeRemaining() != 0 {
		t.Error("Expected 0 for expired auction")
	}
}

func TestAuction_ExtendOnLateBid(t *testing.T) {
	auction := &domain.Auction{
		ID:           "a1",
		StartingBid:  100000,
		CurrentBid:   100000,
		MinIncrement: 10000,
		StartTime:    time.Now().Add(-1 * time.Hour),
		EndTime:      time.Now().Add(1 * time.Minute), // 1 minute remaining
		Status:       domain.AuctionEnding,
	}

	originalEnd := auction.EndTime
	_, err := auction.PlaceBid("user1", "User", "avatar.jpg", 150000)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should extend by 2 minutes
	if !auction.EndTime.After(originalEnd) {
		t.Error("Expected auction to be extended")
	}
}

func BenchmarkAuction_PlaceBid(b *testing.B) {
	auction := domain.NewAuction("p1", "Test", "img.jpg", 1000000, 100000, 24*time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		auction.PlaceBid("user", "User", "avatar.jpg", float64(100000+i+1)*1.1)
	}
}
