# Go Code Structure & DDD Patterns

## Domain-Driven Design Implementation Guide

**Purpose**: Define standard code organization, design patterns, and best practices for building Golang microservices in the Titan Commerce Platform.

---

## 📁 Project Structure (Per Service)

```
order-service/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
├── internal/
│   ├── domain/                     # Core business logic
│   │   ├── order/
│   │   │   ├── order.go           # Aggregate root
│   │   │   ├── order_item.go      # Entity
│   │   │   ├── order_status.go    # Value object
│   │   │   └── repository.go      # Interface
│   │   └── events/
│   │       ├── order_created.go
│   │       └── order_paid.go
│   ├── application/                # Use cases/services
│   │   ├── commands/
│   │   │   ├── create_order.go
│   │   │   └── cancel_order.go
│   │   ├── queries/
│   │   │   ├── get_order.go
│   │   │   └── list_orders.go
│   │   └── services/
│   │       └── order_service.go
│   ├── infrastructure/             # External concerns
│   │   ├── persistence/
│   │   │   ├── postgres/
│   │   │   │   └── order_repository.go
│   │   │   └── scylla/
│   │   │       └── event_store.go
│   │   ├── messaging/
│   │   │   └── kafka/
│   │   │       └── producer.go
│   │   └── grpc/
│   │       ├── server.go
│   │       └── interceptors.go
│   └── interfaces/                 # API layer
│       ├── grpc/
│       │   └── order_handler.go
│       └── http/
│           └── rest_handler.go
├── proto/
│   └── order/
│       └── v1/
│           └── order.proto         # gRPC contract
├── config/
│   ├── config.go
│   └── config.yaml
├── pkg/                            # Reusable utilities
│   ├── logger/
│   ├── validator/
│   └── errors/
├── migrations/
│   ├── 001_create_orders.up.sql
│   └── 001_create_orders.down.sql
├── tests/
│   ├── integration/
│   └── e2e/
├── Dockerfile
├── go.mod
└── go.sum
```

---

## 🎯 Layer Responsibilities

### 1. Domain Layer (`internal/domain/`)

**Purpose**: Contains pure business logic, no external dependencies

**Rules**:
- ✅ No database imports
- ✅ No HTTP/gRPC imports
- ✅ Define interfaces, implement in infrastructure
- ✅ Use dependency injection

**Example - Order Aggregate**:

```go
// internal/domain/order/order.go
package order

import (
    "errors"
    "time"
    
    "github.com/google/uuid"
)

// Order is an aggregate root
type Order struct {
    id         uuid.UUID
    userID     uuid.UUID
    items      []OrderItem
    status     OrderStatus
    totalAmount float64
    createdAt  time.Time
    updatedAt  time.Time
    version    int  // For optimistic locking
}

// NewOrder creates a new order (factory method)
func NewOrder(userID uuid.UUID, items []OrderItem) (*Order, error) {
    if len(items) == 0 {
        return nil, errors.New("order must have at least one item")
    }
    
    order := &Order{
        id:        uuid.New(),
        userID:    userID,
        items:     items,
        status:    StatusPending,
        createdAt: time.Now(),
        version:   1,
    }
    
    // Calculate total
    order.calculateTotal()
    
    return order, nil
}

// Confirm transitions order to confirmed state
func (o *Order) Confirm() error {
    if o.status != StatusPending {
        return ErrInvalidStateTransition
    }
    
    o.status = StatusConfirmed
    o.updatedAt = time.Now()
    o.version++
    
    return nil
}

// Ship transitions order to shipped state
func (o *Order) Ship(trackingNumber string) error {
    if o.status != StatusConfirmed {
        return ErrInvalidStateTransition
    }
    
    o.status = StatusShipped
    o.updatedAt = time.Now()
    o.version++
    
    // Emit domain event
    return o.addEvent(&OrderShippedEvent{
        OrderID:        o.id,
        TrackingNumber: trackingNumber,
        ShippedAt:      time.Now(),
    })
}

// Value Objects
type OrderStatus string

const (
    StatusPending   OrderStatus = "PENDING"
    StatusConfirmed OrderStatus = "CONFIRMED"
    StatusShipped   OrderStatus = "SHIPPED"
    StatusDelivered OrderStatus = "DELIVERED"
    StatusCancelled OrderStatus = "CANCELLED"
)

// Repository interface (implemented in infrastructure)
type Repository interface {
    Save(order *Order) error
    FindByID(id uuid.UUID) (*Order, error)
    FindByUserID(userID uuid.UUID) ([]*Order, error)
    Update(order *Order) error
}
```

