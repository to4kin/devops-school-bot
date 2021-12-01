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
func (r *CallbackRepository) Create(callback *model.Callback) error {
	if err := callback.Validate(); err != nil {
		return err
	}

	if err := r.store.callbackRepository.FindByCallback(callback); err != nil && err != store.ErrRecordNotFound {
		return err
	}

	if callback.ID != 0 {
		return store.ErrRecordIsExist
	}

	callback.ID = int64(len(r.callbacks) + 1)
	r.callbacks = append(r.callbacks, callback)
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
func (r *CallbackRepository) FindByCallback(callback *model.Callback) error {
	for _, c := range r.callbacks {
		if c.Type == callback.Type && c.TypeID == callback.TypeID &&
			c.Command == callback.Command && c.ListCommand == callback.ListCommand {
			callback.ID = c.ID
			callback.Created = c.Created
			return nil
		}
	}

	return store.ErrRecordNotFound
}
