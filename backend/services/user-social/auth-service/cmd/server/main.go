package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/auth-service/internal/application"
	"github.com/titan-commerce/backend/auth-service/internal/infrastructure/postgres"
	"github.com/titan-commerce/backend/auth-service/internal/infrastructure/redis"
	"github.com/titan-commerce/backend/auth-service/internal/infrastructure/token"
	"github.com/titan-commerce/backend/auth-service/internal/interface/grpc"
	pb "github.com/titan-commerce/backend/auth-service/proto/auth/v1"
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

	log.Info("🔐 Auth Service starting...")

	// Initialize PostgreSQL repository
	authRepo, err := postgres.NewAuthRepository(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize auth repository")
	}

	// Initialize Redis repository
	tokenRepo, err := redis.NewRedisRepository(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatal(err, "Failed to initialize redis repository")
	}

	// Initialize Token Service
	tokenService := token.NewTokenService(cfg.JWTSecret, cfg.JWTRefreshSecret)

	// Initialize Application Service
	authService := application.NewAuthService(authRepo, tokenRepo, tokenService, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, grpc.NewAuthServiceServer(authService, log))

	// Start server
	go func() {
		log.Infof("gRPC server listening on :%d", cfg.GRPCPort)
		log.Info("JWT + MFA + Redis Blacklist operational")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err, "Failed to serve")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Auth Service")
	grpcServer.GracefulStop()
	log.Info("Auth Service stopped")
}

