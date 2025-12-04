package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/order-service/internal/application"
	"github.com/titan-commerce/backend/order-service/internal/infrastructure/postgres"
	handler "github.com/titan-commerce/backend/order-service/internal/interfaces/grpc"
	"github.com/titan-commerce/backend/pkg/config"
	"github.com/titan-commerce/backend/pkg/logger"
	grpcLib "google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(logger.Config{
		Level:       cfg.LogLevel,
		ServiceName: cfg.ServiceName,
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("Starting Order Service")

	// Initialize repositories
	eventStore, err := postgres.NewEventStoreRepository(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal(err, "Failed to connect to event store")
	}

	orderRepo, err := postgres.NewOrderReadModelRepository(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal(err, "Failed to connect to read model")
	}

	// Initialize application service
	orderService := application.NewOrderService(orderRepo, eventStore, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	handler.NewOrderServiceServer(grpcServer, orderService, log)

	// Start server in goroutine
	go func() {
		log.Infof("gRPC server listening on :%d", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err, "Failed to serve")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Order Service")
	grpcServer.GracefulStop()
	log.Info("Order Service stopped")
}
