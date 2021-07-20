CREATE TABLE homeworks (
    id bigserial not null primary key,
    title varchar not null unique
);