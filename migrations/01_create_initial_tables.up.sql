CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TYPE role AS ENUM ('admin', 'user');
CREATE TYPE status AS ENUM ('pending', 'accepted');

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users
(
    user_id    UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    first_name VARCHAR(32)              NOT NULL CHECK ( first_name <> '' ),
    last_name  VARCHAR(32)              NOT NULL CHECK ( last_name <> '' ),
    email      VARCHAR(64) UNIQUE       NOT NULL CHECK ( email <> '' ),
    avatar     VARCHAR(250),
    password   VARCHAR(250)             NOT NULL CHECK ( octet_length(password) <> 0 ),
    role       role                     NOT NULL DEFAULT 'user',
    delivery_address      VARCHAR(250)  NOT NULL CHECK ( delivery_address <> '' ),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE          DEFAULT CURRENT_TIMESTAMP
);


DROP TABLE IF EXISTS brands CASCADE;
CREATE TABLE brands
(
    brand_id    UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    brand_name  VARCHAR(32)           NOT NULL CHECK ( first_name <> '' ),
    logo        VARCHAR(250),
    pickup_address      VARCHAR(250)  NOT NULL CHECK ( pickup_address <> '' ),

    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS products CASCADE;
CREATE TABLE products
(
    product_id  UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    name        VARCHAR(250)  NOT NULL CHECK ( name <> '' ),
    description VARCHAR(5000) NOT NULL CHECK ( description <> '' ),
    price       NUMERIC       NOT NULL,
    brand_id    UUID REFERENCES brands (brand_id),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_products__brand_id ON products(brand_id);

DROP TABLE IF EXISTS orders CASCADE;
CREATE TABLE orders
(
    order_id    UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    user_id     UUID REFERENCES users (user_id),
    brand_id    UUID REFERENCES brands (brand_id),
    item        JSONB,
    quantity    NUMERIC       NOT NULL,
    total_price NUMERIC       NOT NULL,
    status      status        NOT NULL DEFAULT 'pending',

    delivery_source_address      VARCHAR(250)  NOT NULL CHECK ( delivery_source_address <> '' ),
    delivery_destination_address VARCHAR(250)  NOT NULL CHECK ( delivery_destination_address <> '' ),

    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_orders__user_id ON orders(user_id);
CREATE INDEX idx_orders__brand_id ON orders(brand_id);