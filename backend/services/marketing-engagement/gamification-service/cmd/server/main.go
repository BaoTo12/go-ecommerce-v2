package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/gamification-service/internal/application"
	"github.com/titan-commerce/backend/gamification-service/internal/infrastructure/postgres"
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
		ServiceName: "gamification-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Gamification Service starting...")

	// Initialize repository
	repo := postgres.NewGamificationRepository()

	// Initialize application service
	gamificationService := application.NewGamificationService(repo, log)

	// HTTP endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Get coin balance
	http.HandleFunc("/api/v1/coins/balance", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		wallet, _ := gamificationService.GetBalance(r.Context(), userID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wallet)
	})

	// Daily check-in
	http.HandleFunc("/api/v1/check-in", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			UserID string `json:"user_id"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		reward, streak, err := gamificationService.DailyCheckIn(r.Context(), req.UserID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"reward": reward,
			"streak": streak,
		})
	})

	// Spin lucky draw
	http.HandleFunc("/api/v1/lucky-draw/spin", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			UserID   string `json:"user_id"`
			SpinCost int    `json:"spin_cost"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		result, err := gamificationService.SpinLuckyDraw(r.Context(), req.UserID, req.SpinCost)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// Get missions
	http.HandleFunc("/api/v1/missions", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		missions, userMissions, _ := gamificationService.GetMissions(r.Context(), userID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"missions":      missions,
			"user_progress": userMissions,
		})
	})

	go func() {
		addr := fmt.Sprintf(":%d", cfg.HTTPPort)
		log.Infof("Gamification service listening on %s", addr)
		log.Info("Shopee Coins, Check-in, Missions, Lucky Draw ready!")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Gamification Service")
}
