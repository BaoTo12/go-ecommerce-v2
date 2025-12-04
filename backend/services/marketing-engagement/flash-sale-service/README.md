# Flash Sale Service

High-concurrency limited-time sales (11.11, Black Friday style).

## Purpose
Manages flash sale events with high traffic, limited stock, and time constraints. Optimized for events like Singles' Day (11.11).

## Technology Stack
- **Database**: Redis (inventory, rate limiting)
- **Cache**: Redis for hot data
- **API**: gRPC

## Key Features
- ✅ High-concurrency stock management
- ✅ Time-based sale activation
- ✅ Stock reservation with TTL
- ✅ Per-user purchase limits
- ✅ Real-time stock updates
- ✅ Conversion tracking
- ✅ Countdown timers
- ✅ Pre-sale notifications
- ✅ Redis Lua scripts for atomic operations

## Performance
- Handles 10,000+ TPS
- Redis-based atomic stock deduction
- Optimistic locking
- CDN for static assets

## API
- `CreateFlashSale`: Setup flash sale
- `StartFlashSale`: Activate sale
- `RecordPurchase`: Atomic stock deduction
- `GetActiveFlashSales`: List live sales
