CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       user_id BIGINT NOT NULL,
                       task TEXT NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       completed BOOLEAN DEFAULT FALSE
);