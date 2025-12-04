# Complete API Reference - gRPC & REST

## Protocol Buffer Definitions & Service Contracts

**Purpose**: Complete API specifications for all microservices with Protocol Buffer definitions, gRPC service contracts, and REST endpoint mappings.

---

## 📋 Service Catalog

| Service | Port | Purpose | Protocol |
|---------|------|---------|----------|
| Order Service | 5001 | Order management | gRPC + REST |
| Inventory Service | 5002 | Stock management | gRPC |
| Payment Service | 5003 | Payment processing | gRPC |
| User Service | 5004 | User profiles | gRPC + REST |
| Product Service | 5005 | Product catalog | gRPC + REST |
| Chat Service | 5006 | Real-time messaging | gRPC + WebSocket |
| Notification Service | 5007 | Push notifications | gRPC |
| Search Service | 5008 | Product search | gRPC + REST |

---

## 🔧 Order Service API

### Protocol Buffer Definition

```protobuf
// proto/order/v1/order.proto
syntax = "proto3";

package order.v1;

option go_package = "github.com/titan/order-service/proto/order/v1;orderpb";

import "google/protobuf/timestamp.proto";
import "google/api/annotations.proto";

service OrderService {
  // Create a new order
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse) {
    option (google.api.http) = {
      post: "/v1/orders"
      body: "*"
    };
  }
  
  // Get order by ID
  rpc GetOrder(GetOrderRequest) returns (GetOrderResponse) {
    option (google.api.http) = {
      get: "/v1/orders/{order_id}"
    };
  }
  
  // List orders for user
  rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse) {
    option (google.api.http) = {
      get: "/v1/users/{user_id}/orders"
    };
  }
  
  // Update order status
  rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (UpdateOrderStatusResponse) {
    option (google.api.http) = {
      patch: "/v1/orders/{order_id}/status"
      body: "*"
    };
  }
  
  // Cancel order
  rpc CancelOrder(CancelOrderRequest) returns (CancelOrderResponse) {
    option (google.api.http) = {
      post: "/v1/orders/{order_id}/cancel"
      body: "*"
    };
  }
  
  // Stream order updates
  rpc StreamOrderUpdates(StreamOrderUpdatesRequest) returns (stream OrderUpdate);
}

// Messages
message Order {
  string order_id = 1;
  string user_id = 2;
  repeated OrderItem items = 3;
  OrderStatus status = 4;
  double total_amount = 5;
  string currency = 6;
  Address shipping_address = 7;
  google.protobuf.Timestamp created_at = 8;
  google.protobuf.Timestamp updated_at = 9;
  int32 version = 10;  // For optimistic locking
}

message OrderItem {
  string item_id = 1;
  string product_id = 2;
  string product_name = 3;
  int32 quantity = 4;
  double unit_price = 5;
  double total_price = 6;
  string seller_id = 7;
}

enum OrderStatus {
  ORDER_STATUS_UNSPECIFIED = 0;
  ORDER_STATUS_PENDING = 1;
  ORDER_STATUS_CONFIRMED = 2;
  ORDER_STATUS_SHIPPED = 3;
  ORDER_STATUS_DELIVERED = 4;
  ORDER_STATUS_CANCELLED = 5;
  ORDER_STATUS_REFUNDED = 6;
}

message Address {
  string street = 1;
  string city = 2;
  string state = 3;
  string zip_code = 4;
  string country = 5;
}

message CreateOrderRequest {
  string user_id = 1;
  repeated OrderItemInput items = 2;
  Address shipping_address = 3;
  string payment_method_id = 4;
  repeated string voucher_codes = 5;
}

message OrderItemInput {
  string product_id = 1;
  int32 quantity = 2;
}

message CreateOrderResponse {
  Order order = 1;
  string payment_intent_id = 2;
}

message GetOrderRequest {
  string order_id = 1;
}

message GetOrderResponse {
  Order order = 1;
}

message ListOrdersRequest {
  string user_id = 1;
  int32 page_size = 2;
  string page_token = 3;
  OrderStatus status_filter = 4;
}

message ListOrdersResponse {
  repeated Order orders = 1;
  string next_page_token = 2;
  int32 total_count = 3;
}

message UpdateOrderStatusRequest {
  string order_id = 1;
  OrderStatus new_status = 2;
  string tracking_number = 3;  // For shipped status
}

message UpdateOrderStatusResponse {
  Order order = 1;
}

message CancelOrderRequest {
  string order_id = 1;
  string reason = 2;
}

message CancelOrderResponse {
  Order order = 1;
  bool refund_initiated = 2;
}

message StreamOrderUpdatesRequest {
  string user_id = 1;
}

message OrderUpdate {
  string order_id = 1;
  OrderStatus  new_status = 2;
  google.protobuf.Timestamp updated_at = 3;
  string message = 4;
}
```

