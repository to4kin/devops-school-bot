package teststore

import (
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type AccountRepository struct {
	store    *Store
	accounts map[int64]*model.Account
}

func (r *AccountRepository) Create(a *model.Account) error {
	if err := a.Validate(); err != nil {
		return err
	}

	r.accounts[a.TelegramID] = a
	a.ID = int64(len(r.accounts))

	return nil
}

func (r *AccountRepository) FindByTelegramID(telegramID int64) (*model.Account, error) {
	a, ok := r.accounts[telegramID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return a, nil
}
