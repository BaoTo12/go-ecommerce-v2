package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/titan-commerce/backend/services/marketing-engagement/auction-service/internal/application"
	"github.com/titan-commerce/backend/services/marketing-engagement/auction-service/internal/infrastructure/memory"
	"github.com/titan-commerce/backend/pkg/security"
	"github.com/titan-commerce/backend/pkg/server"
)

/*
Auction Service - Updated to use Unified Server Package

Features automatically enabled:
- JWT Authentication
- CSRF Protection  
- Security Headers (CSP, HSTS, X-Frame-Options)
- Rate Limiting
- Response Compression
- Distributed Tracing
- Metrics Collection
- Graceful Shutdown
*/

func main() {
	// Configuration
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("auction-service-dev-secret")
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8082"
	}

	// Create unified server with all integrations
	config := server.DefaultConfig()
	config.JWTSecret = jwtSecret
	config.Port = port
	config.ServiceName = "auction-service"
	config.CORSOrigins = []string{"http://localhost:3000", "*"}
	config.CSRFEnabled = false // Disable CSRF for API-only service

	srv := server.New(config)

	// Initialize service layer
	repo := memory.NewAuctionRepository()
	auctionService := application.NewAuctionService(repo, nil)

	// Register handlers with automatic middleware (security, perf, tracing, metrics)
	srv.HandleFunc("GET /api/v1/auctions", handleListAuctions(auctionService, srv))
	srv.HandleFunc("GET /api/v1/auctions/{id}", handleGetAuction(auctionService, srv))
	srv.HandleFunc("POST /api/v1/auctions/{id}/bid", handlePlaceBid(auctionService, srv))
	srv.HandleFunc("GET /api/v1/auctions/{id}/bids", handleGetBids(auctionService, srv))

	log.Println("Auction Service starting with unified server...")
	log.Println("Endpoints:")
	log.Println("  GET  /api/v1/auctions              - List active auctions")
	log.Println("  GET  /api/v1/auctions/{id}         - Get auction details")
	log.Println("  POST /api/v1/auctions/{id}/bid     - Place bid")
	log.Println("  GET  /api/v1/auctions/{id}/bids    - Bid history")
	log.Println("  GET  /health                       - Health check (built-in)")
	log.Println("  GET  /metrics                      - Metrics (built-in)")

	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Handlers

func handleListAuctions(svc *application.AuctionService, srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Increment custom metric
		srv.Metrics().Counter("auctions_list_requests").Inc()

		auctions, err := svc.GetActiveAuctions(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auctions": auctions,
			"total":    len(auctions),
		})
	}
}

func handleGetAuction(svc *application.AuctionService, srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auctionID := r.PathValue("id")
		if auctionID == "" {
			http.Error(w, "auction id required", http.StatusBadRequest)
			return
		}

		auction, err := svc.GetAuctionByID(r.Context(), auctionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if auction == nil {
			http.Error(w, "auction not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(auction)
	}
}

func handlePlaceBid(svc *application.AuctionService, srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auctionID := r.PathValue("id")
		if auctionID == "" {
			http.Error(w, "auction id required", http.StatusBadRequest)
			return
		}

		var req struct {
			UserID   string  `json:"user_id"`
			UserName string  `json:"user_name"`
			Avatar   string  `json:"avatar"`
			Amount   float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Input validation using security package
		v := security.NewValidator()
		v.Required("user_id", req.UserID)
		v.Required("user_name", req.UserName).NoXSS("user_name", req.UserName)
		if req.Amount <= 0 {
			v.Result().AddError("amount", "must be positive")
		}

		if !v.Result().Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": v.Result().Errors,
			})
			return
		}

		// Place bid
		bid, err := svc.PlaceBid(r.Context(), auctionID, req.UserID, req.UserName, req.Avatar, req.Amount)
		if err != nil {
			// Increment error metric
			srv.Metrics().Counter("bids_failed").Inc()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Increment success metric
		srv.Metrics().Counter("bids_placed").Inc()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bid)
	}
}

func handleGetBids(svc *application.AuctionService, srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auctionID := r.PathValue("id")
		if auctionID == "" {
			http.Error(w, "auction id required", http.StatusBadRequest)
			return
		}

		bids, err := svc.GetBidHistory(r.Context(), auctionID, 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"bids": bids,
		})
	}
}