### REST API Mapping (via gRPC-Gateway)

**Create Order**:
```http
POST /v1/orders
Content-Type: application/json

{
  "user_id": "user-123",
  "items": [
    {
      "product_id": "prod-456",
      "quantity": 2
    }
  ],
  "shipping_address": {
    "street": "123 Main St",
    "city": "San Francisco",
    "state": "CA",
    "zip_code": "94102",
    "country": "US"
  },
  "payment_method_id": "pm_xxx",
  "voucher_codes": ["SAVE10"]
}
```

**Response**:
```json
{
  "order": {
    "order_id": "order-789",
    "user_id": "user-123",
    "items": [
      {
        "item_id": "item-001",
        "product_id": "prod-456",
        "product_name": "iPhone 15",
        "quantity": 2,
        "unit_price": 999.00,
        "total_price": 1998.00,
        "seller_id": "seller-x"
      }
    ],
    "status": "ORDER_STATUS_PENDING",
    "total_amount": 1798.20,
    "currency": "USD",
    "created_at": "2025-12-04T00:00:00Z",
    "version": 1
  },
  "payment_intent_id": "pi_xxx"
}
```

**Get Order**:
```http
GET /v1/orders/order-789
```

**List Orders**:
```http
GET /v1/users/user-123/orders?page_size=20&status_filter=ORDER_STATUS_DELIVERED
```

---

## 💰 Payment Service API

```protobuf
// proto/payment/v1/payment.proto
syntax = "proto3";

package payment.v1;

option go_package = "github.com/titan/payment-service/proto/payment/v1;paymentpb";

service PaymentService {
  rpc CreatePaymentIntent(CreatePaymentIntentRequest) returns (CreatePaymentIntentResponse);
  rpc ConfirmPayment(ConfirmPaymentRequest) returns (ConfirmPaymentResponse);
  rpc RefundPayment(RefundPaymentRequest) returns (RefundPaymentResponse);
  rpc GetPaymentStatus(GetPaymentStatusRequest) returns (GetPaymentStatusResponse);
}

message CreatePaymentIntentRequest {
  string order_id = 1;
  double amount = 2;
  string currency = 3;
  string payment_method_id = 4;
  string idempotency_key = 5;
  map<string, string> metadata = 6;
}

message CreatePaymentIntentResponse {
  string payment_intent_id = 1;
  string client_secret = 2;
  PaymentStatus status = 3;
}

message ConfirmPaymentRequest {
  string payment_intent_id = 1;
}

message ConfirmPaymentResponse {
  PaymentStatus status = 1;
  string transaction_id = 2;
}

message RefundPaymentRequest {
  string payment_intent_id = 1;
  double amount = 2;  // Partial refund if less than total
  string reason = 3;
}

message RefundPaymentResponse {
  string refund_id = 1;
  RefundStatus status = 2;
}

enum PaymentStatus {
  PAYMENT_STATUS_UNSPECIFIED = 0;
  PAYMENT_STATUS_PENDING = 1;
  PAYMENT_STATUS_PROCESSING = 2;
  PAYMENT_STATUS_SUCCEEDED = 3;
  PAYMENT_STATUS_FAILED = 4;
  PAYMENT_STATUS_CANCELLED = 5;
}

enum RefundStatus {
  REFUND_STATUS_UNSPECIFIED = 0;
  REFUND_STATUS_PENDING = 1;
  REFUND_STATUS_SUCCEEDED = 2;
  REFUND_STATUS_FAILED = 3;
}
```

---

## 📦 Inventory Service API

```protobuf
// proto/inventory/v1/inventory.proto
syntax = "proto3";

package inventory.v1;

service InventoryService {
  // Check stock availability
  rpc CheckStock(CheckStockRequest) returns (CheckStockResponse);
  
  // Reserve stock (atomic operation)
  rpc ReserveStock(ReserveStockRequest) returns (ReserveStockResponse);
  
  // Release reserved stock (compensation)
  rpc ReleaseStock(ReleaseStockRequest) returns (ReleaseStockResponse);
  
  // Commit reservation (after payment)
  rpc CommitReservation(CommitReservationRequest) returns (CommitReservationResponse);
  
  // Update stock (seller operation)
  rpc UpdateStock(UpdateStockRequest) returns (UpdateStockResponse);
  
  // Bulk stock check
  rpc BulkCheckStock(BulkCheckStockRequest) returns (BulkCheckStockResponse);
}

message CheckStockRequest {
  string product_id = 1;
  int32 quantity = 2;
}

message CheckStockResponse {
  bool available = 1;
  int32 current_stock = 2;
  int32 reserved_stock = 3;
}

message ReserveStockRequest {
  string product_id = 1;
  int32 quantity = 2;
  string order_id = 3;
  int32 ttl_seconds = 4;  // Reservation expires after
}

message ReserveStockResponse {
  bool success = 1;
  string reservation_id = 2;
  string error_message = 3;
}

message ReleaseStockRequest {
  string reservation_id = 1;
}

message ReleaseStockResponse {
  bool success = 1;
}

message CommitReservationRequest {
  string reservation_id = 1;
}

message CommitReservationResponse {
  bool success = 1;
}
```

