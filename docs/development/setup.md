# Development Setup Guide

## Quick Start

This guide will get you up and running with the Titan Commerce Platform development environment in under 30 minutes.

---

## Prerequisites

### Required Software

| Tool | Version | Purpose |
|------|---------|---------|
| **Node.js** | 20+ | Frontend development |
| **Go** | 1.23+ | Backend development |
| **Docker** | 24+ | Local infrastructure |
| **Docker Compose** | 2.20+ | Orchestration |
| **Git** | 2.40+ | Version control |

### Optional (Recommended)

- **VS Code** or **IntelliJ IDEA**
- **Postman** or **Insomnia** (API testing)
- **k9s** (Kubernetes management)
- **kubectx/kubens** (Context switching)

---

## Installation

### 1. Clone Repository

```bash
git clone https://github.com/titan-commerce/platform.git
cd platform
```

### 2. Install Frontend Dependencies

```bash
cd frontend
npm install

# Verify installation
npm run type-check
```

### 3. Install Backend Dependencies

```bash
cd backend
go mod download

# Verify installation
go version
go build ./...
```

### 4. Start Infrastructure (Docker Compose)

```bash
# From project root
docker-compose up -d

# Verify services
docker-compose ps
```

**Services Started:**
- PostgreSQL (port 5432)
- ScyllaDB (port 9042)
- Redis (port 6379)
- Kafka (port 9092)
- Elasticsearch (port 9200)
- Jaeger UI (port 16686)
- Prometheus (port 9090)
- Grafana (port 3000)

---

## Project Structure

```
titan-commerce-platform/
├── frontend/                # Next.js micro-frontends
│   ├── shell/              # Host application
│   ├── apps/
│   │   ├── discovery/      # Product browsing
│   │   ├── checkout/       # Cart & payment
│   │   └── seller-centre/  # Vendor portal
│   ├── packages/
│   │   ├── ui/             # Shared components
│   │   └── shared/         # Common utilities
│   ├── package.json
│   └── turbo.json
│
├── backend/                 # Golang microservices
│   ├── services/
│   │   ├── ingestion-engine/
│   │   ├── transaction-core/
│   │   ├── intelligence-layer/
│   │   ├── chat/
│   │   └── catalog/
│   ├── pkg/                # Shared libraries
│   ├── go.mod
│   └── Makefile
│
├── infrastructure/          # Deployment configs
│   ├── kubernetes/
│   ├── helm/
│   ├── terraform/
│   └── docker/
│
├── docs/                    # Documentation
│   ├── architecture/
│   ├── api/
│   ├── deployment/
│   └── development/
│
└── docker-compose.yml       # Local development stack
```

---

## Frontend Development

### Running Development Server

```bash
cd frontend

# Start all apps
npm run dev

# Start specific app
npm run dev --filter=discovery
```

**Access:**
- Shell (Host): http://localhost:3000
- Discovery: http://localhost:3001
- Checkout: http://localhost:3002
- Seller Centre: http://localhost:3003

### Creating a New Component

**1. Navigate to shared UI package:**
```bash
cd frontend/packages/ui
```

**2. Create component:**
```tsx
// src/Button/Button.tsx
import React from 'react';
import { cn } from '../utils';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ variant = 'primary', size = 'md', className, ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={cn(
          'btn',
          `btn-${variant}`,
          `btn-${size}`,
          className
        )}
        {...props}
      />
    );
  }
);

Button.displayName = 'Button';
```

**3. Export from index:**
```tsx
// src/index.ts
export * from './Button/Button';
```

**4. Use in app:**
```tsx
// apps/discovery/src/app/page.tsx
import { Button } from '@titan/ui';

export default function HomePage() {
  return (
    <Button variant="primary" onClick={() => console.log('Clicked')}>
      Browse Products
    </Button>
  );
}
```

### Testing

```bash
# Run all tests
npm run test

# Run tests for specific app
npm run test --filter=discovery

# Run tests in watch mode
npm run test:watch

# Generate coverage report
npm run test:coverage
```

