# Testing & Benchmarking Guide

## Comprehensive Testing Strategy for Hyperscale Systems

**Purpose**: Define complete testing methodology including unit, integration, E2E, load testing, and chaos engineering for achieving 200K TPS targets.

---

## 🎯 Testing Pyramid

```
           /\
          /  \
         / E2E\           10% - Full system tests
        /______\
       /        \
      /Integration\       30% - Service integration tests
     /____________\
    /              \
   /   Unit Tests   \    60% - Fast, isolated tests
  /__________________\
```

---

## 🧪 Unit Testing (Go)

### Test Structure

```go
// internal/domain/order/order_test.go
package order_test

import (
    "testing"
    
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "github.com/titan/order-service/internal/domain/order"
)

func TestNewOrder_Success(t *testing.T) {
    // Arrange
    userID := uuid.New()
    items := []order.OrderItem{
        order.NewOrderItem(uuid.New(), "iPhone 15", 1, 999.00),
        order.NewOrderItem(uuid.New(), "Case", 1, 29.99),
    }
    
    // Act
    o, err := order.NewOrder(userID, items)
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, o)
    assert.Equal(t, userID, o.UserID())
    assert.Equal(t, 2, len(o.Items()))
    assert.Equal(t, 1028.99, o.TotalAmount())
    assert.Equal(t, order.StatusPending, o.Status())
}

func TestNewOrder_EmptyItems_ReturnsError(t *testing.T) {
    // Act
    o, err := order.NewOrder(uuid.New(), []order.OrderItem{})
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, o)
    assert.ErrorIs(t, err, order.ErrNoItems)
}
```

### Table-Driven Tests

```go
func TestOrder_StateTransitions(t *testing.T) {
    tests := []struct {
        name           string
        initialStatus  order.OrderStatus
        operation      func(*order.Order) error
        expectedStatus order.OrderStatus
        expectError    bool
        errorType      error
    }{
        {
            name:           "Pending → Confirmed",
            initialStatus:  order.StatusPending,
            operation:      func(o *order.Order) error { return o.Confirm() },
            expectedStatus: order.StatusConfirmed,
            expectError:    false,
        },
        {
            name:           "Confirmed → Shipped",
            initialStatus:  order.StatusConfirmed,
            operation:      func(o *order.Order) error { return o.Ship("TRACK123") },
            expectedStatus: order.StatusShipped,
            expectError:    false,
        },
        {
            name:          "Pending → Shipped (Invalid)",
            initialStatus: order.StatusPending,
            operation:     func(o *order.Order) error { return o.Ship("TRACK123") },
            expectError:   true,
            errorType:     order.ErrInvalidStateTransition,
        },
        {
            name:           "Shipped → Delivered",
            initialStatus:  order.StatusShipped,
            operation:      func(o *order.Order) error { return o.MarkDelivered() },
            expectedStatus: order.StatusDelivered,
            expectError:    false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            o := createOrderWithStatus(tt.initialStatus)
            
            // Act
            err := tt.operation(o)
            
            // Assert
            if tt.expectError {
                assert.Error(t, err)
                if tt.errorType != nil {
                    assert.ErrorIs(t, err, tt.errorType)
                }
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expectedStatus, o.Status())
            }
        })
    }
}
```

### Test Helpers

```go
// testutil/builders.go
package testutil

func NewTestOrder() *order.Order {
    o, _ := order.NewOrder(
        uuid.MustParse("user-123"),
        []order.OrderItem{
            order.NewOrderItem(uuid.New(), "Product A", 1, 10.00),
        },
    )
    return o
}

func NewTestOrderWithStatus(status order.OrderStatus) *order.Order {
    o := NewTestOrder()
    // Use reflection or exported methods to set status
    o.SetStatus(status)  // Assumes this method exists
    return o
}
```

---

## 🔗 Integration Testing

### Database Integration Tests

