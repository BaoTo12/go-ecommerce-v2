# Complete Microservices & Cell Architecture

## All 25+ Microservices with Cell-Based Deployment

**IMPORTANT**: This shows the COMPLETE microservices architecture where **each cell contains ALL services**, not just one service.

---

## 🏗️ Complete Service Catalog

### 1. Transaction Core Cluster (High Integrity)

| Service | Port | Database | Purpose |
|---------|------|----------|---------|
| **order-service** | 50001 | CockroachDB | Order lifecycle management |
| **payment-service** | 50002 | CockroachDB | Payment processing (Stripe, PayPal) |
| **cart-service** | 50003 | Redis | Shopping cart (ephemeral) |
| **checkout-service** | 50004 | CockroachDB | Checkout orchestration (Saga) |
| **wallet-service** | 50005 | CockroachDB | Escrow, seller payouts |
| **refund-service** | 50006 | CockroachDB | Refund processing |
| **voucher-service** | 50007 | Redis + Postgres | Voucher management |

### 2. Catalog & Discovery Cluster (Read-Heavy)

| Service | Port | Database | Purpose |
|---------|------|----------|---------|
| **product-service** | 50010 | MongoDB | Product catalog (100M+ products) |
| **search-service** | 50011 | Elasticsearch | Full-text search, filters |
| **recommendation-service** | 50012 | ScyllaDB | ML-powered recommendations |
| **category-service** | 50013 | PostgreSQL | Category hierarchy |
| **seller-service** | 50014 | PostgreSQL | Seller profiles, shops |
| **review-service** | 50015 | MongoDB | Product reviews, ratings |

### 3. User & Social Cluster

| Service | Port | Database | Purpose |
|---------|------|----------|---------|
| **user-service** | 50020 | PostgreSQL | User profiles, preferences |
| **auth-service** | 50021 | Redis + Postgres | Authentication (JWT) |
| **social-service** | 50022 | MongoDB | User following, social graph |
| **feed-service** | 50023 | ScyllaDB | Activity feed (TikTok-style) |
| **notification-service** | 50024 | Redis + FCM | Push notifications |

### 4. Communication Cluster

| Service | Port | Database | Purpose |
|---------|------|----------|---------|
| **chat-service** | 50030 | ScyllaDB | Buyer-seller chat |
| **livestream-service** | 50031 | Redis + S3 | Live shopping streams |
| **videocall-service** | 50032 | WebRTC mesh | Video calls (customer support) |

### 5. Logistics & Fulfillment Cluster

| Service | Port | Database | Purpose |
|---------|------|----------|---------|
| **shipping-service** | 50040 | PostgreSQL | Shipping calculations |
| **tracking-service** | 50041 | ScyllaDB | Real-time package tracking |
| **warehouse-service** | 50042 | PostgreSQL | Warehouse management |
| **inventory-service** | 50043 | Redis + Postgres | Stock management |

### 6. Marketing & Engagement Cluster

| Service | Port | Database | Purpose |
|---------|------|----------|---------|
| **flash-sale-service** | 50050 | Redis | Limited-time sales (11.11) |
| **gamification-service** | 50051 | PostgreSQL | Coins, badges, rewards |
| **campaign-service** | 50052 | PostgreSQL | Marketing campaigns |
| **coupon-service** | 50053 | Redis | Coupon generation |

### 7. Intelligence & Analytics Cluster

| Service | Port | Database | Purpose |
|---------|------|----------|---------|
| **pricing-service** | 50060 | ClickHouse | Dynamic pricing (ML) |
| **fraud-service** | 50061 | ClickHouse | Fraud detection |
| **analytics-service** | 50062 | ClickHouse | Business intelligence |
| **ab-testing-service** | 50063 | PostgreSQL | A/B experiment management |

---

## 🔥 Cell-Based Architecture (The Key Innovation)

### What is a Cell?

**A cell is a COMPLETE, self-contained deployment of ALL microservices above.**

```
┌─────────────────────────────────────────────────────────────┐
│                      CELL #1 (Users 1-10,000)                │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Order Svc   │  │ Payment Svc │  │ Cart Svc    │         │
│  │ (Pod 1-3)   │  │ (Pod 1-3)   │  │ (Pod 1-3)   │  ...30+ │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Product Svc │  │ Search Svc  │  │ Chat Svc    │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│                                                               │
│  Databases (Cell-Local):                                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ ScyllaDB    │  │ Redis       │  │ Postgres    │         │
│  │ (Replicas)  │  │ (Cluster)   │  │ (Cell DB)   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                      CELL #2 (Users 10,001-20,000)           │
│         (Exact same 30+ services, isolated data)             │
└─────────────────────────────────────────────────────────────┘

... (498 more cells)

┌─────────────────────────────────────────────────────────────┐
│                      CELL #500 (Users 4,990,000-5,000,000)   │
│         (Exact same 30+ services, isolated data)             │
└─────────────────────────────────────────────────────────────┘
```

