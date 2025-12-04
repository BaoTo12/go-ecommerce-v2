package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/flash-sale-service/internal/application"
	grpcServer "github.com/titan-commerce/backend/flash-sale-service/internal/interface/grpc"
	"github.com/titan-commerce/backend/flash-sale-service/internal/infrastructure/postgres"
	"github.com/titan-commerce/backend/flash-sale-service/internal/infrastructure/redis"
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
		ServiceName: "flash-sale-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Flash Sale Service starting...")

	// Initialize PostgreSQL repository for persistence
	pgRepo, err := postgres.NewFlashSalePostgresRepository(cfg.DatabaseURL, log)
	if err != nil {
		log.Warnf("PostgreSQL not available, using Redis only: %v", err)
	}

	// Initialize Redis repository for atomic operations
	redisRepo, err := redis.NewFlashSaleRepository(cfg.RedisAddr, cfg.RedisPassword, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize Redis repository")
	}

	// Use PostgreSQL if available, otherwise fallback to Redis-based persistence
	var repo application.FlashSaleRepository
	if pgRepo != nil {
		repo = pgRepo
	} else {
		repo = redisRepo
	}

	// Initialize application service
	flashSaleService := application.NewFlashSaleService(repo, log)

	// Start gRPC server
	go func() {
		grpcAddr := fmt.Sprintf(":%d", cfg.GRPCPort)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatal(err, "Failed to listen for gRPC")
		}

		server := grpc.NewServer()
		grpcServer.NewFlashSaleServer(flashSaleService).Register(server)
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

	http.HandleFunc("/api/v1/flash-sale/challenge", func(w http.ResponseWriter, r *http.Request) {
		saleID := r.URL.Query().Get("sale_id")
		userID := r.URL.Query().Get("user_id")
		challenge := flashSaleService.GetChallenge(saleID, userID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"challenge":  challenge,
			"difficulty": 4,
		})
	})

	http.HandleFunc("/api/v1/flash-sale/purchase", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			SaleID    string `json:"flash_sale_id"`
			UserID    string `json:"user_id"`
			Quantity  int    `json:"quantity"`
			Challenge string `json:"challenge"`
			Nonce     string `json:"nonce"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		reservation, err := flashSaleService.AttemptPurchase(r.Context(), req.SaleID, req.UserID, req.Quantity, req.Challenge, req.Nonce)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"reservation_id": reservation.ID,
			"expires_at":     reservation.ExpiresAt,
		})
	})

	http.HandleFunc("/api/v1/flash-sales/active", func(w http.ResponseWriter, r *http.Request) {
		sales, _ := flashSaleService.GetActiveFlashSales(r.Context())
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sales)
	})

	go func() {
		addr := fmt.Sprintf(":%d", cfg.HTTPPort)
		log.Infof("HTTP server listening on %s", addr)
		log.Info("PoW-protected flash sale ready for 11.11!")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Flash Sale Service")
}
