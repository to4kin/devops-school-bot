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

	return r.store.db.QueryRow(
		"INSERT INTO school (created, title, chat_id, active) VALUES ($1, $2, $3, $4) RETURNING id",
		s.Created,
		s.Title,
		s.ChatID,
		s.Active,
	).Scan(
		&s.ID,
	)
}

// Update ...
func (r *SchoolRepository) Update(s *model.School) error {
	if err := r.store.db.QueryRow(
		"UPDATE school SET title = $2, active = $3 WHERE id = $1 RETURNING id",
		s.ID,
		s.Title,
		s.Active,
	).Scan(
		&s.ID,
	); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}

		return err
	}

	return nil
}

// FindAll ...
func (r *SchoolRepository) FindAll() ([]*model.School, error) {
	rowsCount := 0
	schools := []*model.School{}

	rows, err := r.store.db.Query(`
		SELECT id, created, title, chat_id, active FROM school ORDER BY created DESC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsCount++

		s := &model.School{}

		if err := rows.Scan(
			&s.ID,
			&s.Created,
			&s.Title,
			&s.ChatID,
			&s.Active,
		); err != nil {
			return nil, err
		}

		schools = append(schools, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if rowsCount == 0 {
		return nil, store.ErrRecordNotFound
	}

	return schools, nil
}

// FindByID ...
func (r *SchoolRepository) FindByID(id int64) (*model.School, error) {
	s := &model.School{}
	if err := r.store.db.QueryRow(
		"SELECT id, created, title, chat_id, active FROM school WHERE id = $1",
		id,
	).Scan(
		&s.ID,
		&s.Created,
		&s.Title,
		&s.ChatID,
		&s.Active,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return s, nil

}

// FindByTitle ...
func (r *SchoolRepository) FindByTitle(title string) (*model.School, error) {
	s := &model.School{}
	if err := r.store.db.QueryRow(
		"SELECT id, created, title, chat_id, active FROM school WHERE title = $1",
		title,
	).Scan(
		&s.ID,
		&s.Created,
		&s.Title,
		&s.ChatID,
		&s.Active,
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
		"SELECT id, created, title, chat_id, active FROM school WHERE chat_id = $1",
		chatID,
	).Scan(
		&s.ID,
		&s.Created,
		&s.Title,
		&s.ChatID,
		&s.Active,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return s, nil

}

// FindByActive ...
func (r *SchoolRepository) FindByActive(active bool) ([]*model.School, error) {
	rowsCount := 0
	schools := []*model.School{}

	rows, err := r.store.db.Query(`
		SELECT id, created, title, chat_id, active FROM school WHERE active = $1 ORDER BY created DESC
		`,
		active,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsCount++

		s := &model.School{}

		if err := rows.Scan(
			&s.ID,
			&s.Created,
			&s.Title,
			&s.ChatID,
			&s.Active,
		); err != nil {
			return nil, err
		}

		schools = append(schools, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if rowsCount == 0 {
		return nil, store.ErrRecordNotFound
	}

	return schools, nil
}
