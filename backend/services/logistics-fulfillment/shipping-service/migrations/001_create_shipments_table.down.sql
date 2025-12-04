-- Drop indexes
DROP INDEX IF EXISTS idx_shipments_created;
DROP INDEX IF EXISTS idx_shipments_status;
DROP INDEX IF EXISTS idx_shipments_tracking;
DROP INDEX IF EXISTS idx_shipments_order;

-- Drop table
DROP TABLE IF EXISTS shipments;