**Entity vs Value Object**:

```go
// Entity - Has identity, mutable
type OrderItem struct {
    id        uuid.UUID  // Unique identifier
    productID uuid.UUID
    quantity  int
    price     float64
}

// Value Object - No identity, immutable
type Address struct {
    street  string
    city    string
    zipCode string
    country string
}

// Immutability enforced through methods
func (a Address) WithStreet(street string) Address {
    return Address{
        street:  street,
        city:    a.city,
        zipCode: a.zipCode,
        country: a.country,
    }
}
```

---

### 2. Application Layer (`internal/application/`)

**Purpose**: Orchestrate use cases, coordinate domain logic

**Command** (Write operation):

```go
// internal/application/commands/create_order.go
package commands

import (
    "context"
    "fmt"
    
    "github.com/titan/order-service/internal/domain/order"
    "github.com/titan/order-service/internal/domain/events"
)

type CreateOrderCommand struct {
    UserID uuid.UUID
    Items  []OrderItemDTO
}

type CreateOrderHandler struct {
    orderRepo    order.Repository
    eventBus     events.EventBus
    inventorySvc InventoryService  // External service
}

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd *CreateOrderCommand) (*OrderDTO, error) {
    // 1. Convert DTO to domain model
    items := make([]order.OrderItem, len(cmd.Items))
    for i, item := range cmd.Items {
        items[i] = order.NewOrderItem(item.ProductID, item.Quantity, item.Price)
    }
    
    // 2. Create order (domain logic)
    o, err := order.NewOrder(cmd.UserID, items)
    if err != nil {
        return nil, fmt.Errorf("failed to create order: %w", err)
    }
    
    // 3. Reserve inventory (external call)
    if err := h.inventorySvc.ReserveStock(ctx, o.Items()); err != nil {
        return nil, fmt.Errorf("inventory reservation failed: %w", err)
    }
    
    // 4. Persist
    if err := h.orderRepo.Save(o); err != nil {
        // Compensate: release inventory
        h.inventorySvc.ReleaseStock(ctx, o.Items())
        return nil, err
    }
    
    // 5. Publish domain event
    event := &events.OrderCreatedEvent{
        OrderID:    o.ID(),
        UserID:     o.UserID(),
        TotalAmount: o.TotalAmount(),
        CreatedAt:  o.CreatedAt(),
    }
    h.eventBus.Publish(ctx, event)
    
    // 6. Return DTO
    return ToOrderDTO(o), nil
}
```

**Query** (Read operation):

```go
// internal/application/queries/get_order.go
package queries

type GetOrderQuery struct {
    OrderID uuid.UUID
}

type GetOrderHandler struct {
    orderRepo order.Repository
}

func (h *GetOrderHandler) Handle(ctx context.Context, q *GetOrderQuery) (*OrderDTO, error) {
    o, err := h.orderRepo.FindByID(q.OrderID)
    if err != nil {
        return nil, err
    }
    
    return ToOrderDTO(o), nil
}
```

---

### 3. Infrastructure Layer (`internal/infrastructure/`)

**Purpose**: Implement interfaces defined in domain, handle external systems

**Repository Implementation**:

