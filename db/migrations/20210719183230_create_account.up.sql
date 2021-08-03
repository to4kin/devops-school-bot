CREATE TABLE account (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created TIMESTAMP NOT NULL,
    telegram_id BIGINT NOT NULL UNIQUE,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR,
    username VARCHAR,
    superuser BOOLEAN NOT NULL DEFAULT FALSE
);