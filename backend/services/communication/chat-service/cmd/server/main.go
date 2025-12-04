﻿package main

import (
	"fmt"
	"os"

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
		ServiceName: "chat-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("💬 Chat Service starting...")
	log.Infof("Cell: %s", cfg.CellID)
	log.Infof("gRPC server would listen on :%d", cfg.GRPCPort)
	log.Infof("WebSocket server would listen on :%d", cfg.HTTPPort)

	// TODO: Initialize ScyllaDB connection
	// TODO: Initialize repositories
	// TODO: Initialize application service
	// TODO: Initialize gRPC server
	// TODO: Initialize WebSocket server
	// TODO: Register handlers
	// TODO: Start servers

	log.Info("Chat Service ready for real-time messaging")

	select {}
}
