package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type HomeworkRepository struct {
	store *Store
}

func (r *HomeworkRepository) Create(h *model.Homework) error {
	if err := h.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO homework (student_id, lesson_id, accept) VALUES ($1, $2, $3) RETURNING id",
		h.Student.ID,
		h.Lesson.ID,
		h.Accept,
	).Scan(
		&h.ID,
	)
}

func (r *HomeworkRepository) FindByStudentIDLessonID(studentID int64, lessonID int64) (*model.Homework, error) {
	h := &model.Homework{}
	if err := r.store.db.QueryRow(
		"SELECT id, accept FROM homework WHERE student_id = $1 AND lesson_id = $2",
		studentID,
		lessonID,
	).Scan(
		&h.ID,
		&h.Accept,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return h, nil
}
