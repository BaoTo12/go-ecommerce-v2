package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/videocall-service/internal/application"
	"github.com/titan-commerce/backend/videocall-service/internal/infrastructure/postgres"
	ws "github.com/titan-commerce/backend/videocall-service/internal/interface/websocket"
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
		ServiceName: "videocall-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Video Call Service starting...")

	// Initialize repository
	repo := postgres.NewVideoCallRepository()

	// Initialize application service
	callService := application.NewVideoCallService(repo, log)

	// Initialize WebSocket handler for signaling
	wsHandler := ws.NewSignalingHandler(callService, log)

	// HTTP endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// WebSocket endpoint for WebRTC signaling
	http.HandleFunc("/ws/signaling", wsHandler.HandleSignaling)

	// REST API for call history
	http.HandleFunc("/api/v1/calls/history", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		calls, _ := callService.GetCallHistory(r.Context(), userID, 20)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"calls":%d}`, len(calls))
	})

	go func() {
		addr := fmt.Sprintf(":%d", cfg.HTTPPort)
		log.Infof("WebSocket signaling server listening on %s", addr)
		log.Info("WebRTC signaling ready for peer-to-peer video calls")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Video Call Service")
}
