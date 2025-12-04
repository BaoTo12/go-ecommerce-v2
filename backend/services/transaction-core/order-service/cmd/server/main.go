package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/order-service/internal/application"
	"github.com/titan-commerce/backend/order-service/internal/infrastructure/postgres"
	"github.com/titan-commerce/backend/order-service/internal/interfaces/grpc"
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
		Pretty:      true, // Set to false in production
	})

	log.Info("Starting Order Service")

	// Initialize database
	db, err := postgres.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err, "Failed to connect to database")
	}
	defer db.Close()

	// Initialize repositories
	orderRepo := postgres.NewOrderRepository(db)
	eventStore := postgres.NewEventStore(db)

	// Initialize application service
	orderService := application.NewOrderService(orderRepo, eventStore, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	grpc.NewOrderServiceServer(grpcServer, orderService, log)

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
