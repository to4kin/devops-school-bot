CREATE TABLE lesson (
    id bigserial NOT NULL PRIMARY KEY,
    title VARCHAR NOT NULL UNIQUE,
    module_id BIGINT NOT NULL
);