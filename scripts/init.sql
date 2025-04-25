CREATE DATABASE  payments_events;

CREATE TABLE IF NOT EXISTS payment_events (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    amount INTEGER NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
