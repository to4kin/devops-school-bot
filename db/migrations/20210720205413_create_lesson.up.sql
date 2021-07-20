CREATE TABLE lesson (
    id bigserial not null primary key,
    title varchar not null unique
);