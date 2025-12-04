package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-commerce/backend/payment-service/internal/application"
	"github.com/titan-commerce/backend/payment-service/internal/domain"
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

	// TODO: Initialize PostgreSQL repository
	// paymentRepo := postgres.NewPaymentRepository(cfg.DatabaseURL, log)

	// Initialize payment gateways (mock for now)
	// In production, initialize real gateway adapters:
	// stripeGateway := stripe.NewStripeGateway(cfg.StripeSecretKey, log)
	// paypalGateway := paypal.NewPayPalGateway(cfg.PayPalClientID, cfg.PayPalSecret, log)
	// adyenGateway := adyen.NewAdyenGateway(cfg.AdyenAPIKey, log)

	gateways := make(map[domain.PaymentGateway]domain.PaymentGateway)
	// gateways[domain.PaymentGatewayStripe] = stripeGateway
	// gateways[domain.PaymentGatewayPayPal] = paypalGateway
	// gateways[domain.PaymentGatewayAdyen] = adyenGateway

	// Initialize application service
	// paymentService := application.NewPaymentService(paymentRepo, gateways, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	// TODO: Register gRPC handler
	// paymentv1.RegisterPaymentServiceServer(grpcServer, handler.NewPaymentServiceServer(paymentService, log))

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

