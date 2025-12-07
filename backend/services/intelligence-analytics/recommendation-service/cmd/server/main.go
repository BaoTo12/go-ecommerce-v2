package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/recommendation-service/internal/application"
	"github.com/titan-commerce/backend/recommendation-service/internal/infrastructure/memory"
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
		ServiceName: "recommendation-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Recommendation Service starting...")

	// Initialize in-memory repository (use PostgreSQL in production)
	repo := memory.NewRecommendationRepository()

	// Initialize application service
	recService := application.NewRecommendationService(repo, log)

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

	// GET /api/v1/recommendations/also-bought?product_id=
	http.HandleFunc("/api/v1/recommendations/also-bought", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		productID := r.URL.Query().Get("product_id")
		if productID == "" {
			http.Error(w, "product_id is required", http.StatusBadRequest)
			return
		}

		result, err := recService.GetAlsoBought(r.Context(), productID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}))

	// GET /api/v1/recommendations/personalized?user_id=
	http.HandleFunc("/api/v1/recommendations/personalized", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")

		result, err := recService.GetPersonalized(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}))

	// GET /api/v1/recommendations/frequently-bought?product_id=
	http.HandleFunc("/api/v1/recommendations/frequently-bought", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		productID := r.URL.Query().Get("product_id")

		result, err := recService.GetFrequentlyBoughtTogether(r.Context(), productID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}))

	// GET /api/v1/recommendations/because-viewed?product_id=
	http.HandleFunc("/api/v1/recommendations/because-viewed", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		productID := r.URL.Query().Get("product_id")

		result, err := recService.GetBecauseYouViewed(r.Context(), productID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}))

	// GET /api/v1/recommendations/trending?location=
	http.HandleFunc("/api/v1/recommendations/trending", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		location := r.URL.Query().Get("location")

		result, err := recService.GetTrendingNearYou(r.Context(), location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}))

	// POST /api/v1/recommendations/record-view
	http.HandleFunc("/api/v1/recommendations/record-view", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			UserID    string `json:"user_id"`
			ProductID string `json:"product_id"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		err := recService.RecordView(r.Context(), req.UserID, req.ProductID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "recorded"})
	}))

	// Start HTTP server
	port := cfg.HTTPPort
	if port == 0 {
		port = 8081
	}

	go func() {
		addr := fmt.Sprintf(":%d", port)
		log.Infof("HTTP server listening on %s", addr)
		log.Info("Endpoints:")
		log.Info("  GET /api/v1/recommendations/also-bought?product_id=")
		log.Info("  GET /api/v1/recommendations/personalized?user_id=")
		log.Info("  GET /api/v1/recommendations/frequently-bought?product_id=")
		log.Info("  GET /api/v1/recommendations/because-viewed?product_id=")
		log.Info("  GET /api/v1/recommendations/trending?location=")
		log.Info("  POST /api/v1/recommendations/record-view")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Recommendation Service")
}
