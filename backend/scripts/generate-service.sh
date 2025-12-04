#!/bin/bash

# Service Generator Script
# Generates skeleton structure for a new microservice

SERVICE_NAME=$1
SERVICE_CATEGORY=$2
SERVICE_DESCRIPTION=$3

if [ -z "$SERVICE_NAME" ] || [ -z "$SERVICE_CATEGORY" ]; then
    echo "Usage: ./generate-service.sh <service-name> <category> <description>"
    echo "Categories: transaction-core, catalog-discovery, user-social, communication, logistics-fulfillment, marketing-engagement, intelligence-analytics"
    exit 1
fi

SERVICE_DIR="services/${SERVICE_CATEGORY}/${SERVICE_NAME}"

echo "Generating ${SERVICE_NAME} in ${SERVICE_CATEGORY}..."

# Create directory structure
mkdir -p "${SERVICE_DIR}/cmd/server"
mkdir -p "${SERVICE_DIR}/internal/domain"
mkdir -p "${SERVICE_DIR}/internal/application"
mkdir -p "${SERVICE_DIR}/internal/infrastructure"
mkdir -p "${SERVICE_DIR}/internal/interfaces"
mkdir -p "${SERVICE_DIR}/proto/${SERVICE_NAME}/v1"

# Create go.mod
cat > "${SERVICE_DIR}/go.mod" <<EOF
module github.com/titan-commerce/backend/${SERVICE_NAME}

go 1.23

require (
	github.com/google/uuid v1.5.0
	github.com/titan-commerce/backend/pkg v0.0.0
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
)

replace github.com/titan-commerce/backend/pkg => ../../pkg
EOF

# Create main.go
cat > "${SERVICE_DIR}/cmd/server/main.go" <<EOF
package main

import (
	"fmt"
	"os"

	"github.com/titan-commerce/backend/pkg/config"
	"github.com/titan-commerce/backend/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\\n", err)
		os.Exit(1)
	}

	log := logger.New(logger.Config{
		Level:       cfg.LogLevel,
		ServiceName: cfg.ServiceName,
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("${SERVICE_NAME} starting...")
	log.Infof("gRPC server would listen on :%d", cfg.GRPCPort)
	
	// TODO: Implement full service logic following Order Service pattern
	
	select {}
}
EOF

# Create README.md
cat > "${SERVICE_DIR}/README.md" <<EOF
# ${SERVICE_NAME}

${SERVICE_DESCRIPTION}

## Quick Start

\`\`\`bash
export SERVICE_NAME=${SERVICE_NAME}
export CELL_ID=cell-001
go run cmd/server/main.go
\`\`\`

## Status

🚧 **Under Development** - Skeleton structure created
EOF

# Create Dockerfile
cat > "${SERVICE_DIR}/Dockerfile" <<EOF
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${SERVICE_NAME} ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/${SERVICE_NAME} .
EXPOSE 9000 8080
CMD ["./${SERVICE_NAME}"]
EOF

echo "✓ ${SERVICE_NAME} skeleton created at ${SERVICE_DIR}"
