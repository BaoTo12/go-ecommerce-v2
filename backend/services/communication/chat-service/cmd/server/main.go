﻿package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/chat-service/internal/application"
	"github.com/titan-commerce/backend/chat-service/internal/infrastructure/mongodb"
	ws "github.com/titan-commerce/backend/chat-service/internal/interface/websocket"
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
		ServiceName: cfg.ServiceName,
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Chat Service starting...")

	// Initialize MongoDB repository
	repo, err := mongodb.NewChatRepository(cfg.MongoURI, cfg.MongoDatabase, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize chat repository")
	}

	// Initialize application service
	chatService := application.NewChatService(repo, log)

	// Initialize WebSocket handler
	wsHandler := ws.NewChatWebSocketHandler(chatService, log)

	// HTTP server with WebSocket endpoint
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start HTTP server
	go func() {
		addr := fmt.Sprintf(":%d", cfg.HTTPPort)
		log.Infof("WebSocket server listening on %s", addr)
		log.Info("Real-time chat with WebSocket ready")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err, "Failed to serve HTTP")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Chat Service")
	log.Info("Chat Service stopped")
}
