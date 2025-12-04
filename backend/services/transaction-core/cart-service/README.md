# Cart Service

Redis-based shopping cart with sub-10ms latency and auto-save.

## Features

- ✅ Redis for <10ms latency
- ✅ Auto-save every 5 seconds
- ✅ TTL management (7 days expiry)
- ✅ Atomic cart operations (add, remove, update quantity)
- ✅ Real-time sync across devices

## Quick Start

```bash
export SERVICE_NAME=cart-service
export CELL_ID=cell-001
export REDIS_ADDR=localhost:6379
go run cmd/server/main.go
```

## Status

🚧 **Under Development** - Skeleton structure created
