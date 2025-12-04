# Driver Service

Last-mile delivery driver management and route optimization.

## Purpose
Manages delivery drivers, routes, and real-time delivery status updates for last-mile fulfillment.

## Technology Stack
- **Database**: PostgreSQL (for transactional driver and route data)
- **API**: gRPC

## Key Features
- ✅ Driver profile and availability management
- ✅ Route assignment and optimization
- ✅ Real-time delivery status updates
- ✅ Driver location tracking
- ✅ Performance metrics and ratings
- ✅ Multi-delivery batch routing

## Domain Model

### Driver
- Profile information (name, phone, vehicle type)
- Current status (AVAILABLE, ON_DELIVERY, OFF_DUTY)
- Location tracking
- Performance metrics

### DeliveryRoute
- Assigned driver
- Multiple delivery stops
- Optimized route sequence
- Status tracking

### Delivery
- Shipment association
- Customer location
- Time windows
- Proof of delivery

## Quick Start

```bash
export SERVICE_NAME=driver-service
export CELL_ID=cell-001
export DB_HOST=localhost
export DB_PORT=5432
go run cmd/server/main.go
```

## API Overview

### Commands
- `RegisterDriver`: Register new driver
- `UpdateDriverStatus`: Update availability status
- `AssignRoute`: Assign deliveries to driver
- `UpdateDeliveryStatus`: Update delivery status
- `RecordProofOfDelivery`: Record POD (photo, signature)

### Queries
- `GetDriver`: Get driver details
- `GetAvailableDrivers`: List drivers ready for assignment
- `GetDriverRoute`: Get current route
- `GetDeliveryStatus`: Get delivery status

## Integration

### Events Published
- `DriverRegistered`: New driver onboarded
- `RouteAssigned`: Route assigned to driver
- `DeliveryCompleted`: Delivery successfully completed
- `DeliveryFailed`: Delivery attempt failed

### Events Consumed
- `ShipmentReadyForDelivery`: From shipping-service
- `DeliveryRescheduled`: From customer service

## Database Schema

See `migrations/001_init.sql` for complete schema.

