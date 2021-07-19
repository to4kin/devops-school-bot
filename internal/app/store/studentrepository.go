package store

import "gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"

type StudentRepository struct {
	store *Store
}

// Create ...
func (r *StudentRepository) Create(s *model.Student) (*model.Student, error) {
	if err := r.store.db.QueryRow(
		"INSERT INTO students (telegram_id, first_name, last_name, username) VALUES ($1, $2, $3, $4) RETURNING id",
		s.TelegramID,
		s.FirstName,
		s.LastName,
		s.Username,
	).Scan(
		&s.ID,
	); err != nil {
		return nil, err
	}

	return s, nil
}

// FindByTelegramID ...
func (r *StudentRepository) FindByTelegramID(telegramID int) (*model.Student, error) {
	s := &model.Student{}
	if err := r.store.db.QueryRow(
		"SELECT id, telegram_id, first_name, last_name, username FROM students WHERE telegram_id = $1",
		telegramID,
	).Scan(
		&s.ID,
		&s.TelegramID,
		&s.FirstName,
		&s.LastName,
		&s.Username,
	); err != nil {
		return nil, err
	}

	return s, nil
}