### Directory Structure Showing Cell Architecture

```
backend/
├── services/                           # All microservices (deployed in each cell)
│   │
│   ├── transaction-core/               # Critical transaction services
│   │   ├── order-service/
│   │   │   ├── cmd/
│   │   │   │   └── server/
│   │   │   │       └── main.go
│   │   │   ├── internal/
│   │   │   │   ├── domain/
│   │   │   │   │   ├── order/
│   │   │   │   │   │   ├── order.go          # Aggregate
│   │   │   │   │   │   ├── order_item.go
│   │   │   │   │   │   └── repository.go
│   │   │   │   │   └── events/
│   │   │   │   ├── application/
│   │   │   │   │   ├── commands/
│   │   │   │   │   │   ├── create_order.go
│   │   │   │   │   │   ├── cancel_order.go
│   │   │   │   │   │   └── update_status.go
│   │   │   │   │   └── queries/
│   │   │   │   │       ├── get_order.go
│   │   │   │   │       └── list_orders.go
│   │   │   │   ├── infrastructure/
│   │   │   │   │   ├── persistence/
│   │   │   │   │   │   ├── cockroachdb/
│   │   │   │   │   │   │   └── order_repository.go
│   │   │   │   │   │   └── scylladb/
│   │   │   │   │   │       └── event_store.go
│   │   │   │   │   ├── messaging/
│   │   │   │   │   │   └── kafka/
│   │   │   │   │   │       ├── producer.go
│   │   │   │   │   │       └── consumer.go
│   │   │   │   │   └── grpc/
│   │   │   │   │       ├── server.go
│   │   │   │   │       └── interceptors.go
│   │   │   │   └── interfaces/
│   │   │   │       └── grpc/
│   │   │   │           └── order_handler.go
│   │   │   ├── proto/
│   │   │   │   └── order/
│   │   │   │       └── v1/
│   │   │   │           └── order.proto
│   │   │   ├── migrations/
│   │   │   ├── Dockerfile
│   │   │   └── go.mod
│   │   │
│   │   ├── payment-service/            # Same structure
│   │   ├── cart-service/
│   │   ├── checkout-service/           # Saga coordinator
│   │   ├── wallet-service/             # Escrow management
│   │   ├── refund-service/
│   │   └── voucher-service/
│   │
│   ├── catalog-discovery/              # Product catalog cluster
│   │   ├── product-service/
│   │   │   ├── cmd/server/main.go
│   │   │   ├── internal/
│   │   │   │   ├── domain/
│   │   │   │   │   └── product/
│   │   │   │   │       ├── product.go
│   │   │   │   │       ├── variant.go      # Size, color
│   │   │   │   │       ├── pricing.go
│   │   │   │   │       └── repository.go
│   │   │   │   ├── application/
│   │   │   │   ├── infrastructure/
│   │   │   │   │   ├── persistence/
│   │   │   │   │   │   └── mongodb/
│   │   │   │   │   │       └── product_repository.go
│   │   │   │   │   └── cache/
│   │   │   │   │       └── redis_cache.go
│   │   │   │   └── interfaces/
│   │   │   ├── proto/product/v1/
│   │   │   └── ...
│   │   ├── search-service/
│   │   │   ├── internal/
│   │   │   │   ├── elasticsearch/
│   │   │   │   │   ├── indexer.go
│   │   │   │   │   └── search.go
│   │   │   │   └── ranking/
│   │   │   │       └── ml_ranker.go
│   │   │   └── ...
│   │   ├── recommendation-service/
│   │   │   ├── internal/
│   │   │   │   ├── ml/
│   │   │   │   │   ├── collaborative_filtering.go
│   │   │   │   │   ├── content_based.go
│   │   │   │   │   └── hybrid_recommender.go
│   │   │   │   └── serving/
│   │   │   │       └── recommendation_engine.go
│   │   │   └── ...
│   │   ├── category-service/
│   │   ├── seller-service/
│   │   │   ├── internal/
│   │   │   │   └── domain/
│   │   │   │       └── seller/
│   │   │   │           ├── seller.go
│   │   │   │           ├── shop.go
│   │   │   │           ├── performance_metrics.go
│   │   │   │           └── repository.go
│   │   │   └── ...
│   │   └── review-service/
│   │       ├── internal/
│   │       │   └── domain/
│   │       │       └── review/
│   │       │           ├── review.go
│   │       │           ├── rating.go
│   │       │           ├── moderation.go      # AI-powered
│   │       │           └── repository.go
│   │       └── ...
│   │
│   ├── user-social/                    # User & social features
│   │   ├── user-service/
│   │   ├── auth-service/
│   │   │   ├── internal/
│   │   │   │   ├── jwt/
│   │   │   │   │   ├── generator.go
│   │   │   │   │   └── validator.go
│   │   │   │   ├── oauth/
│   │   │   │   │   ├── google.go
│   │   │   │   │   └── facebook.go
│   │   │   │   └── 2fa/
│   │   │   │       └── totp.go
│   │   │   └── ...
│   │   ├── social-service/
│   │   │   ├── internal/
│   │   │   │   └── domain/
│   │   │   │       └── social/
│   │   │   │           ├── follow.go
│   │   │   │           ├── friend.go
│   │   │   │           └── social_graph.go
│   │   │   └── ...
│   │   ├── feed-service/               # TikTok-style feed
│   │   │   ├── internal/
│   │   │   │   ├── algorithm/
│   │   │   │   │   ├── ranking.go      # Feed ranking
│   │   │   │   │   └── personalization.go
│   │   │   │   └── cache/
│   │   │   │       └── feed_cache.go
│   │   │   └── ...
│   │   └── notification-service/
│   │       ├── internal/
│   │       │   ├── fcm/
│   │       │   │   └── push.go
│   │       │   └── apns/
│   │       │       └── push.go
│   │       └── ...
│   │
│   ├── communication/                  # Real-time communication
│   │   ├── chat-service/
│   │   │   ├── internal/
│   │   │   │   ├── websocket/
│   │   │   │   │   ├── gateway.go
│   │   │   │   │   ├── connection.go
│   │   │   │   │   └── router.go
│   │   │   │   ├── persistence/
│   │   │   │   │   └── scylla/
│   │   │   │   │       └── message_repository.go
│   │   │   │   └── presence/
│   │   │   │       └── online.go
│   │   │   └── ...
│   │   ├── livestream-service/         # 🔥 Shopee Live
│   │   │   ├── internal/
│   │   │   │   ├── streaming/
│   │   │   │   │   ├── rtmp_server.go
│   │   │   │   │   ├── hls_packager.go
│   │   │   │   │   └── cdn_pusher.go
│   │   │   │   ├── chat/
│   │   │   │   │   └── live_chat.go    # Chat during live
│   │   │   │   ├── products/
│   │   │   │   │   └── pinned_products.go
│   │   │   │   └── analytics/
│   │   │   │       └── viewer_stats.go
│   │   │   └── ...
│   │   └── videocall-service/
│   │       ├── internal/
│   │       │   ├── webrtc/
│   │       │   │   ├── signaling.go
│   │       │   │   └── turn_server.go
│   │       │   └── recording/
│   │       │       └── call_recorder.go
│   │       └── ...
│   │
│   ├── logistics-fulfillment/          # Shipping & tracking
│   │   ├── shipping-service/
│   │   │   ├── internal/
│   │   │   │   ├── carriers/
│   │   │   │   │   ├── fedex.go
│   │   │   │   │   ├── ups.go
│   │   │   │   │   └── dhl.go
│   │   │   │   └── rate_calculator.go
│   │   │   └── ...
│   │   ├── tracking-service/
│   │   │   ├── internal/
│   │   │   │   ├── realtime/
│   │   │   │   │   └── gps_tracker.go
│   │   │   │   └── webhook/
│   │   │   │       └── carrier_webhook.go
│   │   │   └── ...
│   │   ├── warehouse-service/
│   │   └── inventory-service/
│   │       ├── internal/
│   │       │   ├── atomic/
│   │       │   │   └── redis_lua.go    # Lua scripts
│   │       │   └── reservation/
│   │       │       └── ttl_manager.go
│   │       └── ...
│   │
│   ├── marketing-engagement/           # Marketing cluster
│   │   ├── flash-sale-service/         # 🔥 11.11 Flash Sales
│   │   │   ├── internal/
│   │   │   │   ├── countdown/
│   │   │   │   │   └── websocket_sync.go
│   │   │   │   ├── admission/
│   │   │   │   │   ├── token_bucket.go
│   │   │   │   │   └── proof_of_work.go
│   │   │   │   ├── inventory/
│   │   │   │   │   └── atomic_decrement.go
│   │   │   │   └── queue/
│   │   │   │       └── kafka_leveling.go
│   │   │   └── ...
│   │   ├── gamification-service/       # 🎮 Shopee Coins, Games
│   │   │   ├── internal/
│   │   │   │   ├── coins/
│   │   │   │   │   ├── balance.go
│   │   │   │   │   └── transaction.go
│   │   │   │   ├── games/
│   │   │   │   │   ├── shake_shake.go  # Shake phone game
│   │   │   │   │   ├── lucky_draw.go
│   │   │   │   │   └── daily_checkin.go
│   │   │   │   └── rewards/
│   │   │   │       └── reward_engine.go
│   │   │   └── ...
│   │   ├── campaign-service/
│   │   └── coupon-service/
│   │
│   └── intelligence-analytics/         # ML & Analytics
│       ├── pricing-service/
│       │   ├── internal/
│       │   │   ├── ml/
│       │   │   │   ├── price_elasticity.go
│       │   │   │   ├── competitor_scraper.go
│       │   │   │   └── onnx_model.go
│       │   │   └── realtime/
│       │   │       └── dynamic_pricer.go
│       │   └── ...
│       ├── fraud-service/
│       │   ├── internal/
│       │   │   ├── detection/
│       │   │   │   ├── rule_engine.go
│       │   │   │   ├── ml_scorer.go
│       │   │   │   └── device_fingerprint.go
│       │   │   └── graph/
│       │   │       └── fraud_ring_detection.go
│       │   └── ...
│       ├── analytics-service/
│       └── ab-testing-service/
│
├── cell-router/                        # 🔥 Routes users to cells
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── routing/
│   │   │   ├── consistent_hash.go      # Hash(UserID) → Cell
│   │   │   ├── health_check.go         # Cell health monitoring
│   │   │   └── failover.go             # Reroute on cell failure
│   │   └── registry/
│   │       └── cell_registry.go        # List of all 500 cells
│   └── ...
│
└── pkg/                                # Shared libraries
    ├── logger/
    ├── metrics/
    ├── tracing/
    └── errors/
```

