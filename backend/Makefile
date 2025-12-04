.PHONY: help proto build test clean docker-build

help:
	@echo "Titan Commerce Platform - Backend Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  proto          - Generate Protocol Buffer code for all services"
	@echo "  build          - Build all microservices"
	@echo "  test           - Run tests for all services"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker images for all services"

# Generate Protocol Buffer code
proto:
	@echo "Generating Protocol Buffer code..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		services/transaction-core/order-service/proto/order/v1/order.proto

# Build all services
build:
	@echo "Building all services..."
	@cd services/transaction-core/order-service && go build -o bin/order-service ./cmd/server
	@echo "✓ Order Service built"

# Run tests
test:
	@echo "Running tests..."
	@go test ./pkg/...
	@go test ./services/transaction-core/order-service/internal/...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@find . -name "bin" -type d -exec rm -rf {} + 2>/dev/null ||  true
	@find . -name "*.test" -delete
	@echo "✓ Clean complete"

# Build Docker images
docker-build:
	@echo "Building Docker images..."
	@cd services/transaction-core/order-service && docker build -t titan-commerce/order-service:latest .
	@echo "✓ Docker images built"

# Run Order Service locally
run-order-service:
	@echo "Starting Order Service..."
	@cd services/transaction-core/order-service && go run ./cmd/server
