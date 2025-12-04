package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/fraud-service/internal/application"
	grpcServer "github.com/titan-commerce/backend/fraud-service/internal/interface/grpc"
	"github.com/titan-commerce/backend/fraud-service/internal/infrastructure/postgres"
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
		ServiceName: "fraud-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Fraud Detection Service starting...")

	// Initialize PostgreSQL repository
	repo, err := postgres.NewFraudRepository(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize repository")
	}
	defer repo.Close()

	// Initialize application service with ML scoring
	fraudService := application.NewFraudService(repo, log)

	// Start gRPC server
	go func() {
		grpcAddr := fmt.Sprintf(":%d", cfg.GRPCPort)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatal(err, "Failed to listen for gRPC")
		}

		server := grpc.NewServer()
		grpcServer.NewFraudServer(fraudService).Register(server)
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

	http.HandleFunc("/api/v1/fraud/check", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			TransactionID string  `json:"transaction_id"`
			UserID        string  `json:"user_id"`
			Amount        float64 `json:"amount"`
			Currency      string  `json:"currency"`
			IP            string  `json:"ip"`
			DeviceID      string  `json:"device_id"`
			UserAgent     string  `json:"user_agent"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		check, err := fraudService.CheckTransaction(
			r.Context(),
			req.TransactionID,
			req.UserID,
			req.Amount,
			req.Currency,
			req.IP,
			req.DeviceID,
			req.UserAgent,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"check_id":        check.ID,
			"score":           check.Score,
			"risk_level":      check.RiskLevel,
			"decision":        check.Decision,
			"reasons":         check.Reasons,
			"processing_time": check.ProcessingTime,
		})
	})

	http.HandleFunc("/api/v1/fraud/alerts", func(w http.ResponseWriter, r *http.Request) {
		alerts, _ := fraudService.GetPendingAlerts(r.Context())
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(alerts)
	})

	http.HandleFunc("/api/v1/fraud/history", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		history, _ := fraudService.GetUserFraudHistory(r.Context(), userID, 20)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(history)
	})

	go func() {
		addr := fmt.Sprintf(":%d", cfg.HTTPPort)
		log.Infof("HTTP server listening on %s", addr)
		log.Info("ML-based fraud scoring engine ready!")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Fraud Detection Service")
}
