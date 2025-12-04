package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/wallet-service/internal/application"
	"github.com/titan-commerce/backend/wallet-service/internal/infrastructure/postgres"
	"github.com/titan-commerce/backend/wallet-service/internal/interface/grpc"
	pb "github.com/titan-commerce/backend/wallet-service/proto/wallet/v1"
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

	log.Info("Wallet Service starting...")

	// Initialize PostgreSQL repositories
	walletRepo, err := postgres.NewWalletRepository(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize wallet repository")
	}

	txnRepo, err := postgres.NewTransactionRepository(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize transaction repository")
	}

	// Initialize application service
	walletService := application.NewWalletService(walletRepo, txnRepo, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	pb.RegisterWalletServiceServer(grpcServer, grpc.NewWalletServiceServer(walletService, log))

	// Start server
	go func() {
		log.Infof("gRPC server listening on :%d", cfg.GRPCPort)
		log.Info("Escrow system operational - Available + Held balances")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err, "Failed to serve")
		}
	}()

	log.Info("Wallet Service ready - Escrow system operational")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Wallet Service")
	grpcServer.GracefulStop()
	log.Info("Wallet Service stopped")
}

