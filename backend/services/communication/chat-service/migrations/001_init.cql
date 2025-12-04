-- Chat Service ScyllaDB Schema

-- Messages table (time-series optimized)
CREATE TABLE IF NOT EXISTS messages (
    conversation_id TEXT,
    created_at TIMESTAMP,
    message_id TEXT,
    sender_id TEXT,
    sender_name TEXT,
    message_type TEXT,
    content TEXT,
    metadata TEXT, -- JSON
    status TEXT,
    is_edited BOOLEAN,
    is_deleted BOOLEAN,
    read_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (conversation_id, created_at, message_id)
) WITH CLUSTERING ORDER BY (created_at DESC)
  AND default_time_to_live = 0
  AND gc_grace_seconds = 86400;

-- Conversations table
CREATE TABLE IF NOT EXISTS conversations (
    conversation_id TEXT PRIMARY KEY,
    type TEXT,
    participant_ids LIST<TEXT>,
    participant_names MAP<TEXT, TEXT>,
    last_message_id TEXT,
    last_message_content TEXT,
    last_message_at TIMESTAMP,
    unread_counts MAP<TEXT, INT>,
    is_muted MAP<TEXT, BOOLEAN>,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- User conversations index (for fast user lookup)
CREATE TABLE IF NOT EXISTS user_conversations (
    user_id TEXT,
    updated_at TIMESTAMP,
    conversation_id TEXT,
    type TEXT,
    last_message_at TIMESTAMP,
    unread_count INT,
    is_muted BOOLEAN,
    PRIMARY KEY (user_id, updated_at, conversation_id)
) WITH CLUSTERING ORDER BY (updated_at DESC);

-- Typing indicators (short-lived, stored in memory with TTL)
CREATE TABLE IF NOT EXISTS typing_indicators (
    conversation_id TEXT,
    user_id TEXT,
    user_name TEXT,
    is_typing BOOLEAN,
    timestamp TIMESTAMP,
    PRIMARY KEY (conversation_id, user_id)
) WITH default_time_to_live = 10; -- Auto-expire after 10 seconds

-- Online status (short-lived with TTL)
CREATE TABLE IF NOT EXISTS online_status (
    user_id TEXT PRIMARY KEY,
    is_online BOOLEAN,
    last_seen_at TIMESTAMP,
    device_info TEXT
) WITH default_time_to_live = 300; -- Auto-expire after 5 minutes

-- Message delivery tracking
CREATE TABLE IF NOT EXISTS message_delivery (
    message_id TEXT,
    user_id TEXT,
    delivered_at TIMESTAMP,
    read_at TIMESTAMP,
    PRIMARY KEY (message_id, user_id)
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_messages_sender ON messages (sender_id);
CREATE INDEX IF NOT EXISTS idx_conversations_participants ON conversations (participant_ids);

