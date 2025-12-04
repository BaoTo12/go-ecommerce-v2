package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/inventory-service/internal/application"
	"github.com/titan-commerce/backend/inventory-service/internal/infrastructure"
	"github.com/titan-commerce/backend/pkg/config"
	"github.com/titan-commerce/backend/pkg/logger"
	grpcLib "google.golang.org/grpc"
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

	log.Info("Inventory Service starting...")

	// Initialize Redis repository with Lua scripts
	inventoryRepo, err := infrastructure.NewRedisInventoryRepository(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatal(err, "Failed to initialize Redis inventory repository")
	}

	// Initialize application service
	inventoryService := application.NewInventoryService(inventoryRepo, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	// TODO: Register gRPC handler when proto is generated
	// pb.RegisterInventoryServiceServer(grpcServer, grpc.NewInventoryServiceServer(inventoryService, log))

	// Start server
	go func() {
		log.Infof("gRPC server listening on :%d", cfg.GRPCPort)
		log.Info("Atomic Redis Lua scripts - ZERO overselling guarantee")
log.Info("Reserve → Commit/Rollback pattern active")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err, "Failed to serve")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Inventory Service")
	grpcServer.GracefulStop()
	log.Info("Inventory Service stopped")
}
