package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type SchoolRepository struct {
	store *Store
}

func (r *SchoolRepository) Create(s *model.School) error {
	if err := s.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO school (title, active, finished) VALUES ($1, $2, $3) RETURNING id",
		s.Title,
		s.Active,
		s.Finished,
	).Scan(
		&s.ID,
	)
}

func (r *SchoolRepository) FindByTitle(title string) (*model.School, error) {
	s := &model.School{}
	if err := r.store.db.QueryRow(
		"SELECT id, title, active, finished FROM school WHERE title = $1",
		title,
	).Scan(
		&s.ID,
		&s.Title,
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