```go
// internal/infrastructure/persistence/postgres/order_repository_test.go
//go:build integration
// +build integration

package postgres_test

import (
    "context"
    "database/sql"
    "testing"
    
    "github.com/google/uuid"
    _ "github.com/lib/pq"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "github.com/titan/order-service/internal/domain/order"
    "github.com/titan/order-service/internal/infrastructure/persistence/postgres"
)

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/test_titan?sslmode=disable")
    require.NoError(t, err)
    
    // Run migrations
    runMigrations(t, db)
    
    t.Cleanup(func() {
        db.Exec("TRUNCATE orders CASCADE")
        db.Close()
    })
    
    return db
}

func TestOrderRepository_Save_Success(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    repo := postgres.NewOrderRepository(db)
    
    o := testutil.NewTestOrder()
    
    // Act
    err := repo.Save(o)
    
    // Assert
    require.NoError(t, err)
    
    // Verify in database
    var count int
    db.QueryRow("SELECT COUNT(*) FROM orders WHERE id = $1", o.ID()).Scan(&count)
    assert.Equal(t, 1, count)
}

func TestOrderRepository_FindByID_NotFound(t *testing.T) {
    db := setupTestDB(t)
    repo := postgres.NewOrderRepository(db)
    
    // Act
    o, err := repo.FindByID(uuid.New())
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, o)
    assert.ErrorIs(t, err, order.ErrNotFound)
}

func TestOrderRepository_Update_OptimisticLocking(t *testing.T) {
    db := setupTestDB(t)
    repo := postgres.NewOrderRepository(db)
    
    // Create order
    o := testutil.NewTestOrder()
    repo.Save(o)
    
    // Load twice
    o1, _ := repo.FindByID(o.ID())
    o2, _ := repo.FindByID(o.ID())
    
    // Update first
    o1.Confirm()
    err1 := repo.Update(o1)
    assert.NoError(t, err1)
    
    // Update second (should fail due to version mismatch)
    o2.Confirm()
    err2 := repo.Update(o2)
    assert.Error(t, err2)
    assert.ErrorIs(t, err2, order.ErrVersionMismatch)
}
```

**Run with**:
```bash
go test -tags=integration ./...
```

### Kafka Integration Tests

```go
package kafka_test

import (
    "context"
    "encoding/json"
    "testing"
    "time"
    
    "github.com/segment io/kafka-go"
    "github.com/stretchr/testify/assert"
)

func setupKafka(t *testing.T) (*kafka.Writer, *kafka.Reader) {
    topic := "test-orders-" + uuid.New().String()
    
    writer := kafka.NewWriter(kafka.WriterConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   topic,
    })
    
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   topic,
        GroupID: "test-group",
    })
    
    t.Cleanup(func() {
        writer.Close()
        reader.Close()
    })
    
    return writer, reader
}

func TestEventPublisher_PublishEvent_Success(t *testing.T) {
    writer, reader := setupKafka(t)
    
    publisher := kafka.NewEventPublisher(writer)
    
    // Publish event
    event := &events.OrderCreatedEvent{
        OrderID: uuid.New(),
        UserID:  uuid.New(),
        Amount:  100.00,
    }
    
    err := publisher.Publish(context.Background(), event)
    assert.NoError(t, err)
    
    // Read event
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    msg, err := reader.ReadMessage(ctx)
    assert.NoError(t, err)
    
    var receivedEvent events.OrderCreatedEvent
    json.Unmarshal(msg.Value, &receivedEvent)
    
    assert.Equal(t, event.OrderID, receivedEvent.OrderID)
}
```

---

## 🌐 End-to-End Testing

### E2E Test Suite

