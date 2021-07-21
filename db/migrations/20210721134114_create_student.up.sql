CREATE TABLE student (
    id bigserial NOT NULL PRIMARY KEY,
    account_id bigint NOT NULL,
    school_id bigint NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (account_id, school_id),
    FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE RESTRICT,
    FOREIGN KEY (school_id) REFERENCES school(id) ON DELETE RESTRICT
);