### Linting & Formatting

```bash
# Lint all code
npm run lint

# Fix lint issues
npm run lint:fix

# Format code
npm run format
```

---

## Backend Development

### Running Services Locally

**Option 1: Run All Services**
```bash
cd backend
make dev
```

**Option 2: Run Specific Service**
```bash
# Order service
cd backend/services/transaction-core/order-processing
go run cmd/main.go

# With hot reload (requires air)
air
```

**Option 3: Docker Compose**
```bash
docker-compose up order-service inventory-service
```

### Creating a New Service

**1. Generate boilerplate:**
```bash
cd backend
mkdir -p services/my-service/{cmd,internal/{domain,ports,adapters},proto}
```

**2. Define protobuf:**
```protobuf
// services/my-service/proto/my_service.proto
syntax = "proto3";

package myservice.v1;

option go_package = "github.com/titan-commerce/backend/services/my-service/proto;myservicepb";

service MyService {
  rpc DoSomething(DoSomethingRequest) returns (DoSomethingResponse);
}

message DoSomethingRequest {
  string id = 1;
}

message DoSomethingResponse {
  string result = 1;
}
```

**3. Generate code:**
```bash
make proto
```

**4. Implement service:**
```go
// internal/adapters/grpc/server.go
package grpc

import (
    "context"
    pb "github.com/titan-commerce/backend/services/my-service/proto"
)

type Server struct {
    pb.UnimplementedMyServiceServer
}

func (s *Server) DoSomething(ctx context.Context, req *pb.DoSomethingRequest) (*pb.DoSomethingResponse, error) {
    // Business logic here
    return &pb.DoSomethingResponse{
        Result: "Success: " + req.Id,
    }, nil
}
```

**5. Main entry point:**
```go
// cmd/main.go
package main

import (
    "flag"
    "fmt"
    "log"
    "net"
    
    "google.golang.org/grpc"
    pb "github.com/titan-commerce/backend/services/my-service/proto"
    "github.com/titan-commerce/backend/services/my-service/internal/adapters/grpc"
)

func main() {
    port := flag.Int("port", 50051, "gRPC port")
    flag.Parse()
    
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    
    grpcServer := grpc.NewServer()
    pb.RegisterMyServiceServer(grpcServer, &grpc.Server{})
    
    log.Printf("Server listening on :%d", *port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
```

### Testing

**Unit Tests:**
```go
// internal/domain/order_test.go
package domain_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
    order := NewOrder("user123")
    
    assert.NotNil(t, order)
    assert.Equal(t, "user123", order.UserID)
    assert.Equal(t, "CREATED", order.Status)
}
```

**Run tests:**
```bash
# All tests
make test

# With coverage
make coverage

# Specific package
go test ./internal/domain/...

# Verbose
go test -v ./...
```

**Integration Tests:**
```bash
# Requires Docker
make test-integration
```

### Linting

```bash
# Run linter
make lint

# Auto-fix issues
golangci-lint run --fix
```

---

## Database Migrations

### Create Migration

```bash
cd backend

# Create new migration
migrate create -ext sql -dir migrations -seq add_users_table
```

**Files created:**
- `migrations/000001_add_users_table.up.sql`
- `migrations/000001_add_users_table.down.sql`

**Up migration:**
```sql
-- migrations/000001_add_users_table.up.sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Down migration:**
```sql
-- migrations/000001_add_users_table.down.sql
DROP TABLE users;
```

### Run Migrations

```bash
# Up
make migrate-up

# Down (rollback)
make migrate-down

# Specific version
migrate -path migrations -database "postgres://localhost/titan?sslmode=disable" goto 1
```

---

## Environment Variables

### Frontend (.env.local)

Create `frontend/.env.local`:

```env
# API Endpoints
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080

# Feature Flags
NEXT_PUBLIC_ENABLE_FLASH_SALES=true
NEXT_PUBLIC_ENABLE_CHAT=true

