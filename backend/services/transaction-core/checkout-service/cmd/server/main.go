package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/checkout-service/internal/domain"
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

	// Initialize gRPC clients for service orchestration
	// TODO: Connect to:
	// - Inventory Service (Reserve stock)
	// - Payment Service (Process payment)
	// - Order Service (Create order)
	// - Cart Service (Clear cart)

	// Saga Orchestration Flow:
	// 1. Reserve inventory (compensate: rollback reservation)
	// 2. Process payment (compensate: refund)
	// 3. Create order (compensate: cancel order)
	// 4. Commit inventory reservation
	// 5. Clear cart

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	// TODO: Register gRPC handler
	// checkoutv1.RegisterCheckoutServiceServer(grpcServer, handler.NewCheckoutServiceServer(sagaOrchestrator, log))

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
