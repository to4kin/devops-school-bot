package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/sqlstore"
)

func TestCallbackRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("callback")

	s := sqlstore.New(db)
	callback := model.TestAccountCallback(t)

	assert.NoError(t, s.Callback().Create(callback))
	assert.NotNil(t, callback)

	err := s.Callback().Create(callback)
	assert.Error(t, err)
}

func TestCallbackRepository_FindByID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("callback")

	s := sqlstore.New(db)
	testCallback := model.TestAccountCallback(t)

	_, err := s.Callback().FindByID(testCallback.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Callback().Create(testCallback))

	callback, err := s.Callback().FindByID(testCallback.ID)
	assert.NoError(t, err)
	assert.NotNil(t, callback)
}

func TestCallbackRepository_FindByCallback(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("callback")

	s := sqlstore.New(db)
	testCallback := model.TestAccountCallback(t)

	_, err := s.Callback().FindByCallback(testCallback)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Callback().Create(testCallback))

	callback, err := s.Callback().FindByCallback(testCallback)
	assert.NoError(t, err)
	assert.NotNil(t, callback)
}
