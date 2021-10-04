package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// LessonRepository ...
type LessonRepository struct {
	store *Store
}

// Create ...
func (r *LessonRepository) Create(l *model.Lesson) error {
	if err := l.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO lesson (title, module_id) VALUES ($1, $2) RETURNING id",
		l.Title,
		l.Module.ID,
	).Scan(
		&l.ID,
	)
}

// FindAll ...
func (r *LessonRepository) FindAll() ([]*model.Lesson, error) {
	rowsCount := 0
	lessons := []*model.Lesson{}

	rows, err := r.store.db.Query(`
		SELECT les.id, les.title, mod.id, mod.title
		FROM lesson les
		JOIN module mod ON mod.id = les.module_id
		ORDER BY les.title ASC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsCount++

		l := &model.Lesson{
			Module: &model.Module{},
		}

		if err := rows.Scan(
			&l.ID,
			&l.Title,
			&l.Module.ID,
			&l.Module.Title,
		); err != nil {
			return nil, err
		}

		lessons = append(lessons, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if rowsCount == 0 {
		return nil, store.ErrRecordNotFound
	}

	return lessons, nil
}

// FindByID ...
func (r *LessonRepository) FindByID(lessonID int64) (*model.Lesson, error) {
	l := &model.Lesson{
		Module: &model.Module{},
	}
	if err := r.store.db.QueryRow(`
		SELECT les.id, les.title, mod.id, mod.title
		FROM lesson les
		JOIN module mod ON mod.id = les.module_id
		WHERE les.id = $1`,
		lessonID,
	).Scan(
		&l.ID,
		&l.Title,
		&l.Module.ID,
		&l.Module.Title,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return l, nil
}

// FindByTitle ...
func (r *LessonRepository) FindByTitle(title string) (*model.Lesson, error) {
	l := &model.Lesson{
		Module: &model.Module{},
	}
	if err := r.store.db.QueryRow(
		`SELECT les.id, les.title, mod.id, mod.title
		FROM lesson les
		JOIN module mod ON mod.id = les.module_id
		WHERE les.title = $1`,
		title,
	).Scan(
		&l.ID,
		&l.Title,
		&l.Module.ID,
		&l.Module.Title,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return l, nil
}

// FindBySchoolID ...
func (r *LessonRepository) FindBySchoolID(schoolID int64) ([]*model.Lesson, error) {
	rowsCount := 0
	l := []*model.Lesson{}

	rows, err := r.store.db.Query(`
		SELECT lesson.id, lesson.title, mod.id, mod.title 
		FROM lesson
		JOIN module mod ON mod.id = lesson.module_id
		JOIN homework ON homework.lesson_id = lesson.id
		JOIN student ON student.id = homework.student_id
		WHERE student.school_id = $1 AND homework.active = true
		GROUP BY lesson.id, lesson.title, mod.id, mod.title
		ORDER BY lesson.title ASC
		`,
		schoolID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsCount++
		lesson := &model.Lesson{
			Module: &model.Module{},
		}

		if err := rows.Scan(&lesson.ID, &lesson.Title, &lesson.Module.ID, &lesson.Module.Title); err != nil {
			return nil, err
		}

		l = append(l, lesson)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if rowsCount == 0 {
		return nil, store.ErrRecordNotFound
	}

	return l, nil
}
