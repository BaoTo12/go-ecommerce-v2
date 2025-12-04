package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/coupon-service/internal/application"
	"github.com/titan-commerce/backend/coupon-service/internal/infrastructure/postgres"
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
		ServiceName: "coupon-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Coupon Service starting...")

	// Initialize repository
	repo := postgres.NewCouponRepository()

	// Initialize application service
	couponService := application.NewCouponService(repo, log)

	// HTTP endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Validate coupon
	http.HandleFunc("/api/v1/coupons/validate", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Code       string   `json:"code"`
			UserID     string   `json:"user_id"`
			OrderValue float64  `json:"order_value"`
			Categories []string `json:"categories"`
			Products   []string `json:"products"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		coupon, discount, err := couponService.ValidateCoupon(r.Context(), req.Code, req.UserID, req.OrderValue, req.Categories, req.Products)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"coupon":   coupon,
			"discount": discount,
		})
	})

	// Apply coupon
	http.HandleFunc("/api/v1/coupons/apply", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Code       string  `json:"code"`
			UserID     string  `json:"user_id"`
			OrderID    string  `json:"order_id"`
			OrderValue float64 `json:"order_value"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		discount, err := couponService.ApplyCoupon(r.Context(), req.Code, req.UserID, req.OrderID, req.OrderValue)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]float64{"discount": discount})
	})

	// List active coupons
	http.HandleFunc("/api/v1/coupons", func(w http.ResponseWriter, r *http.Request) {
		coupons, _ := couponService.GetActiveCoupons(r.Context())
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(coupons)
	})

	go func() {
		addr := fmt.Sprintf(":%d", cfg.HTTPPort)
		log.Infof("Coupon service listening on %s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Coupon Service")
}
