-- Drop indexes
DROP INDEX IF EXISTS idx_stock_movements_created;
DROP INDEX IF EXISTS idx_stock_movements_product;
DROP INDEX IF EXISTS idx_stock_movements_warehouse;
DROP INDEX IF EXISTS idx_warehouse_stock_product;
DROP INDEX IF EXISTS idx_warehouses_priority;
DROP INDEX IF EXISTS idx_warehouses_status;

-- Drop tables
DROP TABLE IF EXISTS stock_movements;
DROP TABLE IF EXISTS warehouse_stock;
DROP TABLE IF EXISTS warehouses;

