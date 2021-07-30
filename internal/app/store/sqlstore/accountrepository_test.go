package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/sqlstore"
)

func TestAccountRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("account")

	s := sqlstore.New(db)
	a := model.TestAccount(t)

	assert.NoError(t, s.Account().Create(a))
	assert.NotNil(t, a)
}

func TestAccountRepository_FindByTelegramID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("account")

	s := sqlstore.New(db)
	a := model.TestAccount(t)

	_, err := s.Account().FindByTelegramID(a.TelegramID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(a))

	account, err := s.Account().FindByTelegramID(a.TelegramID)
	assert.NoError(t, err)
	assert.NotNil(t, account)
}
