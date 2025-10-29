CREATE TABLE analytics_events (
    id BIGSERIAL PRIMARY KEY,
    tg_user_id BIGINT NOT NULL,
    event_name TEXT NOT NULL,
    event_timestamp TIMESTAMP NOT NULL
);