package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type StudentRepository struct {
	store *Store
}

func (r *StudentRepository) Create(student *model.Student) error {
	if err := student.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO students (telegram_id, first_name, last_name, username) VALUES ($1, $2, $3, $4) RETURNING id",
		student.TelegramID,
		student.FirstName,
		student.LastName,
		student.Username,
	).Scan(
		&student.ID,
	)
}

func (r *StudentRepository) FindByTelegramID(telegramID int64) (*model.Student, error) {
	student := &model.Student{}
	if err := r.store.db.QueryRow(
		"SELECT id, telegram_id, first_name, last_name, username FROM students WHERE telegram_id = $1",
		telegramID,
	).Scan(
		&student.ID,
		&student.TelegramID,
		&student.FirstName,
		&student.LastName,
		&student.Username,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return student, nil
}
