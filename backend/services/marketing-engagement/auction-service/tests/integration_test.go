package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/titan-commerce/backend/auction-service/internal/application"
	"github.com/titan-commerce/backend/auction-service/internal/domain"
	"github.com/titan-commerce/backend/auction-service/internal/infrastructure/memory"
)

// MockLogger for testing
type MockLogger struct{}

func (m *MockLogger) Info(args ...interface{})                    {}
func (m *MockLogger) Infof(format string, args ...interface{})    {}
func (m *MockLogger) Error(args ...interface{})                   {}
func (m *MockLogger) Errorf(format string, args ...interface{})   {}
func (m *MockLogger) Debug(args ...interface{})                   {}
func (m *MockLogger) Debugf(format string, args ...interface{})   {}
func (m *MockLogger) Warn(args ...interface{})                    {}
func (m *MockLogger) Warnf(format string, args ...interface{})    {}
func (m *MockLogger) Fatal(err error, msg string)                 {}

func setupTestServer() (*httptest.Server, *application.AuctionService) {
	repo := memory.NewAuctionRepository()
	// Note: In real tests, we'd use a proper mock logger
	service := application.NewAuctionService(repo, nil)

	mux := http.NewServeMux()

	// GET /auctions
	mux.HandleFunc("/api/v1/auctions", func(w http.ResponseWriter, r *http.Request) {
		auctions, err := service.GetActiveAuctions(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auctions": auctions,
			"total":    len(auctions),
		})
	})

	// GET /auctions/:id
	mux.HandleFunc("/api/v1/auctions/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/v1/auctions/"):]
		if id == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}

		auction, err := service.GetAuctionByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if auction == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(auction)
	})

	server := httptest.NewServer(mux)
	return server, service
}

func TestGetAuctions_Integration(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/auctions")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result struct {
		Auctions []*domain.Auction `json:"auctions"`
		Total    int               `json:"total"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result.Auctions) == 0 {
		t.Error("Expected at least one auction")
	}
	if result.Total != len(result.Auctions) {
		t.Errorf("Expected total %d to match auction count %d", result.Total, len(result.Auctions))
	}
}

func TestGetAuctionByID_Integration(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	// First get list of auctions to get a valid ID
	resp, _ := http.Get(server.URL + "/api/v1/auctions")
	defer resp.Body.Close()

	var listResult struct {
		Auctions []*domain.Auction `json:"auctions"`
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &listResult)

	if len(listResult.Auctions) == 0 {
		t.Skip("No auctions available for testing")
	}

	auctionID := listResult.Auctions[0].ID

	// Get specific auction
	resp2, err := http.Get(server.URL + "/api/v1/auctions/" + auctionID)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp2.StatusCode)
	}

	var auction domain.Auction
	body2, _ := io.ReadAll(resp2.Body)
	if err := json.Unmarshal(body2, &auction); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if auction.ID != auctionID {
		t.Errorf("Expected auction ID %s, got %s", auctionID, auction.ID)
	}
}

func TestGetAuction_NotFound(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/auctions/nonexistent-id")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestPlaceBid_Integration(t *testing.T) {
	repo := memory.NewAuctionRepository()
	service := application.NewAuctionService(repo, nil)

	// Create a test auction
	auction := domain.NewAuction("p1", "Test Product", "img.jpg", 1000000, 100000, 1*time.Hour)
	repo.SaveAuction(nil, auction)

	// Place a valid bid
	bid, err := service.PlaceBid(nil, auction.ID, "user1", "Test User", "avatar.jpg", 210000)
	if err != nil {
		t.Fatalf("Failed to place bid: %v", err)
	}

	if bid.Amount != 210000 {
		t.Errorf("Expected bid amount 210000, got %f", bid.Amount)
	}

	// Verify auction was updated
	updatedAuction, _ := service.GetAuctionByID(nil, auction.ID)
	if updatedAuction.CurrentBid != 210000 {
		t.Errorf("Expected current bid 210000, got %f", updatedAuction.CurrentBid)
	}
	if updatedAuction.BidCount != 1 {
		t.Errorf("Expected bid count 1, got %d", updatedAuction.BidCount)
	}
}

func TestPlaceBid_TooLow(t *testing.T) {
	repo := memory.NewAuctionRepository()
	service := application.NewAuctionService(repo, nil)

	auction := domain.NewAuction("p1", "Test", "img.jpg", 1000000, 100000, 1*time.Hour)
	repo.SaveAuction(nil, auction)

	// Try to place a bid that's too low
	_, err := service.PlaceBid(nil, auction.ID, "user1", "User", "avatar.jpg", 150000)
	if err != domain.ErrBidTooLow {
		t.Errorf("Expected ErrBidTooLow, got %v", err)
	}
}

// Benchmark: HTTP request throughput
func BenchmarkGetAuctions(b *testing.B) {
	server, _ := setupTestServer()
	defer server.Close()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, _ := client.Get(server.URL + "/api/v1/auctions")
		if resp != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
}

// Test concurrent bid placing
func TestConcurrentBidding(t *testing.T) {
	repo := memory.NewAuctionRepository()
	service := application.NewAuctionService(repo, nil)

	auction := domain.NewAuction("p1", "Test", "img.jpg", 10000000, 1000000, 1*time.Hour)
	repo.SaveAuction(nil, auction)

	// Simulate 10 concurrent bidders
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			bidAmount := float64(1100000 + (idx+1)*100000)
			service.PlaceBid(nil, auction.ID, string(rune('0'+idx)), "User", "avatar.jpg", bidAmount)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify auction is in valid state
	finalAuction, _ := service.GetAuctionByID(nil, auction.ID)
	if finalAuction.CurrentBid < 1100000 {
		t.Error("Expected current bid to be updated after concurrent bidding")
	}
}

// Helper for POST requests
func doPost(url string, body interface{}) (*http.Response, error) {
	jsonBody, _ := json.Marshal(body)
	return http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
}
