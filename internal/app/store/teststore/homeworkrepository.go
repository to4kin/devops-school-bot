package teststore

import (
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type HomeworkRepository struct {
	store     *Store
	homeworks map[string]*model.Homework
}

func (r *HomeworkRepository) Create(h *model.Homework) error {
	if err := h.Validate(); err != nil {
		return err
	}

	r.homeworks[h.Title] = h
	h.ID = int64(len(r.homeworks))

	return nil
}

func (r *HomeworkRepository) FindByTitle(title string) (*model.Homework, error) {
	h, ok := r.homeworks[title]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return h, nil
}