# Analytics
NEXT_PUBLIC_GA_ID=G-XXXXXXXXXX

# Stripe (Public Key)
NEXT_PUBLIC_STRIPE_KEY=pk_test_xxxxx
```

### Backend (.env)

Create `backend/.env`:

```env
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/titan?sslmode=disable
SCYLLA_NODES=localhost:9042
REDIS_URL=redis://localhost:6379
ELASTICSEARCH_URL=http://localhost:9200

# Kafka
KAFKA_BROKERS=localhost:9092

# Secrets
JWT_SECRET=your-super-secret-key-change-me
STRIPE_SECRET_KEY=sk_test_xxxxx

# Observability
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
JAEGER_ENDPOINT=http://localhost:14268/api/traces

# Service Ports
ORDER_SERVICE_PORT=50051
INVENTORY_SERVICE_PORT=50052
PAYMENT_SERVICE_PORT=50053
```

---

## Debugging

### Frontend (VS Code)

**.vscode/launch.json:**
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Next.js: debug server-side",
      "type": "node-terminal",
      "request": "launch",
      "command": "npm run dev",
      "cwd": "${workspaceFolder}/frontend"
    },
    {
      "name": "Next.js: debug client-side",
      "type": "chrome",
      "request": "launch",
      "url": "http://localhost:3000",
      "webRoot": "${workspaceFolder}/frontend"
    }
  ]
}
```

### Backend (Delve)

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug service
cd backend/services/transaction-core/order-processing
dlv debug cmd/main.go

# Set breakpoint
(dlv) break main.main
(dlv) continue
```

**VS Code launch.json:**
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Order Service",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/backend/services/transaction-core/order-processing/cmd",
      "env": {
        "DATABASE_URL": "postgres://localhost/titan"
      }
    }
  ]
}
```

---

## Common Tasks

### Reset Database

```bash
docker-compose down -v
docker-compose up -d postgres
make migrate-up
```

### Clear Redis Cache

```bash
docker exec -it titan-redis redis-cli FLUSHALL
```

### View Kafka Messages

```bash
docker exec -it titan-kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic order.events \
  --from-beginning
```

### Generate Sample Data

```bash
cd backend
go run scripts/seed-data.go
```

---

## Tips & Best Practices

### Frontend

1. **Use Server Components by default**, Client Components only when needed (interactivity)
2. **Colocate tests** with components: `Button.tsx` + `Button.test.tsx`
3. **Use TanStack Query** for server state, **Zustand** for client state
4. **Optimize images** with `next/image`
5. **Enable Strict Mode** in TypeScript

### Backend

1. **Follow Hexagonal Architecture**: Domain → Ports → Adapters
2. **Use context for cancellation**: Always pass `context.Context`
3. **Log structured data**: Use `zap` with fields
4. **Handle errors explicitly**: Never ignore `err`
5. **Write table-driven tests**

### General

1. **Commit often**, push daily
2. **Write meaningful commit messages**: "feat: add flash sale countdown"
3. **Run tests before pushing**: `npm run test && make test`
4. **Use feature branches**: `feature/flash-sale-ui`
5. **Request code reviews** for all PRs

---

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 3000
lsof -i :3000

# Kill process
kill -9 <PID>
```

### Module Not Found (Frontend)

```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### Go Module Issues

```bash
# Update dependencies
go get -u ./...
go mod tidy
```

### Docker Container Won't Start

```bash
# View logs
docker logs titan-postgres

# Restart container
docker-compose restart postgres
```

---

## Next Steps

1. Read [Architecture Overview](../architecture/overview.md)
2. Explore [API Documentation](../api/README.md)
3. Review [Cell-Based Architecture](../architecture/cell-architecture.md)
4. Learn [Event Sourcing](../architecture/event-sourcing.md)

---

**Document Version:** 1.0  
**Last Updated:** 2025-12-04
