package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type HomeworkRepository struct {
	store *Store
}

func (r *HomeworkRepository) Create(h *model.Homework) error {
	if err := h.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO homeworks (title) VALUES ($1) RETURNING id",
		h.Title,
	).Scan(
		&h.ID,
	)
}

func (r *HomeworkRepository) FindByTitle(title string) (*model.Homework, error) {
	h := &model.Homework{}
	if err := r.store.db.QueryRow(
		"SELECT id, title FROM homeworks WHERE title = $1",
		title,
	).Scan(
		&h.ID,
		&h.Title,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return h, nil
}
