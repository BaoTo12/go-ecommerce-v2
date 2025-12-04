package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/pricing-service/internal/application"
	"github.com/titan-commerce/backend/pricing-service/internal/infrastructure/postgres"
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
		ServiceName: "pricing-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Dynamic Pricing Service starting...")

	// Initialize repository
	repo := postgres.NewPricingRepository()

	// Initialize application service
	pricingService := application.NewPricingService(repo, log)

	// HTTP endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Get product price
	http.HandleFunc("/api/v1/prices", func(w http.ResponseWriter, r *http.Request) {
		productID := r.URL.Query().Get("product_id")
		price, err := pricingService.GetPrice(r.Context(), productID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(price)
	})

	// Set base price
	http.HandleFunc("/api/v1/prices/set", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			ProductID string  `json:"product_id"`
			BasePrice float64 `json:"base_price"`
			MinPrice  float64 `json:"min_price"`
			MaxPrice  float64 `json:"max_price"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		price, err := pricingService.SetBasePrice(r.Context(), req.ProductID, req.BasePrice, req.MinPrice, req.MaxPrice)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(price)
	})

	// Optimize price
	http.HandleFunc("/api/v1/prices/optimize", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		productID := r.URL.Query().Get("product_id")
		price, err := pricingService.OptimizePrice(r.Context(), productID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(price)
	})

	// Get price history
	http.HandleFunc("/api/v1/prices/history", func(w http.ResponseWriter, r *http.Request) {
		productID := r.URL.Query().Get("product_id")
		history, _ := pricingService.GetPriceHistory(r.Context(), productID, 20)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(history)
	})

	go func() {
		addr := fmt.Sprintf(":%d", cfg.HTTPPort)
		log.Infof("Dynamic Pricing service listening on %s", addr)
		log.Info("Demand-based, competitive, and surge pricing ready!")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Dynamic Pricing Service")
}
