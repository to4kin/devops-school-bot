package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// AccountRepository ...
type AccountRepository struct {
	store *Store
}

// Create ...
func (r *AccountRepository) Create(a *model.Account) error {
	if err := a.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO account (created, telegram_id, first_name, last_name, username, superuser) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		a.Created,
		a.TelegramID,
		a.FirstName,
		a.LastName,
		a.Username,
		a.Superuser,
	).Scan(
		&a.ID,
	)
}

// FindByTelegramID ...
func (r *AccountRepository) FindByTelegramID(telegramID int64) (*model.Account, error) {
	a := &model.Account{}
	if err := r.store.db.QueryRow(
		"SELECT id, created, telegram_id, first_name, last_name, username, superuser FROM account WHERE telegram_id = $1",
		telegramID,
	).Scan(
		&a.ID,
		&a.Created,
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
