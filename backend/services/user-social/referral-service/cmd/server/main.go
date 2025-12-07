package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/referral-service/internal/application"
	"github.com/titan-commerce/backend/referral-service/internal/infrastructure/memory"
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
		ServiceName: "referral-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Referral Service starting...")

	// Initialize in-memory repository
	repo := memory.NewReferralRepository()

	// Initialize application service
	referralService := application.NewReferralService(repo, log)

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

	// POST /api/v1/referrals/generate-code
	http.HandleFunc("/api/v1/referrals/generate-code", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			UserID   string `json:"user_id"`
			UserName string `json:"user_name"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		code, err := referralService.GenerateCode(r.Context(), req.UserID, req.UserName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(code)
	}))

	// GET /api/v1/referrals/stats?user_id=
	http.HandleFunc("/api/v1/referrals/stats", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "user_id required", http.StatusBadRequest)
			return
		}

		stats, err := referralService.GetReferralStats(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}))

	// GET /api/v1/referrals?user_id=
	http.HandleFunc("/api/v1/referrals", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")

		referrals, err := referralService.GetReferrals(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"referrals": referrals,
			"total":     len(referrals),
		})
	}))

	// POST /api/v1/referrals/redeem
	http.HandleFunc("/api/v1/referrals/redeem", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Code     string `json:"code"`
			UserID   string `json:"user_id"`
			UserName string `json:"user_name"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		referral, err := referralService.RedeemCode(r.Context(), req.Code, req.UserID, req.UserName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if referral == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid or expired code"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(referral)
	}))

	// Start HTTP server
	port := cfg.HTTPPort
	if port == 0 {
		port = 8083
	}

	go func() {
		addr := fmt.Sprintf(":%d", port)
		log.Infof("HTTP server listening on %s", addr)
		log.Info("Endpoints:")
		log.Info("  POST /api/v1/referrals/generate-code")
		log.Info("  GET  /api/v1/referrals/stats?user_id=")
		log.Info("  GET  /api/v1/referrals?user_id=")
		log.Info("  POST /api/v1/referrals/redeem")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Referral Service")
}
