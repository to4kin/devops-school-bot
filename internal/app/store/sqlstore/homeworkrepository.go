package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// HomeworkRepository ...
type HomeworkRepository struct {
	store *Store
}

// Create ...
func (r *HomeworkRepository) Create(h *model.Homework) error {
	if err := h.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO homework (created, student_id, lesson_id, message_id, verify, active) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		h.Created,
		h.Student.ID,
		h.Lesson.ID,
		h.MessageID,
		h.Verify,
		h.Active,
	).Scan(
		&h.ID,
	)
}

// Update ...
func (r *HomeworkRepository) Update(h *model.Homework) error {
	if err := r.store.db.QueryRow(
		"UPDATE homework SET active = $2, verify = $3 WHERE id = $1 RETURNING id",
		h.ID,
		h.Active,
		h.Verify,
	).Scan(
		&h.ID,
	); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}

		return err
	}

	return nil
}

// DeleteByMessageID ...
func (r *HomeworkRepository) DeleteByMessageID(messageID int64) error {
	res, err := r.store.db.Exec(
		"DELETE FROM homework WHERE message_id = $1",
		messageID,
	)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return store.ErrRecordNotFound
	}

	return nil
}

// FindByID ...
func (r *HomeworkRepository) FindByID(id int64) (*model.Homework, error) {
	h := &model.Homework{
		Student: &model.Student{
			Account: &model.Account{},
			School:  &model.School{},
		},
		Lesson: &model.Lesson{
			Module: &model.Module{},
		},
	}

	if err := r.store.db.QueryRow(`
		SELECT hw.id, hw.created, hw.message_id, hw.verify, hw.active,
			st.id, st.created, st.active, st.full_course,
			acc.id, acc.created, acc.telegram_id, acc.first_name, acc.last_name, acc.username, acc.superuser,
			sch.id, sch.created, sch.title, sch.chat_id, sch.active,
			les.id, les.title,
			mod.id, mod.title
		FROM homework hw
		JOIN student st ON st.id = hw.student_id
		JOIN account acc ON acc.id = st.account_id
		JOIN school sch ON sch.id = st.school_id
		JOIN lesson les ON les.id = hw.lesson_id
		JOIN module mod ON mod.id = les.module_id
		WHERE hw.id = $1
		`,
		id,
	).Scan(
		&h.ID,
		&h.Created,
		&h.MessageID,
		&h.Verify,
		&h.Active,
		&h.Student.ID,
		&h.Student.Created,
		&h.Student.Active,
		&h.Student.FullCourse,
		&h.Student.Account.ID,
		&h.Student.Account.Created,
		&h.Student.Account.TelegramID,
		&h.Student.Account.FirstName,
		&h.Student.Account.LastName,
		&h.Student.Account.Username,
		&h.Student.Account.Superuser,
		&h.Student.School.ID,
		&h.Student.Account.Created,
		&h.Student.School.Title,
		&h.Student.School.ChatID,
		&h.Student.School.Active,
		&h.Lesson.ID,
		&h.Lesson.Title,
		&h.Lesson.Module.ID,
		&h.Lesson.Module.Title,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return h, nil
}

// FindByStudentID ...
func (r *HomeworkRepository) FindByStudentID(studentID int64) ([]*model.Homework, error) {
	rowsCount := 0
	hw := []*model.Homework{}

	rows, err := r.store.db.Query(`
		SELECT hw.id, hw.created, hw.message_id, hw.verify, hw.active,
			st.id, st.created, st.active, st.full_course,
			acc.id, acc.created, acc.telegram_id, acc.first_name, acc.last_name, acc.username, acc.superuser,
			sch.id, sch.created, sch.title, sch.chat_id, sch.active,
			les.id, les.title,
			mod.id, mod.title
		FROM homework hw
		JOIN student st ON st.id = hw.student_id
		JOIN account acc ON acc.id = st.account_id
		JOIN school sch ON sch.id = st.school_id
		JOIN lesson les ON les.id = hw.lesson_id
		JOIN module mod ON mod.id = les.module_id
		WHERE hw.student_id = $1
		ORDER BY les.title ASC
		`,
		studentID,
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
			Lesson: &model.Lesson{
				Module: &model.Module{},
			},
		}

		if err := rows.Scan(
			&h.ID,
			&h.Created,
			&h.MessageID,
			&h.Verify,
			&h.Active,
			&h.Student.ID,
			&h.Student.Created,
			&h.Student.Active,
			&h.Student.FullCourse,
			&h.Student.Account.ID,
			&h.Student.Account.Created,
			&h.Student.Account.TelegramID,
			&h.Student.Account.FirstName,
			&h.Student.Account.LastName,
			&h.Student.Account.Username,
			&h.Student.Account.Superuser,
			&h.Student.School.ID,
			&h.Student.School.Created,
			&h.Student.School.Title,
			&h.Student.School.ChatID,
			&h.Student.School.Active,
			&h.Lesson.ID,
			&h.Lesson.Title,
			&h.Lesson.Module.ID,
			&h.Lesson.Module.Title,
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

// FindBySchoolID ...
func (r *HomeworkRepository) FindBySchoolID(schoolID int64) ([]*model.Homework, error) {
	rowsCount := 0
	hw := []*model.Homework{}

	rows, err := r.store.db.Query(`
		SELECT hw.id, hw.created, hw.message_id, hw.verify, hw.active,
			st.id, st.created, st.active, st.full_course,
			acc.id, acc.created, acc.telegram_id, acc.first_name, acc.last_name, acc.username, acc.superuser,
			sch.id, sch.created, sch.title, sch.chat_id, sch.active,
			les.id, les.title,
			mod.id, mod.title
		FROM homework hw
		JOIN student st ON st.id = hw.student_id
		JOIN account acc ON acc.id = st.account_id
		JOIN school sch ON sch.id = st.school_id
		JOIN lesson les ON les.id = hw.lesson_id
		JOIN module mod ON mod.id = les.module_id
		WHERE sch.id = $1
		ORDER BY les.title ASC
		`,
		schoolID,
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
			Lesson: &model.Lesson{
				Module: &model.Module{},
			},
		}

		if err := rows.Scan(
			&h.ID,
			&h.Created,
			&h.MessageID,
			&h.Verify,
			&h.Active,
			&h.Student.ID,
			&h.Student.Created,
			&h.Student.Active,
			&h.Student.FullCourse,
			&h.Student.Account.ID,
			&h.Student.Account.Created,
			&h.Student.Account.TelegramID,
			&h.Student.Account.FirstName,
			&h.Student.Account.LastName,
			&h.Student.Account.Username,
			&h.Student.Account.Superuser,
			&h.Student.School.ID,
			&h.Student.Account.Created,
			&h.Student.School.Title,
			&h.Student.School.ChatID,
			&h.Student.School.Active,
			&h.Lesson.ID,
			&h.Lesson.Title,
			&h.Lesson.Module.ID,
			&h.Lesson.Module.Title,
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

// FindBySchoolIDLessonID ...
func (r *HomeworkRepository) FindBySchoolIDLessonID(schoolID int64, lessonID int64) ([]*model.Homework, error) {
	rowsCount := 0
	hw := []*model.Homework{}

	rows, err := r.store.db.Query(`
		SELECT hw.id, hw.created, hw.message_id, hw.verify, hw.active,
			st.id, st.created, st.active, st.full_course,
			acc.id, acc.created, acc.telegram_id, acc.first_name, acc.last_name, acc.username, acc.superuser,
			sch.id, sch.created, sch.title, sch.chat_id, sch.active,
			les.id, les.title,
			mod.id, mod.title
		FROM homework hw
		JOIN student st ON st.id = hw.student_id
		JOIN account acc ON acc.id = st.account_id
		JOIN school sch ON sch.id = st.school_id
		JOIN lesson les ON les.id = hw.lesson_id
		JOIN module mod ON mod.id = les.module_id
		WHERE les.id = $1 AND sch.id = $2
		ORDER BY les.title ASC
		`,
		lessonID,
		schoolID,
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
			Lesson: &model.Lesson{
				Module: &model.Module{},
			},
		}

		if err := rows.Scan(
			&h.ID,
			&h.Created,
			&h.MessageID,
			&h.Verify,
			&h.Active,
			&h.Student.ID,
			&h.Student.Created,
			&h.Student.Active,
			&h.Student.FullCourse,
			&h.Student.Account.ID,
			&h.Student.Account.Created,
			&h.Student.Account.TelegramID,
			&h.Student.Account.FirstName,
			&h.Student.Account.LastName,
			&h.Student.Account.Username,
			&h.Student.Account.Superuser,
			&h.Student.School.ID,
			&h.Student.Account.Created,
			&h.Student.School.Title,
			&h.Student.School.ChatID,
			&h.Student.School.Active,
			&h.Lesson.ID,
			&h.Lesson.Title,
			&h.Lesson.Module.ID,
			&h.Lesson.Module.Title,
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

// FindByStudentIDLessonID ...
func (r *HomeworkRepository) FindByStudentIDLessonID(studentID int64, lessonID int64) (*model.Homework, error) {
	h := &model.Homework{
		Student: &model.Student{
			Account: &model.Account{},
			School:  &model.School{},
		},
		Lesson: &model.Lesson{
			Module: &model.Module{},
		},
	}

	if err := r.store.db.QueryRow(`
		SELECT hw.id, hw.created, hw.message_id, hw.verify, hw.active,
			st.id, st.created, st.active, st.full_course,
			acc.id, acc.created, acc.telegram_id, acc.first_name, acc.last_name, acc.username, acc.superuser,
			sch.id, sch.created, sch.title, sch.chat_id, sch.active,
			les.id, les.title,
			mod.id, mod.title
		FROM homework hw
		JOIN student st ON st.id = hw.student_id
		JOIN account acc ON acc.id = st.account_id
		JOIN school sch ON sch.id = st.school_id
		JOIN lesson les ON les.id = hw.lesson_id
		JOIN module mod ON mod.id = les.module_id
		WHERE hw.student_id = $1 AND hw.lesson_id = $2
		`,
		studentID,
		lessonID,
	).Scan(
		&h.ID,
		&h.Created,
		&h.MessageID,
		&h.Verify,
		&h.Active,
		&h.Student.ID,
		&h.Student.Created,
		&h.Student.Active,
		&h.Student.FullCourse,
		&h.Student.Account.ID,
		&h.Student.Account.Created,
		&h.Student.Account.TelegramID,
		&h.Student.Account.FirstName,
		&h.Student.Account.LastName,
		&h.Student.Account.Username,
		&h.Student.Account.Superuser,
		&h.Student.School.ID,
		&h.Student.Account.Created,
		&h.Student.School.Title,
		&h.Student.School.ChatID,
		&h.Student.School.Active,
		&h.Lesson.ID,
		&h.Lesson.Title,
		&h.Lesson.Module.ID,
		&h.Lesson.Module.Title,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return h, nil
}
