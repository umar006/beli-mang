CREATE TABLE IF NOT EXISTS merchant_items (
    id varchar PRIMARY KEY,
    created_at bigint NOT NULL,
    name varchar NOT NULL,
    price numeric NOT NULL,
    category varchar NOT NULL,
    image_url varchar NOT NULL,

    CONSTRAINT category_merchant_items_check CHECK (
        category IN ('Beverage', 'Food', 'Snack', 'Condiments', 'Additions')
    )
);
