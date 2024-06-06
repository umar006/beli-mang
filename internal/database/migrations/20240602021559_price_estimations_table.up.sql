CREATE TABLE IF NOT EXISTS price_estimations (
    id varchar PRIMARY KEY,
    created_at bigint NOT NULL,
    total_price int NOT NULL,
    delivery_time_in_minutes int NOT NULL
)
