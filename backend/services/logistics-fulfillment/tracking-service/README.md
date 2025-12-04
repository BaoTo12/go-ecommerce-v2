﻿# Tracking Service 📍

Real-time package tracking with ScyllaDB for time-series event storage.

## Features

- 🔄 **Real-time updates**: Track package location in real-time
- 📊 **Event history**: Complete event timeline stored in ScyllaDB
- 🗺️ **Location tracking**: GPS coordinates and facility information
- 🔔 **Webhook support**: Carrier webhook integrations
- 📡 **Live streaming**: Subscribe to tracking updates via gRPC streams

## Technology Stack

- **Database**: ScyllaDB (time-series data optimized)
- **Messaging**: Kafka for event streaming
- **API**: gRPC with streaming support

## Key Components

### Domain Layer
- `TrackingInfo`: Main aggregate for tracking information
- `TrackingEvent`: Individual tracking events (time-series)
- `Location`: GPS and address information

### Application Layer
- `CreateTracking`: Initialize tracking for a shipment
- `UpdateLocation`: Add new tracking event
- `GetTrackingHistory`: Retrieve complete event history
- `GetCurrentStatus`: Get latest tracking status

### Infrastructure
- **ScyllaDB Repository**: Time-series optimized storage
- **Kafka Producer**: Publish tracking events
- **Webhook Handler**: Receive carrier updates

## Database Schema

### tracking_info
```cql
CREATE TABLE tracking_info (
    tracking_number TEXT PRIMARY KEY,
    shipment_id TEXT,
    carrier TEXT,
    current_status TEXT,
    current_location_city TEXT,
    estimated_delivery TIMESTAMP
);
```

### tracking_events (Time-Series)
```cql
CREATE TABLE tracking_events (
    tracking_number TEXT,
    timestamp TIMESTAMP,
    event_id TEXT,
    event_type TEXT,
    location_city TEXT,
    PRIMARY KEY (tracking_number, timestamp, event_id)
) WITH CLUSTERING ORDER BY (timestamp DESC);
```

## Event Types

- `PICKED_UP`: Package picked up from sender
- `IN_TRANSIT`: Package in transit
- `AT_FACILITY`: Package at sorting facility
- `OUT_FOR_DELIVERY`: Out for final delivery
- `DELIVERED`: Successfully delivered
- `EXCEPTION`: Delivery exception
- `RETURNED`: Returned to sender

## API Examples

### Create Tracking
```protobuf
CreateTracking(
  tracking_number: "TRK12345",
  shipment_id: "SHP-001",
  carrier: "DHL"
)
```

### Update Location
```protobuf
UpdateLocation(
  tracking_number: "TRK12345",
  event_type: IN_TRANSIT,
  location: { city: "Singapore", country: "SG" }
)
```

## Status

✅ **Implemented** - Core functionality complete
