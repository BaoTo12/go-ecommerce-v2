# Order Service

Production-ready order management microservice with DDD, Event Sourcing, and CQRS.

## Architecture

- **Domain Layer**: Business logic, aggregates, value objects
- **Application Layer**: Use cases, commands, queries  
- **Infrastructure Layer**: Database, messaging, external services
- **Interface Layer**: gRPC, HTTP REST APIs

## Features

- ✅ Event Sourcing for complete audit trail
- ✅ CQRS (Command Query Responsibility Segregation)
- ✅ Domain-Driven Design (DDD) patterns
- ✅ PostgreSQL for persistence
- ✅ Kafka for event publishing
- ✅ gRPC API with Protocol Buffers
- ✅ Optimistic locking for concurrency

## Quick Start

```bash
# Set environment variables
export SERVICE_NAME=order-service
export CELL_ID=cell-001
export DATABASE_URL=postgresql://user:pass@localhost:5432/orders
export KAFKA_BROKERS=localhost:9092

# Run the service
go run cmd/server/main.go
```

## API (gRPC)

See `proto/order/v1/order.proto` for complete API definition.

### Create Order
```bash
grpcurl -plaintext -d '{"user_id":"user-123","items":[...],"shipping_address":"..."}' localhost:9000 order.v1.OrderService/CreateOrder
```

## Database Schema

```sql
CREATE TABLE orders (
  id UUID PRIMARY KEY,
  user_id VARCHAR(255) NOT NULL,
  total_amount DECIMAL(10,2) NOT NULL,
  status VARCHAR(50) NOT NULL,
  shipping_address TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE order_items (
  id UUID PRIMARY KEY,
  order_id UUID REFERENCES orders(id),
  product_id VARCHAR(255) NOT NULL,
  product_name VARCHAR(255) NOT NULL,
  quantity INT NOT NULL,
  unit_price DECIMAL(10,2) NOT NULL,
  subtotal DECIMAL(10,2) NOT NULL
);

CREATE TABLE order_events (
  id UUID PRIMARY KEY,
  aggregate_id UUID NOT NULL,
  event_type VARCHAR(100) NOT NULL,
  data JSONB NOT NULL,
  user_id VARCHAR(255) NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  version INT NOT NULL
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_order_events_aggregate_id ON order_events(aggregate_id);
```

## Testing

```bash
# Unit tests  
go test ./internal/domain/...

# Integration tests
go test -tags=integration ./...
```
