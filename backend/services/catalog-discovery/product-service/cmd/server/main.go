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
		ServiceName: cfg.ServiceName,
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Product Service starting...")
	log.Infof("Using MongoDB at %s", cfg.MongoURI)
	
	// TODO: Implement product catalog
	// - MongoDB for flexible product schema
	// - Multi-variant products (size, color, etc.)
	// - Bulk import APIs for sellers
	// - Image upload to S3/MinIO
	
	select {}
}
