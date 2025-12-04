package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/gamification-service/internal/application"
	grpcServer "github.com/titan-commerce/backend/gamification-service/internal/interface/grpc"
	"github.com/titan-commerce/backend/gamification-service/internal/infrastructure/postgres"
	"github.com/titan-commerce/backend/pkg/config"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	// Initialize PostgreSQL repository
	repo, err := postgres.NewGamificationRepository(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize repository")
	}
	defer repo.Close()

	// Initialize application service
	gamificationService := application.NewGamificationService(repo, log)

	// Start gRPC server
	go func() {
		grpcAddr := fmt.Sprintf(":%d", cfg.GRPCPort)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatal(err, "Failed to listen for gRPC")
		}

		server := grpc.NewServer()
		grpcServer.NewGamificationServer(gamificationService).Register(server)
		reflection.Register(server)

		log.Infof("gRPC server listening on %s", grpcAddr)
		if err := server.Serve(lis); err != nil {
			log.Fatal(err, "Failed to serve gRPC")
		}
	}()

	// HTTP endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	})

	http.HandleFunc("/api/v1/coins/balance", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		wallet, _ := gamificationService.GetBalance(r.Context(), userID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wallet)
	})

	http.HandleFunc("/api/v1/coins/earn", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			UserID      string `json:"user_id"`
			Amount      int    `json:"amount"`
			Source      string `json:"source"`
			Description string `json:"description"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		wallet, err := gamificationService.EarnCoins(r.Context(), req.UserID, req.Amount, req.Source, req.Description)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wallet)
	})

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
			"reward":       reward,
			"streak":       streak,
			"streak_bonus": streak >= 7,
		})
	})

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
		log.Infof("HTTP server listening on %s", addr)
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
