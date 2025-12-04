-- User Service Database Migration

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255),
    phone_number VARCHAR(50),
    avatar_url TEXT,
    addresses JSONB,
    preferences JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_users_email ON users(email);
