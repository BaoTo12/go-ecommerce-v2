package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/flash-sale-service/internal/application"
	"github.com/titan-commerce/backend/flash-sale-service/internal/infrastructure/redis"
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
		ServiceName: "flash-sale-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Flash Sale Service starting...")

	// Initialize Redis repository for atomic operations
	repo, err := redis.NewFlashSaleRepository(cfg.RedisAddr, cfg.RedisPassword, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize flash sale repository")
	}

	// Initialize application service
	flashSaleService := application.NewFlashSaleService(repo, log)

	// HTTP endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Get PoW challenge
	http.HandleFunc("/api/v1/flash-sale/challenge", func(w http.ResponseWriter, r *http.Request) {
		saleID := r.URL.Query().Get("sale_id")
		userID := r.URL.Query().Get("user_id")
		challenge := flashSaleService.GetChallenge(saleID, userID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"challenge": challenge})
	})

	// Attempt purchase with PoW
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

	// List active flash sales
	http.HandleFunc("/api/v1/flash-sales/active", func(w http.ResponseWriter, r *http.Request) {
		sales, _ := flashSaleService.GetActiveFlashSales(r.Context())
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sales)
	})

	go func() {
		addr := fmt.Sprintf(":%d", cfg.HTTPPort)
		log.Infof("Flash Sale service listening on %s", addr)
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
