CREATE TABLE account (
    id bigserial NOT NULL PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR,
    username VARCHAR,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE
);