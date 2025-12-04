package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/notification-service/internal/application"
	"github.com/titan-commerce/backend/notification-service/internal/infrastructure/postgres"
	"github.com/titan-commerce/backend/notification-service/internal/infrastructure/sender"
	"github.com/titan-commerce/backend/notification-service/internal/interface/grpc"
	pb "github.com/titan-commerce/backend/notification-service/proto/notification/v1"
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

	log.Info("Notification Service starting...")

	// Initialize PostgreSQL repository
	repo, err := postgres.NewNotificationRepository(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize notification repository")
	}

	// Initialize mock sender (replace with real Email/SMS/Push providers in production)
	mockSender := sender.NewMockNotificationSender(log)

	// Initialize application service
	notificationService := application.NewNotificationService(repo, mockSender, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, grpc.NewNotificationServiceServer(notificationService, log))

	// Start server
	go func() {
		log.Infof("gRPC server listening on :%d", cfg.GRPCPort)
		log.Info("Multi-channel notification system ready (Email/SMS/Push/In-App)")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err, "Failed to serve")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Notification Service")
	grpcServer.GracefulStop()
	log.Info("Notification Service stopped")
}