```go
// internal/infrastructure/persistence/postgres/order_repository.go
package postgres

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    
    "github.com/google/uuid"
    "github.com/titan/order-service/internal/domain/order"
)

type OrderRepository struct {
    db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
    return &OrderRepository{db: db}
}

func (r *OrderRepository) Save(o *order.Order) error {
    itemsJSON, _ := json.Marshal(o.Items())
    
    query := `
        INSERT INTO orders (id, user_id, items, status, total_amount, created_at, version)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    
    _, err := r.db.Exec(query,
        o.ID(),
        o.UserID(),
        itemsJSON,
        o.Status(),
        o.TotalAmount(),
        o.CreatedAt(),
        o.Version(),
    )
    
    return err
}

func (r *OrderRepository) FindByID(id uuid.UUID) (*order.Order, error) {
    var (
        userID      uuid.UUID
        itemsJSON   []byte
        status      string
        totalAmount float64
        createdAt   time.Time
        version     int
    )
    
    query := `
        SELECT user_id, items, status, total_amount, created_at, version
        FROM orders
        WHERE id = $1
    `
    
    err := r.db.QueryRow(query, id).Scan(
        &userID, &itemsJSON, &status, &totalAmount, &createdAt, &version,
    )
    
    if err == sql.ErrNoRows {
        return nil, order.ErrNotFound
    }
    if err != nil {
        return nil, err
    }
    
    var items []order.OrderItem
    json.Unmarshal(itemsJSON, &items)
    
    // Reconstruct order (using package-level constructor)
    return order.ReconstructOrder(id, userID, items, order.OrderStatus(status), totalAmount, createdAt, version), nil
}

func (r *OrderRepository) Update(o *order.Order) error {
    // Optimistic locking
    result, err := r.db.Exec(`
        UPDATE orders
        SET status = $1, updated_at = $2, version = version + 1
        WHERE id = $3 AND version = $4
    `, o.Status(), time.Now(), o.ID(), o.Version())
    
    if err != nil {
        return err
    }
    
    rows, _ := result.RowsAffected()
    if rows == 0 {
        return order.ErrVersionMismatch  // Concurrent update detected
    }
    
    return nil
}
```

**Kafka Producer**:

```go
// internal/infrastructure/messaging/kafka/producer.go
package kafka

import (
    "context"
    "encoding/json"
    
    "github.com/segmentio/kafka-go"
)

type EventPublisher struct {
    writer *kafka.Writer
}

func NewEventPublisher(brokers []string) *EventPublisher {
    return &EventPublisher{
        writer: &kafka.Writer{
            Addr:     kafka.TCP(brokers...),
            Topic:    "order-events",
            Balancer: &kafka.LeastBytes{},
        },
    }
}

func (p *EventPublisher) Publish(ctx context.Context, event interface{}) error {
    payload, _ := json.Marshal(event)
    
    return p.writer.WriteMessages(ctx, kafka.Message{
        Key:   []byte(getEventKey(event)),
        Value: payload,
    })
}
```

---

### 4. Interface Layer (`internal/interfaces/`)

**Purpose**: Handle HTTP/gRPC requests, validate input, call application layer

**gRPC Handler**:

```go
// internal/interfaces/grpc/order_handler.go
package grpc

import (
    "context"
    
    pb "github.com/titan/order-service/proto/order/v1"
    "github.com/titan/order-service/internal/application/commands"
    "github.com/titan/order-service/internal/application/queries"
)

type OrderServer struct {
    pb.UnimplementedOrderServiceServer
    
    createOrderHandler *commands.CreateOrderHandler
    getOrderHandler    *queries.GetOrderHandler
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
    // Validate input
    if err := validateCreateOrderRequest(req); err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
    }
    
    // Convert protobuf to command
    cmd := &commands.CreateOrderCommand{
        UserID: uuid.MustParse(req.UserId),
        Items:  convertItems(req.Items),
    }
    
    // Execute command
    orderDTO, err := s.createOrderHandler.Handle(ctx, cmd)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
    }
    
    // Convert DTO to protobuf
    return &pb.CreateOrderResponse{
        Order: toProtoOrder(orderDTO),
    }, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
    query := &queries.GetOrderQuery{
        OrderID: uuid.MustParse(req.OrderId),
    }
    
    orderDTO, err := s.getOrderHandler.Handle(ctx, query)
    if err != nil {
        return nil, handleError(err)
    }
    
    return &pb.GetOrderResponse{
        Order: toProtoOrder(orderDTO),
    }, nil
}
```

---

## 🔌 Dependency Injection (Wire)

**wire.go**:

```go
//go:build wireinject
// +build wireinject

