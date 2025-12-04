﻿# Warehouse Service 🏭

Multi-warehouse inventory management with intelligent stock allocation.

## Features

- 🏢 **Multi-warehouse**: Manage multiple warehouse locations
- 📦 **Stock allocation**: Priority-based and proximity-based allocation
- 🔄 **Stock transfers**: Transfer stock between warehouses
- 📊 **Movement tracking**: Complete audit trail of stock movements
- 🗺️ **Location-aware**: Allocate from nearest warehouse

## Technology Stack

- **Database**: PostgreSQL (transactional consistency)
- **Messaging**: Kafka for stock events
- **API**: gRPC

## Key Components

### Domain Layer
- `Warehouse`: Warehouse aggregate with location and capacity
- `WarehouseStock`: Stock levels per warehouse/product
- `StockMovement`: Audit trail of all stock movements

### Application Layer
- `CreateWarehouse`: Create new warehouse location
- `AllocateStock`: Smart allocation across warehouses
- `TransferStock`: Transfer between warehouses
- `RecordStockMovement`: Track inbound/outbound movements

### Allocation Strategy

1. **Priority-based**: Lower priority number = higher preference
2. **Proximity-based**: Calculate distance to customer
3. **Availability**: Only allocate from active warehouses

## Database Schema

### warehouses
- Warehouse locations with GPS coordinates
- Status: ACTIVE, INACTIVE, MAINTENANCE
- Priority for allocation ordering

### warehouse_stock
- Stock levels per warehouse/product combination
- Available vs Reserved quantities
- Zone and bin location tracking

### stock_movements
- Complete audit trail
- Movement types: INBOUND, OUTBOUND, TRANSFER, ADJUSTMENT, RETURN

## Movement Types

- `INBOUND`: Stock arriving at warehouse
- `OUTBOUND`: Stock leaving for orders
- `TRANSFER`: Inter-warehouse transfers
- `ADJUSTMENT`: Inventory adjustments
- `RETURN`: Customer returns

## Stock Allocation Example

```go
// Product needs 60 units
Warehouse A (Priority 1, Singapore): 50 units available
Warehouse B (Priority 2, Malaysia): 100 units available

// Result:
Allocation = [
  { warehouse_id: "A", quantity: 50 },
  { warehouse_id: "B", quantity: 10 }
]
```

## API Examples

### Create Warehouse
```protobuf
CreateWarehouse(
  name: "Singapore Main Warehouse",
  code: "SG-001",
  address: { city: "Singapore", latitude: 1.29, longitude: 103.85 },
  capacity: 100000,
  priority: 1
)
```

### Allocate Stock
```protobuf
AllocateStock(
  product_id: "PROD-123",
  required_quantity: 50,
  customer_address: "Singapore"
)
```

## Status

✅ **Implemented** - Core functionality complete
