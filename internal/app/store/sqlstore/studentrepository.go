package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// StudentRepository ...
type StudentRepository struct {
	store *Store
}

// Create ...
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

// FindByAccountIDSchoolID ...
func (r *StudentRepository) FindByAccountIDSchoolID(accountID int64, schoolID int64) (*model.Student, error) {
	s := &model.Student{
		Account: &model.Account{},
		School:  &model.School{},
	}
	if err := r.store.db.QueryRow(`
		SELECT st.id, st.active, 
			acc.id, acc.telegram_id, acc.first_name, acc.last_name, acc.username, acc.superuser,
			sch.id, sch.title, sch.active, sch.finished
		FROM student st 
		JOIN account acc ON acc.id = st.account_id
		JOIN school sch ON sch.id = st.school_id
		WHERE st.account_id = $1 AND st.school_id = $2
		`,
		accountID,
		schoolID,
	).Scan(
		&s.ID,
		&s.Active,
		&s.Account.ID,
		&s.Account.TelegramID,
		&s.Account.FirstName,
		&s.Account.LastName,
		&s.Account.Username,
		&s.Account.Superuser,
		&s.School.ID,
		&s.School.Title,
		&s.School.Active,
		&s.School.Finished,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return s, nil
}
