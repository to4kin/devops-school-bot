package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// CallbackRepository ...
type CallbackRepository struct {
	store *Store
}

// Create ...
func (r *CallbackRepository) Create(c *model.Callback) error {
	if err := c.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO callback (created, type, type_id, command, list_command) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		c.Created,
		c.Type,
		c.TypeID,
		c.Command,
		c.ListCommand,
	).Scan(
		&c.ID,
	)
}

// FindByID ...
func (r *CallbackRepository) FindByID(callbackID int64) (*model.Callback, error) {
	c := &model.Callback{}
	if err := r.store.db.QueryRow(
		"SELECT id, created, type, type_id, command, list_command FROM callback WHERE id = $1",
		callbackID,
	).Scan(
		&c.ID,
		&c.Created,
		&c.Type,
		&c.TypeID,
		&c.Command,
		&c.ListCommand,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return c, nil
}

// FindByCallback ...
func (r *CallbackRepository) FindByCallback(callback *model.Callback) (*model.Callback, error) {
	c := &model.Callback{}
	if err := r.store.db.QueryRow(`
		SELECT id, created, type, type_id, command, list_command FROM callback
		WHERE type = $1 AND type_id = $2 AND command = $3 AND list_command = $4
		`,
		callback.Type,
		callback.TypeID,
		callback.Command,
		callback.ListCommand,
	).Scan(
		&c.ID,
		&c.Created,
		&c.Type,
		&c.TypeID,
		&c.Command,
		&c.ListCommand,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return c, nil
}