---

## 🎯 Complete Service Count

**Total Services per Cell**: 30+ microservices  
**Total Cells**: 500  
**Total Service Instances**: 15,000+ (30 services × 500 cells)  
**Total Pods**: 45,000+ (averaging 3 pods per service)

---

## 🌟 Modern Shopee Features Implemented

### 1. Shopee Live (Livestream Shopping)
```
livestream-service/
├── RTMP streaming ingestion
├── HLS video packaging
├──CDN distribution (CloudFlare Stream)
├── Live chat overlay
├── Pinned products during stream
├── Flash sale during live
└── Viewer analytics
```

### 2. Social Shopping
```
social-service/
├── Follow sellers
├── Share products to social media
├── User-generated content feed
├── Product reviews with photos
└── Social proof (X people bought this)
```

### 3. Gamification (Shopee Coins, Games)
```
gamification-service/
├── Shopee Coins system
├── Daily check-in rewards
├── Shake-shake game
├── Lucky draw
├── Missions & challenges
└── Voucher redemption
```

### 4. Flash Sales & Mega Sales
```
flash-sale-service/
├── Countdown timers (WebSocket sync)
├── 1M concurrent user handling
├── Atomic inventory (Redis Lua)
├── Queue-based load leveling
└── Bot prevention (PoW)
```

