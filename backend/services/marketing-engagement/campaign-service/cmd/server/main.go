package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/titan-commerce/backend/campaign-service/internal/application"
	"github.com/titan-commerce/backend/campaign-service/internal/domain"
	"github.com/titan-commerce/backend/campaign-service/internal/infrastructure/postgres"
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
		ServiceName: "campaign-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Campaign Service starting...")

	// Initialize repository
	repo := postgres.NewCampaignRepository()

	// Initialize application service
	campaignService := application.NewCampaignService(repo, log)

	// HTTP endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create campaign
	http.HandleFunc("/api/v1/campaigns", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var req struct {
				Name        string  `json:"name"`
				Description string  `json:"description"`
				Type        string  `json:"type"`
				StartTime   string  `json:"start_time"`
				EndTime     string  `json:"end_time"`
				Budget      float64 `json:"budget"`
			}
			json.NewDecoder(r.Body).Decode(&req)

			start, _ := time.Parse(time.RFC3339, req.StartTime)
			end, _ := time.Parse(time.RFC3339, req.EndTime)

			campaign, err := campaignService.CreateCampaign(r.Context(), req.Name, req.Description, domain.CampaignType(req.Type), start, end, req.Budget)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(campaign)
		} else {
			campaigns, _ := campaignService.GetActiveCampaigns(r.Context())
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(campaigns)
		}
	})

	// Get campaign stats
	http.HandleFunc("/api/v1/campaigns/stats", func(w http.ResponseWriter, r *http.Request) {
		campaignID := r.URL.Query().Get("campaign_id")
		stats, _ := campaignService.GetCampaignStats(r.Context(), campaignID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	// Record conversion
	http.HandleFunc("/api/v1/campaigns/conversion", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			CampaignID string  `json:"campaign_id"`
			OrderValue float64 `json:"order_value"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		campaignService.RecordConversion(r.Context(), req.CampaignID, req.OrderValue)
		w.WriteHeader(http.StatusOK)
	})

	go func() {
		addr := fmt.Sprintf(":%d", cfg.HTTPPort)
		log.Infof("Campaign service listening on %s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Campaign Service")
}
