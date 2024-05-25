CREATE TABLE IF NOT EXISTS order_items (
    order_id varchar NOT NULL,
    merchant_item_id varchar NOT NULL,

    CONSTRAINT order_id_order_items_orders_fk 
        FOREIGN KEY (order_id) 
        REFERENCES orders (id),
    CONSTRAINT merchant_item_id_orders_merchant_items_fk 
        FOREIGN KEY (merchant_item_id) 
        REFERENCES merchant_items (id)
);
