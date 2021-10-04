package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// ModuleRepository ...
type ModuleRepository struct {
	store *Store
}

// Create ...
func (r *ModuleRepository) Create(m *model.Module) error {
	if err := m.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO module (title) VALUES ($1) RETURNING id",
		m.Title,
	).Scan(
		&m.ID,
	)
}

// FindAll ...
func (r *ModuleRepository) FindAll() ([]*model.Module, error) {
	rowsCount := 0
	modules := []*model.Module{}

	rows, err := r.store.db.Query(`
		SELECT id, title FROM module ORDER BY title ASC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsCount++

		m := &model.Module{}

		if err := rows.Scan(
			&m.ID,
			&m.Title,
		); err != nil {
			return nil, err
		}

		modules = append(modules, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if rowsCount == 0 {
		return nil, store.ErrRecordNotFound
	}

	return modules, nil
}

// FindByID ...
func (r *ModuleRepository) FindByID(moduleID int64) (*model.Module, error) {
	m := &model.Module{}
	if err := r.store.db.QueryRow(
		"SELECT id, title FROM module WHERE id = $1",
		moduleID,
	).Scan(
		&m.ID,
		&m.Title,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return m, nil
}

// FindByTitle ...
func (r *ModuleRepository) FindByTitle(title string) (*model.Module, error) {
	m := &model.Module{}
	if err := r.store.db.QueryRow(
		"SELECT id, title FROM module WHERE title = $1",
		title,
	).Scan(
		&m.ID,
		&m.Title,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return m, nil
}

// FindBySchoolID ...
func (r *ModuleRepository) FindBySchoolID(schoolID int64) ([]*model.Module, error) {
	rowsCount := 0
	m := []*model.Module{}

	rows, err := r.store.db.Query(`
		SELECT module.id, module.title FROM module
		JOIN lesson ON lesson.module_id = module.id
		JOIN homework ON homework.lesson_id = lesson.id
		JOIN student ON student.id = homework.student_id
		WHERE student.school_id = $1 AND homework.active = true
		GROUP BY module.id, module.title
		ORDER BY module.title ASC
		`,
		schoolID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsCount++
		module := &model.Module{}

		if err := rows.Scan(&module.ID, &module.Title); err != nil {
			return nil, err
		}

		m = append(m, module)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if rowsCount == 0 {
		return nil, store.ErrRecordNotFound
	}

	return m, nil
}
