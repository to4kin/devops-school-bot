package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type LessonRepository struct {
	store *Store
}

func (r *LessonRepository) Create(l *model.Leson) error {
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

func (r *LessonRepository) FindByTitle(title string) (*model.Leson, error) {
	l := &model.Leson{}
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
