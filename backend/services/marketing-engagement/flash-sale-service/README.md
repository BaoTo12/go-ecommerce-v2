# Flash Sale Service ⚡

Solving **"The 11.11 Problem"** - 1M concurrent users hitting "Buy" at 00:00:00.

## Features

- ⚡ Handle 1M concurrent users
- 🛡️ Token bucket rate limiting (10K req/sec per user)
- 🤖 Proof-of-Work (PoW) challenge to prevent bots
- 🔒 Redis atomic inventory (Lua scripts)
- ⏱️ WebSocket countdown synchronization
- 📬 Queue-based load leveling (Kafka → worker pool)
- 🚀 <100ms response time (reservation ID)

## Architecture

```
1M Users → PoW Challenge → Rate Limiter → Redis Atomic Decrement
                                              ↓
                                      Reserve Inventory
                                              ↓
                                    Kafka Queue (async)
                                              ↓
                                      Worker Pool → Create Order
```

## Lua Script (Atomic Inventory)

```lua
local key = KEYS[1]
local qty = tonumber(ARGV[1])
local current = tonumber(redis.call('GET', key) or 0)

if current >= qty then
  redis.call('DECRBY', key, qty)
  return 1  -- Success
else
  return 0  -- Out of stock
end
```

## Status

🚧 **Under Development** - Skeleton structure created
