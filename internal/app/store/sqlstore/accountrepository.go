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

// Update ...
func (r *AccountRepository) Update(a *model.Account) error {
	if err := a.Validate(); err != nil {
		return err
	}

	if err := r.store.db.QueryRow(
		"UPDATE account SET first_name = $2, last_name = $3, username = $4, superuser = $5 WHERE telegram_id = $1 RETURNING id",
		a.TelegramID,
		a.FirstName,
		a.LastName,
		a.Username,
		a.Superuser,
	).Scan(
		&a.ID,
	); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}

		return err
	}

	return nil
}

// FindAll ...
func (r *AccountRepository) FindAll() ([]*model.Account, error) {
	rowsCount := 0
	accounts := []*model.Account{}

	rows, err := r.store.db.Query(`
		SELECT id, created, telegram_id, first_name, last_name, username, superuser FROM account ORDER BY username ASC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsCount++

		a := &model.Account{}

		if err := rows.Scan(
			&a.ID,
			&a.Created,
			&a.TelegramID,
			&a.FirstName,
			&a.LastName,
			&a.Username,
			&a.Superuser,
		); err != nil {
			return nil, err
		}

		accounts = append(accounts, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if rowsCount == 0 {
		return nil, store.ErrRecordNotFound
	}

	return accounts, nil
}

// FindByID ...
func (r *AccountRepository) FindByID(id int64) (*model.Account, error) {
	a := &model.Account{}
	if err := r.store.db.QueryRow(
		"SELECT id, created, telegram_id, first_name, last_name, username, superuser FROM account WHERE id = $1",
		id,
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
