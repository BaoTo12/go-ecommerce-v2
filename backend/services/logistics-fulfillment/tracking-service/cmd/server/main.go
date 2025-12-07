package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/titan-commerce/backend/tracking-service/internal/application"
	"github.com/titan-commerce/backend/tracking-service/internal/infrastructure/memory"
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
		ServiceName: "tracking-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Tracking Service starting...")

	// Initialize in-memory repository
	repo := memory.NewTrackingRepository()

	// Initialize application service
	trackingService := application.NewTrackingService(repo, log)

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

	// GET /api/v1/tracking/:order_id
	http.HandleFunc("/api/v1/tracking/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/tracking/")
		parts := strings.Split(path, "/")

		if len(parts) == 0 || parts[0] == "" {
			http.Error(w, "order_id required", http.StatusBadRequest)
			return
		}

		orderID := parts[0]

		// Handle live location endpoint (SSE)
		if len(parts) == 2 && parts[1] == "live" {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			ch := trackingService.Subscribe(orderID)
			defer trackingService.Unsubscribe(orderID, ch)

			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "SSE not supported", http.StatusInternalServerError)
				return
			}

			for update := range ch {
				data, _ := json.Marshal(update)
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			}
			return
		}

		// Get tracking details
		tracking, err := trackingService.GetTracking(r.Context(), orderID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if tracking == nil {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tracking)
	}))

	// POST /api/v1/tracking/:order_id/location - Update driver location
	http.HandleFunc("/api/v1/tracking-update/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		orderID := strings.TrimPrefix(r.URL.Path, "/api/v1/tracking-update/")

		var req struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		err := trackingService.UpdateLocation(r.Context(), orderID, req.Lat, req.Lng)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	}))

	// Start HTTP server
	port := cfg.HTTPPort
	if port == 0 {
		port = 8084
	}

	go func() {
		addr := fmt.Sprintf(":%d", port)
		log.Infof("HTTP server listening on %s", addr)
		log.Info("Endpoints:")
		log.Info("  GET  /api/v1/tracking/:order_id         - Get tracking info")
		log.Info("  GET  /api/v1/tracking/:order_id/live    - SSE live updates")
		log.Info("  POST /api/v1/tracking-update/:order_id  - Update location")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Tracking Service")
}
