package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"strings"

	"github.com/titan-commerce/backend/search-service/internal/application"
	"github.com/titan-commerce/backend/search-service/internal/infrastructure/elasticsearch"
	"github.com/titan-commerce/backend/search-service/internal/interface/grpc"
	pb "github.com/titan-commerce/backend/search-service/proto/search/v1"
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

	log.Info("Search Service starting...")

	// Initialize Elasticsearch repository
	// Assuming cfg.ElasticsearchURL is a comma-separated list of addresses
	addresses := strings.Split(cfg.ElasticsearchURL, ",")
	if len(addresses) == 0 || addresses[0] == "" {
		addresses = []string{"http://localhost:9200"} // Default
	}

	repo, err := elasticsearch.NewSearchRepository(addresses, log)
	if err != nil {
		log.Fatal(err, "Failed to initialize elasticsearch repository")
	}

	// Initialize application service
	searchService := application.NewSearchService(repo, log)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(err, "Failed to listen")
	}

	grpcServer := grpcLib.NewServer()
	pb.RegisterSearchServiceServer(grpcServer, grpc.NewSearchServiceServer(searchService, log))

	// Start server
	go func() {
		log.Infof("gRPC server listening on :%d", cfg.GRPCPort)
		log.Info("Elasticsearch engine ready")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err, "Failed to serve")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Search Service")
	grpcServer.GracefulStop()
	log.Info("Search Service stopped")
}
