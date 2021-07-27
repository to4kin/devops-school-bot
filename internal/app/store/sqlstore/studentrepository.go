package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type StudentRepository struct {
	store *Store
}

func (r *StudentRepository) Create(s *model.Student) error {
	if err := s.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO student (account_id, school_id, active) VALUES ($1, $2, $3) RETURNING id",
		s.Account.ID,
		s.School.ID,
		s.Active,
	).Scan(
		&s.ID,
	)
}

func (r *StudentRepository) FindByAccountSchool(account *model.Account, school *model.School) (*model.Student, error) {
	s := &model.Student{
		Account: account,
		School:  school,
	}
	if err := r.store.db.QueryRow(
		"SELECT id, active FROM student WHERE account_id = $1 AND school_id = $2",
		s.Account.ID,
		s.School.ID,
	).Scan(
		&s.ID,
		&s.Active,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return s, nil
}
