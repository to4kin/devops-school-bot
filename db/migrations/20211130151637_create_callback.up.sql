CREATE TABLE callback (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  created TIMESTAMP NOT NULL,
  type VARCHAR NOT NULL,
  type_id BIGINT NOT NULL,
  command VARCHAR NOT NULL,
  list_command VARCHAR NOT NULL,
  UNIQUE (type, type_id, command, list_command)
);