**Implementation Note - Redis Lua Script**:
```lua
-- Redis script for atomic stock reservation
local product_id = KEYS[1]
local quantity = tonumber(ARGV[1])
local order_id = ARGV[2]
local ttl = tonumber(ARGV[3])

local stock_key = "stock:" .. product_id
local reserved_key = "reserved:" .. product_id
local reservation_key = "reservation:" .. order_id

-- Get current stock
local current_stock = tonumber(redis.call('GET', stock_key) or 0)
local reserved = tonumber(redis.call('GET', reserved_key) or 0)

-- Check availability
local available = current_stock - reserved

if available < quantity then
    return {false, "Insufficient stock", available}
end

-- Reserve stock
redis.call('INCRBY', reserved_key, quantity)
redis.call('SET', reservation_key, quantity, 'EX', ttl)

return {true, reservation_key, available - quantity}
```

---

## 🔍 Search Service API

```protobuf
// proto/search/v1/search.proto
syntax = "proto3";

package search.v1;

service SearchService {
  rpc SearchProducts(SearchProductsRequest) returns (SearchProductsResponse);
  rpc Autocomplete(AutocompleteRequest) returns (AutocompleteResponse);
  rpc IndexProduct(IndexProductRequest) returns (IndexProductResponse);
}

message SearchProductsRequest {
  string query = 1;
  repeated Filter filters = 2;
  SortOption sort = 3;
  int32 page = 4;
  int32 page_size = 5;
}

message Filter {
  string field = 1;
  FilterOperator operator = 2;
  repeated string values = 3;
}

enum FilterOperator {
  FILTER_OPERATOR_UNSPECIFIED = 0;
  FILTER_OPERATOR_EQUALS = 1;
  FILTER_OPERATOR_IN = 2;
  FILTER_OPERATOR_RANGE = 3;
}

message SortOption {
  string field = 1;
  SortOrder order = 2;
}

enum SortOrder {
  SORT_ORDER_UNSPECIFIED = 0;
  SORT_ORDER_ASC = 1;
  SORT_ORDER_DESC = 2;
}

message SearchProductsResponse {
  repeated ProductSearchResult products = 1;
  int32 total_count = 2;
  int32 page = 3;
  repeated Facet facets = 4;
}

message ProductSearchResult {
  string product_id = 1;
  string name = 2;
  string description = 3;
  double price = 4;
  string image_url = 5;
  double relevance_score = 6;
}

message Facet {
  string field = 1;
  repeated FacetValue values = 2;
}

message FacetValue {
  string value = 1;
  int32 count = 2;
}
```

**Elasticsearch Query Example**:
```json
{
  "query": {
    "bool": {
      "must": [
        {
          "multi_match": {
            "query": "iPhone 15",
            "fields": ["name^3", "description", "tags"],
            "type": "best_fields"
          }
        }
      ],
      "filter": [
        {
          "range": {
            "price": {
              "gte": 500,
              "lte": 1500
            }
          }
        },
        {
          "term": {
            "category": "electronics"
          }
        }
      ]
    }
  },
  "aggs": {
    "price_ranges": {
      "range": {
        "field": "price",
        "ranges": [
          { "to": 500 },
          { "from": 500, "to": 1000 },
          { "from": 1000 }
        ]
      }
    },
    "brands": {
      "terms": {
        "field": "brand",
        "size": 10
      }
    }
  },
  "sort": [
    { "_score": "desc" },
    { "price": "asc" }
  ],
  "from": 0,
  "size": 20
}
```

---

## 💬 Chat Service API

