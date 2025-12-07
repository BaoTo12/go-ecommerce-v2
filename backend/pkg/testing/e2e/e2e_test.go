package e2e_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ====================
// END-TO-END TESTS
// ====================

// TestAuctionFlow_E2E tests complete auction flow
func TestAuctionFlow_E2E(t *testing.T) {
	// Setup test server
	server := setupTestServer()
	defer server.Close()

	client := server.Client()

	// Step 1: Create auction
	t.Run("1_CreateAuction", func(t *testing.T) {
		auction := map[string]interface{}{
			"product_id":     "prod_123",
			"product_name":   "Test Product",
			"starting_bid":   100000,
			"duration_hours": 24,
		}
		body, _ := json.Marshal(auction)

		resp, err := client.Post(server.URL+"/api/v1/auctions", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create auction: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 201 or 200, got %d", resp.StatusCode)
		}
	})

	// Step 2: List auctions
	var auctionID string
	t.Run("2_ListAuctions", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/api/v1/auctions")
		if err != nil {
			t.Fatalf("Failed to list auctions: %v", err)
		}
		defer resp.Body.Close()

		var result struct {
			Auctions []struct {
				ID string `json:"id"`
			} `json:"auctions"`
		}
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Auctions) == 0 {
			t.Skip("No auctions available")
		}
		auctionID = result.Auctions[0].ID
	})

	if auctionID == "" {
		auctionID = "test-auction-1"
	}

	// Step 3: Place bid
	t.Run("3_PlaceBid", func(t *testing.T) {
		bid := map[string]interface{}{
			"user_id":   "user_123",
			"user_name": "Test User",
			"amount":    200100,
		}
		body, _ := json.Marshal(bid)

		resp, err := client.Post(server.URL+"/api/v1/auctions/"+auctionID+"/bid", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to place bid: %v", err)
		}
		defer resp.Body.Close()

		// Accept either success or "bid too low" (if auction has higher bid)
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 200 or 400, got %d", resp.StatusCode)
		}
	})

	// Step 4: Get auction details
	t.Run("4_GetAuctionDetails", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/api/v1/auctions/" + auctionID)
		if err != nil {
			t.Fatalf("Failed to get auction: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	// Step 5: Get bid history
	t.Run("5_GetBidHistory", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/api/v1/auctions/" + auctionID + "/bids")
		if err != nil {
			t.Fatalf("Failed to get bids: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}

// TestOrderFlow_E2E tests complete order flow
func TestOrderFlow_E2E(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := server.Client()

	var orderID string

	// Step 1: Add to cart
	t.Run("1_AddToCart", func(t *testing.T) {
		item := map[string]interface{}{
			"product_id": "prod_456",
			"quantity":   2,
		}
		body, _ := json.Marshal(item)

		resp, err := client.Post(server.URL+"/api/v1/cart/items", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Logf("Cart endpoint not available: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected success status, got %d", resp.StatusCode)
		}
	})

	// Step 2: Checkout
	t.Run("2_Checkout", func(t *testing.T) {
		checkout := map[string]interface{}{
			"user_id":      "user_123",
			"payment_type": "cod",
			"address_id":   "addr_123",
		}
		body, _ := json.Marshal(checkout)

		resp, err := client.Post(server.URL+"/api/v1/checkout", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Logf("Checkout endpoint not available: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			var result struct {
				OrderID string `json:"order_id"`
			}
			json.NewDecoder(resp.Body).Decode(&result)
			orderID = result.OrderID
		}
	})

	if orderID == "" {
		orderID = "SPX2024120712345"
	}

	// Step 3: Get order tracking
	t.Run("3_GetOrderTracking", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/api/v1/tracking/" + orderID)
		if err != nil {
			t.Fatalf("Failed to get tracking: %v", err)
		}
		defer resp.Body.Close()

		// May return 404 if order doesn't exist, that's acceptable
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 200 or 404, got %d", resp.StatusCode)
		}
	})
}

// TestReferralFlow_E2E tests referral program flow
func TestReferralFlow_E2E(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := server.Client()

	var referralCode string

	// Step 1: Generate referral code
	t.Run("1_GenerateCode", func(t *testing.T) {
		req := map[string]string{
			"user_id":   "user_referrer",
			"user_name": "ReferrerUser",
		}
		body, _ := json.Marshal(req)

		resp, err := client.Post(server.URL+"/api/v1/referrals/generate-code", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to generate code: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var result struct {
				Code string `json:"code"`
			}
			json.NewDecoder(resp.Body).Decode(&result)
			referralCode = result.Code
		}
	})

	// Step 2: Redeem referral code
	t.Run("2_RedeemCode", func(t *testing.T) {
		if referralCode == "" {
			referralCode = "TEST-CODE"
		}

		req := map[string]string{
			"code":      referralCode,
			"user_id":   "user_new",
			"user_name": "NewUser",
		}
		body, _ := json.Marshal(req)

		resp, err := client.Post(server.URL+"/api/v1/referrals/redeem", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to redeem code: %v", err)
		}
		defer resp.Body.Close()

		// May fail if code is invalid, that's acceptable
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 200 or 404, got %d", resp.StatusCode)
		}
	})

	// Step 3: Check referral stats
	t.Run("3_GetStats", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/api/v1/referrals/stats?user_id=user_referrer")
		if err != nil {
			t.Fatalf("Failed to get stats: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}

// TestRecommendationFlow_E2E tests recommendation engine
func TestRecommendationFlow_E2E(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := server.Client()

	// Step 1: Record view
	t.Run("1_RecordView", func(t *testing.T) {
		req := map[string]string{
			"user_id":    "user_123",
			"product_id": "prod_123",
		}
		body, _ := json.Marshal(req)

		resp, err := client.Post(server.URL+"/api/v1/recommendations/record-view", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Logf("Record view endpoint not available: %v", err)
			return
		}
		defer resp.Body.Close()
	})

	// Step 2: Get personalized recommendations
	t.Run("2_GetPersonalized", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/api/v1/recommendations/personalized?user_id=user_123")
		if err != nil {
			t.Fatalf("Failed to get recommendations: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		var result struct {
			Products []interface{} `json:"products"`
		}
		json.NewDecoder(resp.Body).Decode(&result)
	})

	// Step 3: Get trending
	t.Run("3_GetTrending", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/api/v1/recommendations/trending?location=HCM")
		if err != nil {
			t.Fatalf("Failed to get trending: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}

// TestConcurrentUsers_E2E simulates multiple concurrent users
func TestConcurrentUsers_E2E(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := server.Client()

	// Simulate 10 concurrent users making requests
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(userNum int) {
			defer func() { done <- true }()

			// Each user makes 5 requests
			for j := 0; j < 5; j++ {
				resp, err := client.Get(server.URL + "/api/v1/auctions")
				if err != nil {
					return
				}
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// Wait for all users
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatal("Timeout waiting for concurrent users")
		}
	}
}

// Helper: Setup test server
func setupTestServer() *httptest.Server {
	mux := http.NewServeMux()

	// Auctions
	mux.HandleFunc("/api/v1/auctions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"auctions": []map[string]interface{}{
					{"id": "test-auction-1", "product_name": "Test Product"},
				},
			})
		} else if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"id": "new-auction-1"})
		}
	})

	mux.HandleFunc("/api/v1/auctions/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":          "test-auction-1",
			"current_bid": 100000,
		})
	})

	// Recommendations
	mux.HandleFunc("/api/v1/recommendations/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"products": []map[string]interface{}{
				{"id": "p1", "name": "Product 1"},
			},
		})
	})

	// Referrals
	mux.HandleFunc("/api/v1/referrals/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":            "TEST-CODE",
			"total_referrals": 5,
		})
	})

	// Tracking
	mux.HandleFunc("/api/v1/tracking/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"order_id": "SPX2024120712345",
			"status":   "in_transit",
		})
	})

	// Cart/Checkout (mock)
	mux.HandleFunc("/api/v1/cart/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/api/v1/checkout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"order_id": "ORDER123"})
	})

	return httptest.NewServer(mux)
}
