CREATE TABLE student (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created TIMESTAMP NOT NULL,
    account_id BIGINT NOT NULL,
    school_id BIGINT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (account_id, school_id),
    FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE RESTRICT,
    FOREIGN KEY (school_id) REFERENCES school(id) ON DELETE RESTRICT
);