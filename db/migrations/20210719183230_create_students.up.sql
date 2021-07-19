CREATE TABLE students (
  id bigserial not null primary key,
  telegram_id bigint not null unique,
  first_name varchar not null,
  last_name varchar,
  username varchar
);