### 5. AI-Powered Features
```
recommendation-service/ - Personalized product feed
pricing-service/        - Dynamic pricing
fraud-service/          - Real-time fraud detection
review-service/         - Auto-moderation (spam detection)
search-service/         - Semantic search
```

---

## 📊 Deployment Model

### Kubernetes Manifest for Cell #1

```yaml
# infrastructure/kubernetes/cells/cell-001.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cell-001
  labels:
    cell-id: "001"
    user-range: "1-10000"
---
# Deploy ALL 30+ services in this cell
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-service
  namespace: cell-001
spec:
  replicas: 3
  selector:
    matchLabels:
      app: order-service
      cell-id: "001"
  template:
    metadata:
      labels:
        app: order-service
        cell-id: "001"
    spec:
      containers:
      - name: order-service
        image: titan/order-service:v1.0.0
        env:
        - name: CELL_ID
          value: "001"
        - name: DATABASE_URL
          value: "postgres://cell-001-db:5432/orders"
        ports:
        - containerPort: 50001
---
# Similar deployments for all other 29 services...
apiVersion: apps/v1
kind: Deployment
metadata:
  name: livestream-service
  namespace: cell-001
spec:
  replicas: 5  # More replicas for livestreaming
  # ... (same pattern for all services)
```

**Total YAML files**: 500 (one per cell) × 30 services = 15,000 Kubernetes manifests (generated programmatically)

---

**Document Version**: 2.0  
**Last Updated**: 2025-12-04  
**Coverage**: Complete 30+ microservices with cell architecture
