package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/titan-commerce/backend/auction-service/internal/application"
	"github.com/titan-commerce/backend/auction-service/internal/infrastructure/memory"
	"github.com/titan-commerce/backend/pkg/config"
	"github.com/titan-commerce/backend/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(logger.Config{
		Level:       cfg.LogLevel,
		ServiceName: "auction-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Auction Service starting...")

	// Initialize in-memory repository with sample auctions
	repo := memory.NewAuctionRepository()

	// Initialize application service
	auctionService := application.NewAuctionService(repo, log)

	// CORS middleware
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next(w, r)
		}
	}

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// GET /api/v1/auctions - List active auctions
	http.HandleFunc("/api/v1/auctions", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			// Check if it's a specific auction request
			if strings.HasPrefix(r.URL.Path, "/api/v1/auctions/") {
				return
			}
		}

		auctions, err := auctionService.GetActiveAuctions(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auctions": auctions,
			"total":    len(auctions),
		})
	}))

	// GET /api/v1/auctions/:id - Get auction details
	// POST /api/v1/auctions/:id/bid - Place bid
	http.HandleFunc("/api/v1/auctions/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/auctions/")
		parts := strings.Split(path, "/")

		if len(parts) == 0 || parts[0] == "" {
			http.Error(w, "auction id required", http.StatusBadRequest)
			return
		}

		auctionID := parts[0]

		// Handle bid endpoint
		if len(parts) == 2 && parts[1] == "bid" {
			if r.Method != "POST" {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var req struct {
				UserID   string  `json:"user_id"`
				UserName string  `json:"user_name"`
				Avatar   string  `json:"avatar"`
				Amount   float64 `json:"amount"`
			}
			json.NewDecoder(r.Body).Decode(&req)

			bid, err := auctionService.PlaceBid(r.Context(), auctionID, req.UserID, req.UserName, req.Avatar, req.Amount)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(bid)
			return
		}

		// Handle bid history endpoint
		if len(parts) == 2 && parts[1] == "bids" {
			bids, err := auctionService.GetBidHistory(r.Context(), auctionID, 20)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"bids": bids,
			})
			return
		}

		// Get auction details
		auction, err := auctionService.GetAuctionByID(r.Context(), auctionID)
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
	}))

	// Start HTTP server
	port := cfg.HTTPPort
	if port == 0 {
		port = 8082
	}

	go func() {
		addr := fmt.Sprintf(":%d", port)
		log.Infof("HTTP server listening on %s", addr)
		log.Info("Endpoints:")
		log.Info("  GET  /api/v1/auctions              - List active auctions")
		log.Info("  GET  /api/v1/auctions/:id          - Get auction details")
		log.Info("  POST /api/v1/auctions/:id/bid      - Place bid")
		log.Info("  GET  /api/v1/auctions/:id/bids     - Bid history")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Auction Service")
}
