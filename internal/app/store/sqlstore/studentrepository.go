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
		"INSERT INTO student (created, account_id, school_id, active) VALUES ($1, $2, $3, $4) RETURNING id",
		s.Created,
		s.Account.ID,
		s.School.ID,
		s.Active,
	).Scan(
		&s.ID,
	)
}

// FindBySchoolID ...
func (r *StudentRepository) FindBySchoolID(schoolID int64) ([]*model.Student, error) {
	rowsCount := 0
	students := []*model.Student{}

	rows, err := r.store.db.Query(`
		SELECT st.id, st.created, st.active, 
			acc.id, acc.created, acc.telegram_id, acc.first_name, acc.last_name, acc.username, acc.superuser,
			sch.id, sch.created, sch.title, sch.chat_id, sch.finished
		FROM student st 
		JOIN account acc ON acc.id = st.account_id
		JOIN school sch ON sch.id = st.school_id
		WHERE st.school_id = $1
		`,
		schoolID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsCount++

		s := &model.Student{
			Account: &model.Account{},
			School:  &model.School{},
		}

		if err := rows.Scan(
			&s.ID,
			&s.Created,
			&s.Active,
			&s.Account.ID,
			&s.Account.Created,
			&s.Account.TelegramID,
			&s.Account.FirstName,
			&s.Account.LastName,
			&s.Account.Username,
			&s.Account.Superuser,
			&s.School.ID,
			&s.School.Created,
			&s.School.Title,
			&s.School.ChatID,
			&s.School.Finished,
		); err != nil {
			return nil, err
		}

		students = append(students, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if rowsCount == 0 {
		return nil, store.ErrRecordNotFound
	}

	return students, nil
}

// FindByAccountIDSchoolID ...
func (r *StudentRepository) FindByAccountIDSchoolID(accountID int64, schoolID int64) (*model.Student, error) {
	s := &model.Student{
		Account: &model.Account{},
		School:  &model.School{},
	}
	if err := r.store.db.QueryRow(`
		SELECT st.id, st.created, st.active, 
			acc.id, acc.created, acc.telegram_id, acc.first_name, acc.last_name, acc.username, acc.superuser,
			sch.id, sch.created, sch.title, sch.chat_id, sch.finished
		FROM student st 
		JOIN account acc ON acc.id = st.account_id
		JOIN school sch ON sch.id = st.school_id
		WHERE st.account_id = $1 AND st.school_id = $2
		`,
		accountID,
		schoolID,
	).Scan(
		&s.ID,
		&s.Created,
		&s.Active,
		&s.Account.ID,
		&s.Account.Created,
		&s.Account.TelegramID,
		&s.Account.FirstName,
		&s.Account.LastName,
		&s.Account.Username,
		&s.Account.Superuser,
		&s.School.ID,
		&s.School.Created,
		&s.School.Title,
		&s.School.ChatID,
		&s.School.Finished,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return s, nil
}
