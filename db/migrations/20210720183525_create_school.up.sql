CREATE TABLE school (
    id bigserial not null primary key,
    title varchar not null unique,
    active boolean not null default false,
    finished boolean not null default false
);