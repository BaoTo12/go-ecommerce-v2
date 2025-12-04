package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/refund-service/internal/application"
	"github.com/titan-commerce/backend/refund-service/internal/infrastructure/postgres"
	handler "github.com/titan-commerce/backend/refund-service/internal/interface/grpc"
	pb "github.com/titan-commerce/backend/refund-service/proto/refund/v1"
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

	log.Info("Refund Service starting...")

	// Initialize PostgreSQL repository
	refundRepo, err := postgres.NewRefundRepository(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize refund repository")
	}

	// Initialize application service (no gateway for now, handled internally)
	refundService := application.NewRefundService(refundRepo, nil, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	pb.RegisterRefundServiceServer(grpcServer, handler.NewRefundServiceServer(refundService, log))

	// Start server
	go func() {
		log.Infof("gRPC server listening on :%d", cfg.GRPCPort)
		log.Info("Automated refund processing ready")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err, "Failed to serve")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Refund Service")
	grpcServer.GracefulStop()
	log.Info("Refund Service stopped")
}
