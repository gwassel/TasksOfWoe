-- Handler performance metrics table
CREATE TABLE handler_performance_metrics (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    handler_name TEXT NOT NULL,
    command TEXT NOT NULL,
    duration_ms BIGINT NOT NULL,
    timestamp TIMESTAMP NOT NULL
);

-- Indexes for efficient querying
CREATE INDEX idx_handler_timestamp ON handler_performance_metrics(handler_name, timestamp);
CREATE INDEX idx_user_timestamp ON handler_performance_metrics(user_id, timestamp);
CREATE INDEX idx_timestamp ON handler_performance_metrics(timestamp DESC);

-- Admin users table
CREATE TABLE admin_users (
    id BIGSERIAL PRIMARY KEY,
    telegram_user_id BIGINT UNIQUE NOT NULL,
    registered_at TIMESTAMP DEFAULT NOW()
);

-- Sample admin user (Telegram user ID from existing config)
INSERT INTO admin_users (telegram_user_id) VALUES (351083864);
