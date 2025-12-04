-- Create warehouses table
CREATE INDEX idx_stock_movements_created ON stock_movements(created_at DESC);
CREATE INDEX idx_stock_movements_product ON stock_movements(product_id);
CREATE INDEX idx_stock_movements_warehouse ON stock_movements(warehouse_id);
CREATE INDEX idx_warehouse_stock_product ON warehouse_stock(product_id);
CREATE INDEX idx_warehouses_priority ON warehouses(priority);
CREATE INDEX idx_warehouses_status ON warehouses(status);
-- Create indexes

);
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(warehouse_id) ON DELETE CASCADE
    created_by VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    reference_id VARCHAR(255),
    quantity INTEGER NOT NULL,
    movement_type VARCHAR(50) NOT NULL,
    product_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    movement_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
CREATE TABLE IF NOT EXISTS stock_movements (
-- Create stock_movements table

);
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(warehouse_id) ON DELETE CASCADE
    PRIMARY KEY (warehouse_id, product_id),
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    bin_location VARCHAR(50),
    zone VARCHAR(50),
    total_quantity INTEGER NOT NULL DEFAULT 0,
    reserved_quantity INTEGER NOT NULL DEFAULT 0,
    available_quantity INTEGER NOT NULL DEFAULT 0,
    product_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
CREATE TABLE IF NOT EXISTS warehouse_stock (
-- Create warehouse_stock table

);
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    priority INTEGER NOT NULL DEFAULT 10,
    capacity INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    longitude DECIMAL(11, 8),
    latitude DECIMAL(10, 8),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    state VARCHAR(100),
    city VARCHAR(100),
    street VARCHAR(255),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    warehouse_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
CREATE TABLE IF NOT EXISTS warehouses (

