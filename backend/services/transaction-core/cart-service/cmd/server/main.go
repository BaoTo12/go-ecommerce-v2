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

	log.Info("Cart Service starting...")
	log.Infof("Using Redis at %s", cfg.RedisAddr)
	
	// TODO: Implement cart service
	// - Redis client initialization
	// - gRPC API: AddToCart, RemoveFromCart, GetCart, ClearCart
	// - Auto-save mechanism
	// - TTL management
	
	select {}
}
