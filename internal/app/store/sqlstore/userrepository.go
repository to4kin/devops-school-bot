package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (telegram_id, first_name, last_name, username, is_admin) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		u.TelegramID,
		u.FirstName,
		u.LastName,
		u.Username,
		u.IsAdmin,
	).Scan(
		&u.ID,
	)
}

func (r *UserRepository) FindByTelegramID(telegramID int64) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, telegram_id, first_name, last_name, username, is_admin FROM users WHERE telegram_id = $1",
		telegramID,
	).Scan(
		&u.ID,
		&u.TelegramID,
		&u.FirstName,
		&u.LastName,
		&u.Username,
		&u.IsAdmin,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}
