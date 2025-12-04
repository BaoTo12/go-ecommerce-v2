-- Tracking info table (main tracking record)
CREATE TABLE IF NOT EXISTS tracking_info (
    tracking_number TEXT PRIMARY KEY,
    shipment_id TEXT,
    carrier TEXT,
    origin TEXT,
    destination TEXT,
    current_status TEXT,
    current_location_city TEXT,
    current_location_state TEXT,
    current_location_country TEXT,
    current_location_lat DOUBLE,
    current_location_lon DOUBLE,
    estimated_delivery TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Index for shipment_id lookups
CREATE INDEX IF NOT EXISTS idx_tracking_shipment ON tracking_info (shipment_id);

-- Tracking events table (time-series data)
CREATE TABLE IF NOT EXISTS tracking_events (
    tracking_number TEXT,
    timestamp TIMESTAMP,
    event_id TEXT,
    event_type TEXT,
    location_city TEXT,
    location_state TEXT,
    location_country TEXT,
    location_lat DOUBLE,
    location_lon DOUBLE,
    description TEXT,
    carrier TEXT,
    facility_name TEXT,
    PRIMARY KEY (tracking_number, timestamp, event_id)
) WITH CLUSTERING ORDER BY (timestamp DESC);

