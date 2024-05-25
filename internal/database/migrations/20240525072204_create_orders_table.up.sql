CREATE TABLE IF NOT EXISTS orders (
    id varchar PRIMARY KEY,
    created_at bigint NOT NULL,
    merchant_id varchar NOT NULL,

    CONSTRAINT merchant_id_orders_merchants_fk 
        FOREIGN KEY (merchant_id) 
        REFERENCES merchants (id)
);
