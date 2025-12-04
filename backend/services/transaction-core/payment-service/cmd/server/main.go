package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/payment-service/internal/application"
	"github.com/titan-commerce/backend/payment-service/internal/domain"
	"github.com/titan-commerce/backend/payment-service/internal/infrastructure/gateway/mock"
	"github.com/titan-commerce/backend/payment-service/internal/infrastructure/postgres"
	handler "github.com/titan-commerce/backend/payment-service/internal/interface/grpc"
	pb "github.com/titan-commerce/backend/payment-service/proto/payment/v1"
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

	log.Info("Payment Service starting...")

	// Initialize PostgreSQL repository
	paymentRepo, err := postgres.NewPaymentRepository(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize payment repository")
	}

	// Initialize payment gateways
	mockGateway := mock.NewMockPaymentGateway(log)
	gateways := map[domain.PaymentGateway]domain.PaymentGatewayProvider{
		domain.PaymentGatewayStripe: mockGateway,
		domain.PaymentGatewayPayPal: mockGateway,
		domain.PaymentGatewayAdyen:  mockGateway,
	}

	// Initialize application service
	paymentService := application.NewPaymentService(paymentRepo, gateways, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, handler.NewPaymentServiceServer(paymentService, log))

	// Start server
	go func() {
		log.Infof("gRPC server listening on :%d", cfg.GRPCPort)
		log.Info("Multi-gateway payment processing ready (Stripe/PayPal/Adyen)")
		log.Info("Idempotency protection enabled")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err, "Failed to serve")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Payment Service")
	grpcServer.GracefulStop()
	log.Info("Payment Service stopped")
}