package main

import (
    "database/sql"
    
    "github.com/google/wire"
    "github.com/titan/order-service/internal/application/commands"
    "github.com/titan/order-service/internal/infrastructure/persistence/postgres"
    "github.com/titan/order-service/internal/infrastructure/messaging/kafka"
)

func InitializeServer(db *sql.DB, kafkaBrokers []string) (*grpc.Server, error) {
    wire.Build(
        // Repositories
        postgres.NewOrderRepository,
        wire.Bind(new(order.Repository), new(*postgres.OrderRepository)),
        
        // Event bus
        kafka.NewEventPublisher,
        wire.Bind(new(events.EventBus), new(*kafka.EventPublisher)),
        
        // Handlers
        commands.NewCreateOrderHandler,
        queries.NewGetOrderHandler,
        
        // gRPC server
        grpc.NewOrderServer,
    )
    
    return nil, nil
}
```

**Usage in main.go**:

```go
func main() {
    // Load config
    cfg := loadConfig()
    
    // Create DB connection
    db, _ := sql.Open("postgres", cfg.DatabaseURL)
    
    // Wire dependencies
    server, err := InitializeServer(db, cfg.KafkaBrokers)
    if err != nil {
        log.Fatal(err)
    }
    
    // Start server
    lis, _ := net.Listen("tcp", ":50051")
    server.Serve(lis)
}
```

---

## 🎨 Design Patterns

### 1. Repository Pattern

**Interface**:
```go
type Repository interface {
    Save(entity *Entity) error
    FindByID(id uuid.UUID) (*Entity, error)
    Delete(id uuid.UUID) error
}
```

**Fake for Testing**:
```go
type FakeRepository struct {
    entities map[uuid.UUID]*Entity
}

func (r *FakeRepository) Save(e *Entity) error {
    r.entities[e.ID()] = e
    return nil
}
```

### 2. Factory Pattern

```go
func NewOrder(userID uuid.UUID, items []OrderItem) (*Order, error) {
    // Validation and initialization
    if len(items) == 0 {
        return nil, ErrNoItems
    }
    
    return &Order{
        id:     uuid.New(),
        userID: userID,
        items:  items,
        status: StatusPending,
    }, nil
}
```

### 3. Strategy Pattern (Pricing)

```go
type PricingStrategy interface {
    Calculate(order *Order) float64
}

type StandardPricing struct{}

func (s *StandardPricing) Calculate(o *Order) float64 {
    total := 0.0
    for _, item := range o.Items() {
        total += item.Price() * float64(item.Quantity())
    }
    return total
}

type DiscountPricing struct {
    discountPercent float64
}

func (d *DiscountPricing) Calculate(o *Order) float64 {
    standard := StandardPricing{}.Calculate(o)
    return standard * (1 - d.discountPercent/100)
}
```

### 4. Specification Pattern

```go
type Specification interface {
    IsSatisfiedBy(order *Order) bool
}

type HighValueOrderSpec struct {
    minAmount float64
}

func (s *HighValueOrderSpec) IsSatisfiedBy(o *Order) bool {
    return o.TotalAmount() >= s.minAmount
}

// Composite specifications
type AndSpecification struct {
    specs []Specification
}

