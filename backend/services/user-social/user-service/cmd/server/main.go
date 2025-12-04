package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/user-service/internal/application"
	"github.com/titan-commerce/backend/user-service/internal/infrastructure/postgres"
	"github.com/titan-commerce/backend/user-service/internal/interface/grpc"
	pb "github.com/titan-commerce/backend/user-service/proto/user/v1"
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

	log.Info("User Service starting...")

	// Initialize PostgreSQL repository
	userRepo, err := postgres.NewUserRepository(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize user repository")
	}

	// Initialize application service
	userService := application.NewUserService(userRepo, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	pb.RegisterUserServiceServer(grpcServer, grpc.NewUserServiceServer(userService, log))

	// Start server
	go func() {
		log.Infof("gRPC server listening on :%d", cfg.GRPCPort)
		log.Info("User profile & address management ready")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err, "Failed to serve")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down User Service")
	grpcServer.GracefulStop()
	log.Info("User Service stopped")
}

