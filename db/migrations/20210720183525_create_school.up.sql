CREATE TABLE school (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created TIMESTAMP NOT NULL,
    title VARCHAR NOT NULL UNIQUE,
    chat_id BIGINT NOT NULL UNIQUE,
    finished BOOLEAN NOT NULL DEFAULT FALSE
);