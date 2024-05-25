BEGIN;

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS merchants (
    id varchar PRIMARY KEY,
    created_at bigint NOT NULL,
    category varchar NOT NULL,
    image_url varchar NOT NULL,
    location geometry NOT NULL,

    CONSTRAINT category_merchants_check CHECK (
        category IN (
            'SmallRestaurant', 'MediumRestaurant',
            'LargeRestaurant', 'MerchandiseRestaurant',
            'BoothKiosk', 'ConvenienceStore'
        )
    )
);

COMMIT;
