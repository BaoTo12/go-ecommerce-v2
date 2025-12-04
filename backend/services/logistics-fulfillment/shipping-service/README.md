﻿# Shipping Service 📦

Multi-carrier shipping integration with rate calculation and label generation.

## Features

- 🚚 **Multi-carrier support**: DHL, FedEx, UPS integration
- 💰 **Rate shopping**: Compare rates across carriers
- 🏷️ **Label generation**: Generate shipping labels
- 📍 **Real-time tracking**: Integrate with carrier tracking APIs
- 🔄 **Webhook handling**: Receive carrier status updates

## Technology Stack

- **Database**: PostgreSQL (transactional data)
- **Carriers**: DHL, FedEx, UPS APIs
- **Messaging**: Kafka for shipment events
- **API**: gRPC

## Supported Carriers

### DHL Express
- Standard shipping
- Express shipping
- International shipping

### FedEx
- FedEx Ground
- FedEx Express
- FedEx International

### UPS
- UPS Ground
- UPS Next Day Air
- UPS Worldwide Express

## Key Components

### Domain Layer
- `Shipment`: Shipment aggregate with tracking info
- `ShipmentStatus`: Lifecycle states

### Application Layer
- `CreateShipment`: Create shipment with carrier
- `CalculateShippingCost`: Get rates from carriers
- `UpdateStatus`: Update shipment status
- `GetShipment`: Query shipment details

### Infrastructure - Carrier Integrations
- **DHL Carrier**: DHL API integration
- **FedEx Carrier**: FedEx API integration  
- **UPS Carrier**: UPS API integration
- **Webhook Handler**: Process carrier callbacks

## Shipment Status Flow

```
PENDING → PICKED_UP → IN_TRANSIT → OUT_FOR_DELIVERY → DELIVERED
                           ↓
                       RETURNED (if failed delivery)
```

## Rate Calculation Logic

```go
// Base cost + weight-based cost
DHL:   $5.00 + (weight * $2.50 * 1.2)
FedEx: $6.00 + (weight * $2.20 * 1.1)
UPS:   $5.50 + (weight * $2.00)
```

## API Examples

### Create Shipment
```protobuf
CreateShipment(
  order_id: "ORD-123",
  carrier: "DHL",
  origin_address: "123 Main St, Singapore",
  destination_address: "456 Oak Ave, Malaysia",
  weight: 2.5
)
```

### Calculate Shipping Cost
```protobuf
CalculateShippingCost(
  origin_address: "Singapore",
  destination_address: "Malaysia",
  weight: 2.5,
  carrier: "DHL"
)
// Response: { cost: 11.25, estimated_days: 3 }
```

## Carrier Integration

Each carrier implements the `Carrier` interface:
- `CalculateRate()`: Get shipping rate
- `CreateShipment()`: Generate label and tracking number
- `CancelShipment()`: Cancel shipment
- `GetTrackingInfo()`: Fetch latest tracking status

## Database Schema

```sql
CREATE TABLE shipments (
    shipment_id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    carrier VARCHAR(50),
    tracking_number VARCHAR(100) UNIQUE,
    status VARCHAR(50),
    shipping_cost DECIMAL(10, 2),
    estimated_delivery TIMESTAMP
);
```

## Events Published

- `ShipmentCreated`: When shipment is created
- `ShipmentPickedUp`: Picked up from sender
- `ShipmentInTransit`: In transit
- `ShipmentDelivered`: Successfully delivered

## Status

✅ **Implemented** - Core functionality complete with carrier integrations
