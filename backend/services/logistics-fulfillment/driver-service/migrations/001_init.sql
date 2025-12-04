-- Driver Service Database Schema
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_deliveries_updated_at BEFORE UPDATE ON deliveries

    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_routes_updated_at BEFORE UPDATE ON delivery_routes

    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_drivers_updated_at BEFORE UPDATE ON drivers
-- Triggers

$$ LANGUAGE plpgsql;
END;
    RETURN NEW;
    NEW.updated_at = CURRENT_TIMESTAMP;
BEGIN
RETURNS TRIGGER AS $$
CREATE OR REPLACE FUNCTION update_updated_at_column()
-- Function to update updated_at timestamp

CREATE INDEX idx_deliveries_status ON deliveries(status);
CREATE INDEX idx_deliveries_driver ON deliveries(driver_id);
CREATE INDEX idx_deliveries_route ON deliveries(route_id);
CREATE INDEX idx_deliveries_shipment ON deliveries(shipment_id);

);
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    pod_completed_at TIMESTAMP,
    recipient_name VARCHAR(255),
    proof_signature_url TEXT,
    proof_photo_url TEXT,
    actual_arrival TIMESTAMP,
    estimated_arrival TIMESTAMP,
    sequence_number INT,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    delivery_lng DECIMAL(11, 8) NOT NULL,
    delivery_lat DECIMAL(10, 8) NOT NULL,
    delivery_address TEXT NOT NULL,
    customer_phone VARCHAR(50) NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    driver_id VARCHAR(36) REFERENCES drivers(driver_id),
    route_id VARCHAR(36) REFERENCES delivery_routes(route_id),
    shipment_id VARCHAR(36) NOT NULL,
    delivery_id VARCHAR(36) PRIMARY KEY,
CREATE TABLE IF NOT EXISTS deliveries (
-- Deliveries table

CREATE INDEX idx_routes_status ON delivery_routes(status);
CREATE INDEX idx_routes_driver ON delivery_routes(driver_id);

);
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    started_at TIMESTAMP,
    actual_duration INT,
    estimated_duration INT, -- in minutes
    completed_deliveries INT DEFAULT 0,
    total_deliveries INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'CREATED',
    driver_id VARCHAR(36) NOT NULL REFERENCES drivers(driver_id),
    route_id VARCHAR(36) PRIMARY KEY,
CREATE TABLE IF NOT EXISTS delivery_routes (
-- Delivery routes table

CREATE INDEX idx_drivers_location ON drivers(current_lat, current_lng);
CREATE INDEX idx_drivers_status ON drivers(status);

);
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    successful_deliveries INT DEFAULT 0,
    total_deliveries INT DEFAULT 0,
    rating DECIMAL(3, 2) DEFAULT 5.0,
    current_lng DECIMAL(11, 8),
    current_lat DECIMAL(10, 8),
    status VARCHAR(50) NOT NULL DEFAULT 'OFF_DUTY',
    license_plate VARCHAR(50),
    vehicle_type VARCHAR(50) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    driver_id VARCHAR(36) PRIMARY KEY,
CREATE TABLE IF NOT EXISTS drivers (
-- Drivers table


