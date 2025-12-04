package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/cart-service/internal/application"
	"github.com/titan-commerce/backend/cart-service/internal/infrastructure/redis"
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

	log.Info("Cart Service starting...")

	// Initialize Redis repository
	repo, err := redis.NewRedisCartRepository(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatal(err, "Failed to connect to Redis")
	}
	log.Infof("Connected to Redis at %s", cfg.RedisAddr)

	// Initialize application service
	cartService := application.NewCartService(repo, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	// TODO: Register gRPC handler
	// cartv1.RegisterCartServiceServer(grpcServer, handler.NewCartServiceServer(cartService, log))

	// Start server
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

	log.Info("Shutting down Cart Service")
	grpcServer.GracefulStop()
	log.Info("Cart Service stopped")
}
