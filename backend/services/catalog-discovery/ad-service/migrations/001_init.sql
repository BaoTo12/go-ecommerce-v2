-- Ad Service Database Migration

CREATE TABLE IF NOT EXISTS campaigns (
    id VARCHAR(36) PRIMARY KEY,
    seller_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    budget DECIMAL(10, 2) NOT NULL,
    remaining_budget DECIMAL(10, 2) NOT NULL,
    bid_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ad_events (
    id SERIAL PRIMARY KEY,
    ad_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36),
    event_type VARCHAR(20) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_campaigns_status ON campaigns(status);
CREATE INDEX idx_ad_events_ad_id ON ad_events(ad_id);
