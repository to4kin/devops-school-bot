package teststore

import (
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// CallbackRepository ...
type CallbackRepository struct {
	store     *Store
	callbacks []*model.Callback
}

// Create ...
func (r *CallbackRepository) Create(c *model.Callback) error {
	if err := c.Validate(); err != nil {
		return err
	}

	callback, err := r.store.callbackRepository.FindByCallback(c)
	if err != nil && err != store.ErrRecordNotFound {
		return err
	}

	if callback != nil {
		return store.ErrRecordIsExist
	}

	c.ID = int64(len(r.callbacks) + 1)
	r.callbacks = append(r.callbacks, c)
	return nil
}

// FindByID ...
func (r *CallbackRepository) FindByID(callbackID int64) (*model.Callback, error) {
	for _, c := range r.callbacks {
		if c.ID == callbackID {
			return c, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

// FindByCallback ...
func (r *CallbackRepository) FindByCallback(callback *model.Callback) (*model.Callback, error) {
	for _, c := range r.callbacks {
		if c.Type == callback.Type && c.TypeID == callback.TypeID &&
			c.Command == callback.Command && c.ListCommand == callback.ListCommand {
			return c, nil
		}
	}

	return nil, store.ErrRecordNotFound
}
