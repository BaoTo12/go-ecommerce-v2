package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/titan-commerce/backend/services/logistics-fulfillment/tracking-service/internal/application"
	"github.com/titan-commerce/backend/services/logistics-fulfillment/tracking-service/internal/infrastructure/memory"
	"github.com/titan-commerce/backend/pkg/server"
)

/*
Tracking Service - Updated to use Unified Server Package

Features automatically enabled:
- Security Headers
- Rate Limiting  
- Response Compression
- Distributed Tracing
- Metrics Collection
- Graceful Shutdown
- SSE support for real-time updates
*/

func main() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8084"
	}

	// Create unified server
	config := server.DefaultConfig()
	config.Port = port
	config.ServiceName = "tracking-service"
	config.CORSOrigins = []string{"http://localhost:3000", "*"}
	config.CSRFEnabled = false // API-only service
	config.CompressionEnabled = false // SSE doesn't work with compression

	srv := server.New(config)

	// Initialize service layer
	repo := memory.NewTrackingRepository()
	trackingService := application.NewTrackingService(repo, nil)

	// Register handlers
	srv.HandleFunc("GET /api/v1/tracking/{order_id}", handleGetTracking(trackingService, srv))
	srv.HandleFunc("GET /api/v1/tracking/{order_id}/live", handleLiveTracking(trackingService, srv))
	srv.HandleFunc("POST /api/v1/tracking-update/{order_id}", handleUpdateLocation(trackingService, srv))

	log.Println("Tracking Service starting with unified server...")
	log.Println("Endpoints:")
	log.Println("  GET  /api/v1/tracking/{order_id}        - Get tracking info")
	log.Println("  GET  /api/v1/tracking/{order_id}/live   - SSE live updates")
	log.Println("  POST /api/v1/tracking-update/{order_id} - Update location")
	log.Println("  GET  /health                            - Health check (built-in)")
	log.Println("  GET  /metrics                           - Metrics (built-in)")

	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleGetTracking(svc *application.TrackingService, srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.PathValue("order_id")
		if orderID == "" {
			http.Error(w, "order_id required", http.StatusBadRequest)
			return
		}

		// Track metric
		srv.Metrics().Counter("tracking_requests").Inc()

		tracking, err := svc.GetTracking(r.Context(), orderID)
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
	}
}

func handleLiveTracking(svc *application.TrackingService, srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.PathValue("order_id")
		if orderID == "" {
			http.Error(w, "order_id required", http.StatusBadRequest)
			return
		}

		// Track metric
		srv.Metrics().Counter("sse_connections").Inc()
		srv.Metrics().Gauge("sse_active").Inc()
		defer srv.Metrics().Gauge("sse_active").Dec()

		// SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ch := svc.Subscribe(orderID)
		defer svc.Unsubscribe(orderID, ch)

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
	}
}

func handleUpdateLocation(svc *application.TrackingService, srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.PathValue("order_id")
		if orderID == "" {
			http.Error(w, "order_id required", http.StatusBadRequest)
			return
		}

		var req struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Track metric
		srv.Metrics().Counter("location_updates").Inc()

		err := svc.UpdateLocation(r.Context(), orderID, req.Lat, req.Lng)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	}
}
