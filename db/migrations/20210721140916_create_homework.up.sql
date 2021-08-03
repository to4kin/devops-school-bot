CREATE TABLE homework (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created TIMESTAMP NOT NULL,
    student_id BIGINT NOT NULL,
    lesson_id BIGINT NOT NULL,
    message_id BIGINT NOT NULL,
    verify BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE (student_id, lesson_id),
    FOREIGN KEY (student_id) REFERENCES student(id) ON DELETE RESTRICT,
    FOREIGN KEY (lesson_id) REFERENCES lesson(id) ON DELETE RESTRICT
);