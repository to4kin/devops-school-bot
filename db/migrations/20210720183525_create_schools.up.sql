CREATE TABLE schools (
    id bigserial not null primary key,
    title varchar not null unique,
    in_progress boolean not null default false,
    finished boolean not null default false
);