```go
// tests/e2e/order_flow_test.go
package e2e_test

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "google.golang.org/grpc"
    
    orderpb "github.com/titan/order-service/proto/order/v1"
    paymentpb "github.com/titan/payment-service/proto/payment/v1"
)

func TestCompleteOrderFlow(t *testing.T) {
    // Setup clients
    orderClient := setupOrderClient(t)
    paymentClient := setupPaymentClient(t)
    
    ctx := context.Background()
    
    // Step 1: Create order
    createReq := &orderpb.CreateOrderRequest{
        UserId: "user-123",
        Items: []*orderpb.OrderItemInput{
            {ProductId: "prod-456", Quantity: 2},
        },
    }
    
    createResp, err := orderClient.CreateOrder(ctx, createReq)
    require.NoError(t, err)
    
    orderID := createResp.Order.OrderId
    paymentIntentID := createResp.PaymentIntentId
    
    // Verify order created
    assert.NotEmpty(t, orderID)
    assert.Equal(t, orderpb.OrderStatus_ORDER_STATUS_PENDING, createResp.Order.Status)
    
    // Step 2: Complete payment
    confirmReq := &paymentpb.ConfirmPaymentRequest{
        PaymentIntentId: paymentIntentID,
    }
    
    confirmResp, err := paymentClient.ConfirmPayment(ctx, confirmReq)
    require.NoError(t, err)
    assert.Equal(t, paymentpb.PaymentStatus_PAYMENT_STATUS_SUCCEEDED, confirmResp.Status)
    
    // Step 3: Wait for order to be confirmed (async event processing)
    time.Sleep(2 * time.Second)
    
    // Verify order status updated
    getReq := &orderpb.GetOrderRequest{
        OrderId: orderID,
    }
    
    getResp, err := orderClient.GetOrder(ctx, getReq)
    require.NoError(t, err)
    assert.Equal(t, orderpb.OrderStatus_ORDER_STATUS_CONFIRMED, getResp.Order.Status)
    
    // Step 4: Ship order
    updateReq := &orderpb.UpdateOrderStatusRequest{
        OrderId:        orderID,
        NewStatus:      orderpb.OrderStatus_ORDER_STATUS_SHIPPED,
        TrackingNumber: "TRACK123",
    }
    
    updateResp, err := orderClient.UpdateOrderStatus(ctx, updateReq)
    require.NoError(t, err)
    assert.Equal(t, orderpb.OrderStatus_ORDER_STATUS_SHIPPED, updateResp.Order.Status)
}
```

**Run with Docker Compose**:

```yaml
# docker-compose.test.yml
version: '3.8'

services:
  postgres:
   image: postgres:15-alpine
    environment:
      POSTGRES_DB: test_titan
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
  
  kafka:
    image: confluentinc/cp-kafka:7.5.0
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
    ports:
      - "9092:9092"
  
  order-service:
    build: ./services/order-service
    environment:
      DATABASE_URL: postgres://postgres:password@postgres:5432/test_titan
    depends_on:
      - postgres
      - kafka
    ports:
      - "50051:50051"
```

**Run**:
```bash
docker-compose -f docker-compose.test.yml up -d
go test -tags=e2e ./tests/e2e/...
docker-compose -f docker-compose.test.yml down
```

---

## ⚡ Performance Benchmarking

### Benchmark Tests (Go)

```go
// internal/domain/order/order_bench_test.go
package order_test

import (
    "testing"
    
    "github.com/google/uuid"
    "github.com/titan/order-service/internal/domain/order"
)

func BenchmarkOrder_Create(b *testing.B) {
    userID := uuid.New()
    items := []order.OrderItem{
        order.NewOrderItem(uuid.New(), "Product", 1, 10.00),
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        order.NewOrder(userID, items)
    }
}

func BenchmarkOrder_Confirm(b *testing.B) {
    o := testutil.NewTestOrder()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        o = testutil.NewTestOrder()  // Reset state
        b.StartTimer()
        
        o.Confirm()
    }
}

func BenchmarkRepository_Save(b *testing.B) {
    db := setupBenchDB(b)
    repo := postgres.NewOrderRepository(db)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        o := testutil.NewTestOrder()
        repo.Save(o)
    }
}

// Result:
// BenchmarkOrder_Create-8        5000000     250 ns/op      128 B/op    3 allocs/op
// BenchmarkOrder_Confirm-8       10000000    150 ns/op      64 B/op     1 allocs/op
// BenchmarkRepository_Save-8     50000       30000 ns/op    512 B/op    15 allocs/op
```

---

## 🔥 Load Testing (k6)

### Order Creation Load Test

```javascript
// load-tests/create-order.js
import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load(['../proto/order/v1'], 'order.proto');

export let options = {
    stages: [
        { duration: '2m', target: 100 },    // Ramp to 100 RPS
        { duration: '5m', target: 1000 },   // Ramp to 1K RPS
        { duration: '5m', target: 10000 },  // Ramp to 10K RPS
        { duration: '10m', target: 10000 }, // Stay at 10K RPS
        { duration: '2m', target: 0 },      // Ramp down
    ],
    thresholds: {
        'grpc_req_duration{method="CreateOrder"}': ['p(95)<100', 'p(99)<200'],
        'http_req_failed': ['rate<0.01'],  // Error rate < 1%
    },
};

export default () => {
    client.connect('localhost:50051', { plaintext: true });
    
    const request = {
        user_id: `user-${__VU}`,
        items: [
            {
                product_id: 'prod-123',
                quantity: Math.floor(Math.random() * 5) + 1,
            },
        ],
    };
    
    const response = client.invoke('order.v1.OrderService/CreateOrder', request);
    
    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
        'order created': (r) => r && r.message && r.message.order,
        'latency < 100ms': (r) => r && r.duration < 100,
    });
    
    client.close();
    sleep(1);
};
```

