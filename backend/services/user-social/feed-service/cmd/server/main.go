package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/feed-service/internal/application"
	"github.com/titan-commerce/backend/feed-service/internal/infrastructure/postgres"
	"github.com/titan-commerce/backend/feed-service/internal/interface/grpc"
	pb "github.com/titan-commerce/backend/feed-service/proto/feed/v1"
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

	log.Info("Feed Service starting...")

	// Initialize PostgreSQL repository
	repo, err := postgres.NewFeedRepository(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize feed repository")
	}

	// Initialize application service
	feedService := application.NewFeedService(repo, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	pb.RegisterFeedServiceServer(grpcServer, grpc.NewFeedServiceServer(feedService, log))

	// Start server
	go func() {
		log.Infof("gRPC server listening on :%d", cfg.GRPCPort)
		log.Info("TikTok-style feed engine ready")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err, "Failed to serve")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Feed Service")
	grpcServer.GracefulStop()
	log.Info("Feed Service stopped")
}
