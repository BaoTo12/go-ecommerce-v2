package main

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
		ServiceName: "driver-service",
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("🚗 Driver Service starting...")
	log.Infof("Cell: %s", cfg.CellID)
	log.Infof("gRPC server would listen on :%d", cfg.GRPCPort)

	// TODO: Initialize PostgreSQL connection
	// TODO: Initialize repositories
	// TODO: Initialize application service
	// TODO: Initialize gRPC server
	// TODO: Register gRPC handlers
	// TODO: Start server

	log.Info("Driver Service ready for last-mile delivery management")

	select {}
}

