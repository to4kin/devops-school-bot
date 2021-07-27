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
		"INSERT INTO homework (student_id, lesson_id, message_id, verify) VALUES ($1, $2, $3, $4) RETURNING id",
		h.Student.ID,
		h.Lesson.ID,
		h.MessageID,
		h.Verify,
	).Scan(
		&h.ID,
	)
}

func (r *HomeworkRepository) FindByStudent(student *model.Student) ([]*model.Homework, error) {
	rowsCount := 0
	hw := []*model.Homework{}

	rows, err := r.store.db.Query(
		"SELECT id, lesson_id, message_id, verify FROM homework WHERE student_id = $1",
		student.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsCount++

		h := &model.Homework{
			Student: student,
			Lesson:  &model.Lesson{},
		}

		if err := rows.Scan(&h.ID, &h.Lesson.ID, &h.MessageID, &h.Verify); err != nil {
			return nil, err
		}

		hw = append(hw, h)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if rowsCount == 0 {
		return nil, store.ErrRecordNotFound
	}

	return hw, nil
}

func (r *HomeworkRepository) FindBySchool(school *model.School) ([]*model.Homework, error) {
	rowsCount := 0
	hw := []*model.Homework{}

	rows, err := r.store.db.Query(
		"SELECT id, student_id, lesson_id, message_id, verify FROM homework WHERE student_id IN (SELECT id FROM student WHERE school_id = $1)",
		school.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsCount++

		h := &model.Homework{
			Student: &model.Student{},
			Lesson:  &model.Lesson{},
		}

		if err := rows.Scan(&h.ID, &h.Student.ID, &h.Lesson.ID, &h.MessageID, &h.Verify); err != nil {
			return nil, err
		}

		hw = append(hw, h)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if rowsCount == 0 {
		return nil, store.ErrRecordNotFound
	}

	return hw, nil
}

func (r *HomeworkRepository) FindByStudentLesson(student *model.Student, lesson *model.Lesson) (*model.Homework, error) {
	h := &model.Homework{
		Student: student,
		Lesson:  lesson,
	}
	if err := r.store.db.QueryRow(
		"SELECT id, message_id, verify FROM homework WHERE student_id = $1 AND lesson_id = $2",
		h.Student.ID,
		h.Lesson.ID,
	).Scan(
		&h.ID,
		&h.MessageID,
		&h.Verify,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return h, nil
}
