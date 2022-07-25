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

DROP TABLE IF EXISTS sellers CASCADE;
CREATE TABLE sellers
(
    seller_id  UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    first_name VARCHAR(32)              NOT NULL CHECK ( first_name <> '' ),
    last_name  VARCHAR(32)              NOT NULL CHECK ( last_name <> '' ),
    email      VARCHAR(64) UNIQUE       NOT NULL CHECK ( email <> '' ),
    avatar     VARCHAR(250),
    password   VARCHAR(250)             NOT NULL CHECK ( octet_length(password) <> 0 ),
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
    seller_id   UUID REFERENCES sellers (seller_id),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS orders CASCADE;
CREATE TABLE orders
(
    order_id    UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    user_id     UUID REFERENCES users (user_id),
    seller_id   UUID REFERENCES sellers (seller_id),
    item        JSONB,
    quantity    NUMERIC       NOT NULL,
    total_price NUMERIC       NOT NULL,
    status      status        NOT NULL DEFAULT 'pending',

    delivery_source_address      VARCHAR(250)  NOT NULL CHECK ( delivery_source_address <> '' ),
    delivery_destination_address VARCHAR(250)  NOT NULL CHECK ( delivery_destination_address <> '' ),

    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);