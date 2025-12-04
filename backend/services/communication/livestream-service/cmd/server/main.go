package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/livestream-service/internal/application"
	"github.com/titan-commerce/backend/livestream-service/internal/infrastructure/mongodb"
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
		ServiceName: "livestream-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Livestream Service starting...")

	// Initialize MongoDB repository
	repo, err := mongodb.NewLivestreamRepository(cfg.MongoURI, cfg.MongoDatabase, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize livestream repository")
	}

	// Initialize application service
	livestreamService := application.NewLivestreamService(repo, log)

	// HTTP endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// REST API endpoints
	http.HandleFunc("/api/v1/streams", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Create stream
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"status":"created"}`))
		} else {
			// List streams
			streams, _ := livestreamService.GetLiveStreams(r.Context(), 10, 0)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"streams":%d}`, len(streams))
		}
	})

	// RTMP callback for stream start
	http.HandleFunc("/rtmp/on_publish", func(w http.ResponseWriter, r *http.Request) {
		streamKey := r.URL.Query().Get("name")
		if streamKey == "" {
			http.Error(w, "stream key required", http.StatusBadRequest)
			return
		}
		_, err := livestreamService.StartStream(r.Context(), streamKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	go func() {
		addr := fmt.Sprintf(":%d", cfg.HTTPPort)
		log.Infof("HTTP server listening on %s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Livestream Service")
}
