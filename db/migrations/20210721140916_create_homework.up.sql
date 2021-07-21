CREATE TABLE homework (
    id bigserial NOT NULL PRIMARY KEY,
    student_id bigint NOT NULL,
    lesson_id bigint NOT NULL,
    accept BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE (student_id, lesson_id),
    FOREIGN KEY (student_id) REFERENCES student(id) ON DELETE RESTRICT,
    FOREIGN KEY (lesson_id) REFERENCES lesson(id) ON DELETE RESTRICT
);