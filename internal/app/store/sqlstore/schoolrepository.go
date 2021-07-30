package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// SchoolRepository ...
type SchoolRepository struct {
	store *Store
}

// Create ...
func (r *SchoolRepository) Create(s *model.School) error {
	if err := s.Validate(); err != nil {
		return err
	}

	school, err := r.store.schoolRepository.FindActive()
	if err != nil && err != store.ErrRecordNotFound {
		return err
	}

	if school != nil {
		return store.ErrAnotherSchoolIsActive
	}

	return r.store.db.QueryRow(
		"INSERT INTO school (title, chat_id, active, finished) VALUES ($1, $2, $3, $4) RETURNING id",
		s.Title,
		s.ChatID,
		s.Active,
		s.Finished,
	).Scan(
		&s.ID,
	)
}

// FindByTitle ...
func (r *SchoolRepository) FindByTitle(title string) (*model.School, error) {
	s := &model.School{}
	if err := r.store.db.QueryRow(
		"SELECT id, title, chat_id, active, finished FROM school WHERE title = $1",
		title,
	).Scan(
		&s.ID,
		&s.Title,
		&s.ChatID,
		&s.Active,
		&s.Finished,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return s, nil

}

// FindByChatID ...
func (r *SchoolRepository) FindByChatID(chatID int64) (*model.School, error) {
	s := &model.School{}
	if err := r.store.db.QueryRow(
		"SELECT id, title, chat_id, active, finished FROM school WHERE chat_id = $1",
		chatID,
	).Scan(
		&s.ID,
		&s.Title,
		&s.ChatID,
		&s.Active,
		&s.Finished,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return s, nil

}

// FindActive ...
func (r *SchoolRepository) FindActive() (*model.School, error) {
	s := &model.School{}
	if err := r.store.db.QueryRow(
		"SELECT id, title, chat_id, active, finished FROM school WHERE active = true",
	).Scan(
		&s.ID,
		&s.Title,
		&s.ChatID,
		&s.Active,
		&s.Finished,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return s, nil

}
