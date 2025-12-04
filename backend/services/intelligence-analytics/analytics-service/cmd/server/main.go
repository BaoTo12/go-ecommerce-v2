package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/titan-commerce/backend/analytics-service/internal/application"
	"github.com/titan-commerce/backend/analytics-service/internal/domain"
	"github.com/titan-commerce/backend/analytics-service/internal/infrastructure/clickhouse"
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
		ServiceName: "analytics-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Analytics Service starting...")

	// Initialize repository
	repo := clickhouse.NewAnalyticsRepository()

	// Initialize application service
	analyticsService := application.NewAnalyticsService(repo, log)

	// HTTP endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Track event
	http.HandleFunc("/api/v1/track", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			UserID     string                 `json:"user_id"`
			SessionID  string                 `json:"session_id"`
			EventType  string                 `json:"event_type"`
			Properties map[string]interface{} `json:"properties"`
			DeviceType string                 `json:"device_type"`
			Platform   string                 `json:"platform"`
			Country    string                 `json:"country"`
			City       string                 `json:"city"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		err := analyticsService.TrackEvent(
			r.Context(),
			req.UserID,
			req.SessionID,
			domain.EventType(req.EventType),
			req.Properties,
			req.DeviceType,
			req.Platform,
			req.Country,
			req.City,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	})

	// Get dashboard
	http.HandleFunc("/api/v1/analytics/dashboard", func(w http.ResponseWriter, r *http.Request) {
		dashboard, _ := analyticsService.GetDashboard(r.Context())
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dashboard)
	})

	// Get sales report
	http.HandleFunc("/api/v1/analytics/sales", func(w http.ResponseWriter, r *http.Request) {
		period := r.URL.Query().Get("period")
		if period == "" {
			period = "daily"
		}
		report, _ := analyticsService.GetSalesReport(r.Context(), period)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
	})

	// Get conversion funnel
	http.HandleFunc("/api/v1/analytics/funnel", func(w http.ResponseWriter, r *http.Request) {
		funnelName := r.URL.Query().Get("name")
		period := r.URL.Query().Get("period")
		funnel, _ := analyticsService.GetConversionFunnel(r.Context(), funnelName, period)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(funnel)
	})

	// Calculate metrics
	http.HandleFunc("/api/v1/analytics/metrics", func(w http.ResponseWriter, r *http.Request) {
		end := time.Now()
		start := end.AddDate(0, 0, -1)
		metrics, _ := analyticsService.CalculateMetrics(r.Context(), start, end)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics)
	})

	go func() {
		addr := fmt.Sprintf(":%d", cfg.HTTPPort)
		log.Infof("Analytics service listening on %s", addr)
		log.Info("Real-time dashboards, funnels, and cohorts ready!")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Analytics Service")
}
