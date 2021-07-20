package teststore

import (
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[int64]*model.User
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	r.users[u.TelegramID] = u
	u.ID = int64(len(r.users))

	return nil
}

func (r *UserRepository) FindByTelegramID(telegramID int64) (*model.User, error) {
	u, ok := r.users[telegramID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}
