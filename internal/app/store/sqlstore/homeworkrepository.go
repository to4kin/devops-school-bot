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

func (r *HomeworkRepository) FindByStudentID(student_id int64) ([]*model.Homework, error) {
	rowsCount := 0
	hw := []*model.Homework{}

	rows, err := r.store.db.Query(`
		SELECT hw.id, hw.message_id, hw.verify,
			st.id, st.active,
			acc.id, acc.telegram_id, acc.first_name, acc.last_name, acc.username, acc.superuser,
			sch.id, sch.title, sch.active, sch.finished,
			les.id, les.title
		FROM homework hw
		JOIN student st ON st.id = hw.student_id
		JOIN account acc ON acc.id = st.account_id
		JOIN school sch ON sch.id = st.school_id
		JOIN lesson les ON les.id = hw.lesson_id
		WHERE hw.student_id = $1
		`,
		student_id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsCount++

		h := &model.Homework{
			Student: &model.Student{
				Account: &model.Account{},
				School:  &model.School{},
			},
			Lesson: &model.Lesson{},
		}

		if err := rows.Scan(
			&h.ID,
			&h.MessageID,
			&h.Verify,
			&h.Student.ID,
			&h.Student.Active,
			&h.Student.Account.ID,
			&h.Student.Account.TelegramID,
			&h.Student.Account.FirstName,
			&h.Student.Account.LastName,
			&h.Student.Account.Username,
			&h.Student.Account.Superuser,
			&h.Student.School.ID,
			&h.Student.School.Title,
			&h.Student.School.Active,
			&h.Student.School.Finished,
			&h.Lesson.ID,
			&h.Lesson.Title,
		); err != nil {
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

func (r *HomeworkRepository) FindBySchoolID(school_id int64) ([]*model.Homework, error) {
	rowsCount := 0
	hw := []*model.Homework{}

	rows, err := r.store.db.Query(`
		SELECT hw.id, hw.message_id, hw.verify,
			st.id, st.active,
			acc.id, acc.telegram_id, acc.first_name, acc.last_name, acc.username, acc.superuser,
			sch.id, sch.title, sch.active, sch.finished,
			les.id, les.title
		FROM homework hw
		JOIN student st ON st.id = hw.student_id
		JOIN account acc ON acc.id = st.account_id
		JOIN school sch ON sch.id = st.school_id
		JOIN lesson les ON les.id = hw.lesson_id
		WHERE sch.id = $1
		`,
		school_id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsCount++

		h := &model.Homework{
			Student: &model.Student{
				Account: &model.Account{},
				School:  &model.School{},
			},
			Lesson: &model.Lesson{},
		}

		if err := rows.Scan(
			&h.ID,
			&h.MessageID,
			&h.Verify,
			&h.Student.ID,
			&h.Student.Active,
			&h.Student.Account.ID,
			&h.Student.Account.TelegramID,
			&h.Student.Account.FirstName,
			&h.Student.Account.LastName,
			&h.Student.Account.Username,
			&h.Student.Account.Superuser,
			&h.Student.School.ID,
			&h.Student.School.Title,
			&h.Student.School.Active,
			&h.Student.School.Finished,
			&h.Lesson.ID,
			&h.Lesson.Title,
		); err != nil {
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

func (r *HomeworkRepository) FindByStudentIDLessonID(student_id int64, lesson_id int64) (*model.Homework, error) {
	h := &model.Homework{
		Student: &model.Student{
			Account: &model.Account{},
			School:  &model.School{},
		},
		Lesson: &model.Lesson{},
	}

	if err := r.store.db.QueryRow(`
		SELECT hw.id, hw.message_id, hw.verify,
			st.id, st.active,
			acc.id, acc.telegram_id, acc.first_name, acc.last_name, acc.username, acc.superuser,
			sch.id, sch.title, sch.active, sch.finished,
			les.id, les.title
		FROM homework hw
		JOIN student st ON st.id = hw.student_id
		JOIN account acc ON acc.id = st.account_id
		JOIN school sch ON sch.id = st.school_id
		JOIN lesson les ON les.id = hw.lesson_id
		WHERE hw.student_id = $1 AND hw.lesson_id = $2
		`,
		student_id,
		lesson_id,
	).Scan(
		&h.ID,
		&h.MessageID,
		&h.Verify,
		&h.Student.ID,
		&h.Student.Active,
		&h.Student.Account.ID,
		&h.Student.Account.TelegramID,
		&h.Student.Account.FirstName,
		&h.Student.Account.LastName,
		&h.Student.Account.Username,
		&h.Student.Account.Superuser,
		&h.Student.School.ID,
		&h.Student.School.Title,
		&h.Student.School.Active,
		&h.Student.School.Finished,
		&h.Lesson.ID,
		&h.Lesson.Title,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return h, nil
}