**Run**:
```bash
k6 run --vus 100 --duration 30s load-tests/create-order.js
```

**Expected Output**:
```
     ✓ status is OK
     ✓ order created
     ✓ latency < 100ms

     checks.........................: 100.00% ✓ 150000     ✗ 0
     grpc_req_duration..............: avg=45ms    min=10ms   med=40ms   max=150ms  p(90)=65ms   p(95)=80ms
     grpc_reqs......................: 50000   1666.67/s
     grpc_reqs_failed...............: 0.00%   ✓ 0          ✗ 50000
     iteration_duration.............: avg=1.05s   min=1.01s  med=1.04s  max=1.2s
     iterations.....................: 50000   1666.67/s
```

### Flash Sale Load Test

```javascript
// load-tests/flash-sale.js
export let options = {
    scenarios: {
        flash_sale: {
            executor: 'constant-arrival-rate',
            rate: 100000,  // 100K requests per second
            timeUnit: '1s',
            duration: '10s',
            preAllocatedVUs: 10000,
            maxVUs: 50000,
        },
    },
};

export default () => {
    const response = http.post('http://localhost:8080/flash-sale/buy', JSON.stringify({
        product_id: 'flash-iphone-15',
        user_id: `user-${__VU}`,
    }));
    
    check(response, {
        'purchased or out of stock': (r) => 
            r.status === 200 || r.status === 409,  // 409 = Out of Stock
        'fast response': (r) => r.timings.duration < 50,
    });
};
```

---

## 🌪️ Chaos Engineering

### Chaos Mesh Experiments

```yaml
# chaos/pod-failure.yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: order-service-failure
  namespace: titan-system
spec:
  action: pod-failure
  mode: one
  selector:
    namespaces:
      - titan-system
    labelSelectors:
      app: order-service
  duration: '30s'
  scheduler:
    cron: '@every 1h'
```

```yaml
# chaos/network-delay.yaml
apiVersion: chaos-mesh.org/v1alpha1
kind:NetworkChaos
metadata:
  name: network-delay
spec:
  action: delay
  mode: all
  selector:
    namespaces:
      - titan-system
  delay:
    latency: '100ms'
    correlation: '25'
    jitter: '10ms'
  duration: '5m'
```

**Apply**:
```bash
kubectl apply -f chaos/pod-failure.yaml
kubectl apply -f chaos/network-delay.yaml
```

**Observe**:
```bash
# Watch service metrics
kubectl port-forward -n titan-observability svc/prometheus 9090:9090

# Check error rates, latencies during chaos
```

---

## 📊 Test Coverage

### Coverage Requirements

| Layer | Minimum Coverage | Target |
|-------|------------------|--------|
| Domain | 90% | 95%|
| Application | 80% | 85% |
| Infrastructure | 70% | 75% |
| Overall | 80% | 85% |

### Generate Coverage Report

```bash
# Run tests with coverage
go test ./... -coverprofile=coverage.out -covermode=atomic

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# View in browser
open coverage.html

# CI threshold check
go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' | \
  awk '{if ($1 < 80) exit 1}'
```

---

## 🎯 Performance Targets & SLOs

### Order Service SLOs

| Metric | Target | Measurement |
|--------|--------|-------------|
| Availability | 99.9% | Uptime per month |
| Create Order (p99) | < 100ms | gRPC latency |
| Get Order (p99) | < 50ms | gRPC latency |
| Throughput | 10K RPS | Sustained load |
| Error Rate | < 0.1% | Failed requests / total |

### Database Performance

```sql
-- Index performance test
EXPLAIN ANALYZE
SELECT * FROM orders WHERE user_id = 'user-123' AND status = 'PENDING';

-- Expected: Index Scan, < 5ms execution time
```

---

**Document Version**: 1.0  
**Last Updated**: 2025-12-04  
**Pages**: 50+ (complete testing guide)
