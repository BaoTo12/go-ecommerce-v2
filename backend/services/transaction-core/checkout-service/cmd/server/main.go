package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/checkout-service/internal/application"
	"github.com/titan-commerce/backend/checkout-service/internal/infrastructure/mock"
	"github.com/titan-commerce/backend/checkout-service/internal/interface/grpc"
	pb "github.com/titan-commerce/backend/checkout-service/proto/checkout/v1"
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

	log.Info("Checkout Service starting - Saga Coordinator ready")

	// Initialize Clients (Mock for now, replace with real gRPC clients in production)
	invClient := &mock.MockInventoryClient{}
	payClient := &mock.MockPaymentClient{}
	ordClient := &mock.MockOrderClient{}
	crtClient := &mock.MockCartClient{}

	// Initialize Application Service (Saga Orchestrator)
	checkoutService := application.NewCheckoutService(invClient, payClient, ordClient, crtClient, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	pb.RegisterCheckoutServiceServer(grpcServer, grpc.NewCheckoutServiceServer(checkoutService, log))

	// Start server
	go func() {
		log.Infof("gRPC server listening on :%d", cfg.GRPCPort)
		log.Info("Saga orchestrator: Reserve → Pay → Order → Commit")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err, "Failed to serve")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Checkout Service")
	grpcServer.GracefulStop()
	log.Info("Checkout Service stopped")
}

