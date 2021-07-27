package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type LessonRepository struct {
	store *Store
}

func (r *LessonRepository) Create(l *model.Lesson) error {
	if err := l.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO lesson (title) VALUES ($1) RETURNING id",
		l.Title,
	).Scan(
		&l.ID,
	)
}

func (r *LessonRepository) FindByID(lesson_id int64) (*model.Lesson, error) {
	l := &model.Lesson{}
	if err := r.store.db.QueryRow(
		"SELECT id, title FROM lesson WHERE id = $1",
		lesson_id,
	).Scan(
		&l.ID,
		&l.Title,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return l, nil
}

func (r *LessonRepository) FindByTitle(title string) (*model.Lesson, error) {
	l := &model.Lesson{}
	if err := r.store.db.QueryRow(
		"SELECT id, title FROM lesson WHERE title = $1",
		title,
	).Scan(
		&l.ID,
		&l.Title,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return l, nil
}

func (r *LessonRepository) FindBySchool(school *model.School) ([]*model.Lesson, error) {
	rowsCount := 0
	l := []*model.Lesson{}

	rows, err := r.store.db.Query(`
		SELECT lesson.id, lesson.title FROM lesson
		JOIN homework ON homework.lesson_id = lesson.id
		JOIN student ON student.id = homework.student_id
		WHERE student.school_id = $1
		GROUP BY lesson.id, lesson.title
		`,
		school.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsCount++
		lesson := &model.Lesson{}

		if err := rows.Scan(&lesson.ID, &lesson.Title); err != nil {
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
