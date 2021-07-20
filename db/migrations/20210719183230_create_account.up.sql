CREATE TABLE account (
    id bigserial not null primary key,
    telegram_id bigint not null unique,
    first_name varchar not null,
    last_name varchar,
    username varchar,
    is_admin boolean not null default false
);