func (a *AndSpecification) IsSatisfiedBy(o *Order) bool {
    for _, spec := range a.specs {
        if !spec.IsSatisfiedBy(o) {
            return false
        }
    }
    return true
}
```

---

## 🧪 Testing Patterns

### Unit Test (Domain)

```go
func TestOrder_Confirm(t *testing.T) {
    // Arrange
    items := []OrderItem{
        order.NewOrderItem(uuid.New(), 2, 10.00),
    }
    o, _ := order.NewOrder(uuid.New(), items)
    
    // Act
    err := o.Confirm()
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, order.StatusConfirmed, o.Status())
}

func TestOrder_Confirm_InvalidState(t *testing.T) {
    o := createShippedOrder()
    
    err := o.Confirm()
    
    assert.Error(t, err)
    assert.Equal(t, order.ErrInvalidStateTransition, err)
}
```

### Integration Test (Repository)

```go
func TestOrderRepository_Save(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer db.Close()
    
    repo := postgres.NewOrderRepository(db)
    
    // Create order
    o, _ := order.NewOrder(uuid.New(), []OrderItem{...})
    
    // Act
    err := repo.Save(o)
    
    // Assert
    assert.NoError(t, err)
    
    // Verify in DB
    found, _ := repo.FindByID(o.ID())
    assert.Equal(t, o.ID(), found.ID())
}
```

### Table-Driven Tests

```go
func TestOrderStatus_Transitions(t *testing.T) {
    tests := []struct {
        name        string
        initialStatus OrderStatus
        action      func(*Order) error
        expectError bool
        finalStatus OrderStatus
    }{
        {
            name:          "Pending to Confirmed",
            initialStatus: StatusPending,
            action:        func(o *Order) error { return o.Confirm() },
            expectError:   false,
            finalStatus:   StatusConfirmed,
        },
        {
            name:          "Confirmed to Shipped",
            initialStatus: StatusConfirmed,
            action:        func(o *Order) error { return o.Ship("TRACK123") },
            expectError:   false,
            finalStatus:   StatusShipped,
        },
        {
            name:          "Pending to Shipped (Invalid)",
            initialStatus: StatusPending,
            action:        func(o *Order) error { return o.Ship("TRACK123") },
            expectError:   true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            o := createOrderWithStatus(tt.initialStatus)
            
            err := tt.action(o)
            
            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.finalStatus, o.Status())
            }
        })
    }
}
```

---

## 📏 Best Practices

### 1. Error Handling

```go
// Define domain errors
var (
    ErrNotFound    = errors.New("order not found")
    ErrInvalidState = errors.New("invalid state transition")
)

// Wrap errors with context
func (r *Repository) FindByID(id uuid.UUID) (*Order, error) {
    o, err := r.db.Query(...)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("order %s: %w", id, ErrNotFound)
    }
    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }
    return o, nil
}

// Check error types
if errors.Is(err, ErrNotFound) {
    return status.Error(codes.NotFound, "order not found")
}
```

### 2. Logging

```go
import "go.uber.org/zap"

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd *CreateOrderCommand) (*OrderDTO, error) {
    logger := zap.L().With(
        zap.String("user_id", cmd.UserID.String()),
        zap.Int("item_count", len(cmd.Items)),
    )
    
    logger.Info("creating order")
    
    o, err := order.NewOrder(cmd.UserID, items)
    if err != nil {
        logger.Error("order creation failed", zap.Error(err))
        return nil, err
    }
    
    logger.Info("order created", zap.String("order_id", o.ID().String()))
    return ToOrderDTO(o), nil
}
```

### 3. Context Propagation

```go
func (h *Handler) Handle(ctx context.Context, cmd *Command) error {
    // Extract trace ID from context
    traceID := ctx.Value("trace_id").(string)
    
    // Propagate to downstream services
    ctx = metadata.AppendToOutgoingContext(ctx, "trace-id", traceID)
    
    // Call external service
    resp, err := h.client.SomeMethod(ctx, req)
    
    return err
}
```

---

**Document Version**: 1.0  
**Last Updated**: 2025-12-04  
**Pages**: 50+ (complete guide)