```protobuf
// proto/chat/v1/chat.proto
syntax = "proto3";

package chat.v1;

service ChatService {
  // Send message
  rpc SendMessage(SendMessageRequest) returns (SendMessageResponse);
  
  // Get conversation history
  rpc GetConversation(GetConversationRequest) returns (GetConversationResponse);
  
  // Mark messages as read
  rpc MarkAsRead(MarkAsReadRequest) returns (MarkAsReadResponse);
  
  // Stream messages (bidirectional)
  rpc StreamMessages(stream StreamMessagesRequest) returns (stream StreamMessagesResponse);
}

message SendMessageRequest {
  string conversation_id = 1;
  string from_user_id = 2;
  string to_user_id = 3;
  string content = 4;
  MessageType type = 5;
  bytes media_data = 6;  // For images/files
}

enum MessageType {
  MESSAGE_TYPE_UNSPECIFIED = 0;
  MESSAGE_TYPE_TEXT = 1;
  MESSAGE_TYPE_IMAGE = 2;
  MESSAGE_TYPE_FILE = 3;
}

message SendMessageResponse {
  string message_id = 1;
  google.protobuf.Timestamp sent_at = 2;
}

message GetConversationRequest {
  string conversation_id = 1;
  int32 page_size = 2;
  string cursor = 3;  // Timestamp-based pagination
}

message GetConversationResponse {
  repeated Message messages = 1;
  string next_cursor = 2;
}

message Message {
  string message_id = 1;
  string conversation_id = 2;
  string from_user_id = 3;
  string to_user_id = 4;
  string content = 5;
  MessageType type = 6;
  google.protobuf.Timestamp sent_at = 7;
  google.protobuf.Timestamp delivered_at = 8;
  google.protobuf.Timestamp read_at = 9;
}
```

---

## 🔐 Authentication & Authorization

### JWT Structure

```json
{
  "header": {
    "alg": "RS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "user-123",
    "email": "user@example.com",
    "roles": ["buyer", "seller"],
    "permissions": ["order:create", "order:read"],
    "iat": 1701648000,
    "exp": 1701655200
  }
}
```

### gRPC Interceptor

```go
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "missing metadata")
    }
    
    tokens := md.Get("authorization")
    if len(tokens) == 0 {
        return nil, status.Error(codes.Unauthenticated, "missing token")
    }
    
    token := strings.TrimPrefix(tokens[0], "Bearer ")
    
    claims, err := ValidateJWT(token)
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid token")
    }
    
    // Add user to context
    ctx = context.WithValue(ctx, "user_id", claims.Subject)
    ctx = context.WithValue(ctx, "permissions", claims.Permissions)
    
    return handler(ctx, req)
}
```

---

## 📊 Error Handling Standards

### gRPC Status Codes

| Scenario | gRPC Code | HTTP Code |
|----------|-----------|-----------|
| Success | OK | 200 |
| Resource not found | NOT_FOUND | 404 |
| Invalid input | INVALID_ARGUMENT | 400 |
| Authentication failed | UNAUTHENTICATED | 401 |
| Permission denied | PERMISSION_DENIED | 403 |
| Resource already exists | ALREADY_EXISTS | 409 |
| Resource exhausted (rate limit) | RESOURCE_EXHAUSTED | 429 |
| Internal error | INTERNAL | 500 |
| Service unavailable | UNAVAILABLE | 503 |

### Error Response Format

```protobuf
message Error {
  string code = 1;
  string message = 2;
  repeated ErrorDetail details = 3;
}

message ErrorDetail {
  string field = 1;
  string issue = 2;
}
```

**Example**:
```json
{
  "code": "INVALID_ARGUMENT",
  "message": "Validation failed",
  "details": [
    {
      "field": "items",
      "issue": "must contain at least one item"
    },
    {
      "field": "items[0].quantity",
      "issue": "must be greater than 0"
    }
  ]
}
```

---

## 🧪 API Testing

### grpcurl Examples

**Create Order**:
```bash
grpcurl -plaintext \
  -d '{
    "user_id": "user-123",
    "items": [{"product_id": "prod-456", "quantity": 2}]
  }' \
  localhost:5001 order.v1.OrderService/CreateOrder
```

**Get Order**:
```bash
grpcurl -plaintext \
  -d '{"order_id": "order-789"}' \
  localhost:5001 order.v1.OrderService/GetOrder
```

### Postman Collection (REST)

```json
{
  "info": {
    "name": "Titan Commerce API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Orders",
      "item": [
        {
          "name": "Create Order",
          "request": {
            "method": "POST",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{access_token}}"
              }
            ],
            "url": "{{base_url}}/v1/orders",
            "body": {
              "mode": "raw",
              "raw": "{\n  \"user_id\": \"{{user_id}}\",\n  \"items\": [...]\n}"
            }
          }
        }
      ]
    }
  ]
}
```

---

**Document Version**: 1.0  
**Last Updated**: 2025-12-04  
**Pages**: 60+ (complete API reference)
