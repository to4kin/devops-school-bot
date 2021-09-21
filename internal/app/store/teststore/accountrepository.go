package teststore

import (
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// AccountRepository ...
type AccountRepository struct {
	store    *Store
	accounts []*model.Account
}

// Create ...
func (r *AccountRepository) Create(a *model.Account) error {
	if err := a.Validate(); err != nil {
		return err
	}

	account, err := r.store.accountRepository.FindByTelegramID(a.TelegramID)
	if err != nil && err != store.ErrRecordNotFound {
		return err
	}

	if account != nil {
		return store.ErrRecordIsExist
	}

	a.ID = int64(len(r.accounts) + 1)
	r.accounts = append(r.accounts, a)
	return nil
}

// Update ...
func (r *AccountRepository) Update(a *model.Account) error {
	for _, account := range r.accounts {
		if account.TelegramID == a.TelegramID {
			account.FirstName = a.FirstName
			account.LastName = a.LastName
			account.Username = a.Username
			account.Superuser = a.Superuser
			return nil
		}
	}

	return store.ErrRecordNotFound
}

// FindAll ...
func (r *AccountRepository) FindAll() ([]*model.Account, error) {
	if len(r.accounts) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return r.accounts, nil
}

// FindByID ...
func (r *AccountRepository) FindByID(id int64) (*model.Account, error) {
	for _, account := range r.accounts {
		if account.ID == id {
			return account, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

// FindByTelegramID ...
func (r *AccountRepository) FindByTelegramID(telegramID int64) (*model.Account, error) {
	for _, account := range r.accounts {
		if account.TelegramID == telegramID {
			return account, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

// FindBySuperuser ...
func (r *AccountRepository) FindBySuperuser(superuser bool) ([]*model.Account, error) {
	accounts := []*model.Account{}

	for _, account := range r.accounts {
		if account.Superuser == superuser {
			accounts = append(accounts, account)
		}
	}

	if len(accounts) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return accounts, nil
}
