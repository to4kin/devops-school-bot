package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type AccountRepository struct {
	store *Store
}

func (r *AccountRepository) Create(a *model.Account) error {
	if err := a.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO account (telegram_id, first_name, last_name, username, superuser) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		a.TelegramID,
		a.FirstName,
		a.LastName,
		a.Username,
		a.Superuser,
	).Scan(
		&a.ID,
	)
}

func (r *AccountRepository) FindByTelegramID(telegramID int64) (*model.Account, error) {
	a := &model.Account{}
	if err := r.store.db.QueryRow(
		"SELECT id, telegram_id, first_name, last_name, username, superuser FROM account WHERE telegram_id = $1",
		telegramID,
	).Scan(
		&a.ID,
		&a.TelegramID,
		&a.FirstName,
		&a.LastName,
		&a.Username,
		&a.Superuser,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return a, nil
}
