# Inventory Service 📦

Atomic inventory management with Redis Lua scripts - **NO overselling**.

## Features

- 🔒 **Redis Lua scripts** for atomic stock operations
- 📝 **Reservation system**: Reserve → Commit or Rollback
- 🚨 **Stock alerts** (low stock notifications)
- 🏢 **Multi-warehouse** support
- ⚡ **Real-time** stock synchronization

## Critical Lua Script (Atomic Decrement)

```lua
-- Key: product:<product_id>:stock
-- Arg: quantity to decrement

local key = KEYS[1]
local quantity = tonumber(ARGV[1])
local current = tonumber(redis.call('GET', key) or 0)

if current >= quantity then
  redis.call('DECRBY', key, quantity)
  return 1  -- Success
else
  return 0  -- Insufficient stock
end
```

## Reservation Flow

```
1. Check Stock → 2. Reserve (atomic) → 3. Process Payment
                           ↓                      ↓
                    Reservation ID         Success: Commit
                                          Failure: Rollback (release)
```

## Multi-Warehouse Example

```go
// Product has stock in 3 warehouses
Warehouse A: 50 units (priority: 1, closest to customer)
Warehouse B: 100 units (priority: 2)
Warehouse C: 200 units (priority: 3)

// User orders 60 units
Allocation:
- Warehouse A: 50 units
- Warehouse B: 10 units
Total: 60 units fulfilled
```

## Status

🚧 **Under Development** - Skeleton structure created
